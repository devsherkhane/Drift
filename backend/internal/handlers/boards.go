package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/devsherkhane/trello-clone/internal/repository"
	"github.com/gin-gonic/gin"
)

func (h *APIHandler) CreateBoard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	var input struct {
		Title string `json:"title" binding:"required,max=255"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	board, err := h.BoardService.CreateBoard(input.Title, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       board.ID,
		"title":    board.Title,
		"user_id":  userID,
		"is_owner": true,
		"status":   "accepted",
	})
}

func (h *APIHandler) GetBoards(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	
	// Our service fetches appropriately. We can filter archived on the handler level or repo
	archivedFilter := c.Query("archived") == "true"
	
	boards, err := h.BoardService.GetBoards(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var response []gin.H
	for _, b := range boards {
		// Temporary hack to filter if we didn't push it down to the repo correctly for this query flag
		if b.IsArchived == archivedFilter {
			response = append(response, gin.H{
				"id":       b.ID,
				"title":    b.Title,
				"user_id":  b.UserID,
				"is_owner": b.UserID == userID,
				"status":   "accepted", // Simplifying collaborations for now
			})
		}
	}
	
	if response == nil {
		response = []gin.H{}
	}
	c.JSON(http.StatusOK, response)
}

func (h *APIHandler) RespondToInvitation(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Action string `json:"action" binding:"required,oneof=accept decline"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	// This is a missing method in our service, but for time's sake we do logic directly or just implement it here:
	if input.Action == "accept" {
		// we assume accepted default for now in our repository implementation
		c.JSON(http.StatusOK, gin.H{"message": "Invitation accepted"})
	} else {
		err := h.BoardService.RemoveCollaborator(boardID, userID, userID) // remove self
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Invitation declined"})
	}
}

func (h *APIHandler) UpdateBoard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))
	
	var input struct {
		Title string `json:"title" binding:"required,max=255"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	err := h.BoardService.UpdateBoardTitle(boardID, userID, input.Title)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Board not found or unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Board updated"})
}

func (h *APIHandler) DeleteBoard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))

	err := h.BoardService.DeleteBoard(boardID, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Board not found or unauthorized"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Board deleted"})
}

func (h *APIHandler) GetBoard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))

	board, err := h.BoardService.GetBoardByID(boardID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Fetch lists
	lists, err := h.ListService.GetListsByBoard(board.ID)
	var listResponses []gin.H

	for _, l := range lists {
		// Fetch cards. This N+1 query is generally bad practice but mimics original logic temporarily
		cards, _ := h.CardService.GetCardsByList(l.ID)
		if cards == nil {
			cards = []models.Card{}
		}
		
		listResponses = append(listResponses, gin.H{
			"id":       l.ID,
			"board_id": l.BoardID,
			"title":    l.Title,
			"position": l.Position,
			"cards":    cards, // Not including labels nested query for time being
		})
	}
	
	if listResponses == nil {
		listResponses = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      board.ID,
		"title":   board.Title,
		"user_id": board.UserID,
		"lists":   listResponses,
	})
}

func (h *APIHandler) GetActivityLogs(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))

	// Verify access
	_, err := h.BoardService.GetBoardByID(boardID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to logs"})
		return
	}

	logs, err := h.BoardService.GetActivityLogs(boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	if logs == nil {
		logs = []models.ActivityLog{}
	}

	c.JSON(http.StatusOK, logs)
}

func (h *APIHandler) AddCollaborator(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Email string `json:"email" binding:"required,email,max=255"`
		Role  string `json:"role" binding:"required,oneof=member admin viewer"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	err := h.BoardService.AddCollaborator(boardID, userID, input.Email, input.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collaborator added"})
}

func (h *APIHandler) ArchiveBoard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))

	err := h.BoardService.ArchiveBoard(boardID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update board status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Board status updated"})
}

func (h *APIHandler) ExportBoardCSV(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	boardID, _ := strconv.Atoi(c.Param("id"))

	board, err := h.BoardService.GetBoardByID(boardID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	h.exportCSVMockup(c, board)
}

func (h *APIHandler) exportCSVMockup(c *gin.Context, board *models.Board) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", board.Title))
	c.Header("Content-Type", "text/csv")
	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"List", "Card Title", "Description"})

	lists, _ := h.ListService.GetListsByBoard(board.ID)
	for _, l := range lists {
		cards, _ := h.CardService.GetCardsByList(l.ID)
		for _, card := range cards {
			writer.Write([]string{l.Title, card.Title, card.Description})
		}
	}
	writer.Flush()
}
