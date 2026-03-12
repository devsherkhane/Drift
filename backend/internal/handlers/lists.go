package handlers

import (
	"net/http"
	"strconv"

	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (h *APIHandler) CreateList(c *gin.Context) {
	var input struct {
		BoardID int    `json:"board_id" binding:"required"`
		Title   string `json:"title" binding:"required,max=255"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	// Security: should verify user has access to board
	userID := c.MustGet("userID").(int)
	_, err := h.BoardService.GetBoardByID(input.BoardID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to board"})
		return
	}

	list, err := h.ListService.CreateList(input.BoardID, input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create list"})
		return
	}

	c.JSON(http.StatusCreated, list)
}

func (h *APIHandler) GetLists(c *gin.Context) {
	boardIDStr := c.Query("board_id")
	if boardIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "board_id is required"})
		return
	}
	boardID, _ := strconv.Atoi(boardIDStr)

	// Verify Board Access
	userID := c.MustGet("userID").(int)
	_, err := h.BoardService.GetBoardByID(boardID, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	lists, err := h.ListService.GetListsByBoard(boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch lists"})
		return
	}

	if lists == nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (h *APIHandler) UpdateList(c *gin.Context) {
	listID, _ := strconv.Atoi(c.Param("id"))
	
	var input struct {
		Title   string `json:"title" binding:"required,max=255"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	err := h.ListService.UpdateListTitle(listID, input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List updated successfully"})
}

func (h *APIHandler) DeleteList(c *gin.Context) {
	listID, _ := strconv.Atoi(c.Param("id"))

	err := h.ListService.DeleteList(listID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete list"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "List deleted successfully"})
}
