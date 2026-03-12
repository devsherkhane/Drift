package handlers

import (
	"net/http"
	"strconv"

	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (h *APIHandler) CreateCard(c *gin.Context) {
	var input struct {
		ListID int    `json:"list_id" binding:"required"`
		Title  string `json:"title" binding:"required,max=255"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	card, err := h.CardService.CreateCard(input.ListID, input.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, card)
}

func (h *APIHandler) GetCards(c *gin.Context) {
	listIDStr := c.Query("list_id")
	if listIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "list_id is required"})
		return
	}
	listID, _ := strconv.Atoi(listIDStr)

	cards, err := h.CardService.GetCardsByList(listID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cards"})
		return
	}

	if cards == nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	c.JSON(http.StatusOK, cards)
}

func (h *APIHandler) UpdateCard(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Title       string `json:"title" binding:"omitempty,max=255"`
		Description string `json:"description" binding:"omitempty,max=5000"`
		DueDate     string `json:"due_date"`
		LabelColor  string `json:"label_color" binding:"omitempty,max=50"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	card, err := h.CardService.GetCardByID(cardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	if input.Title != "" { card.Title = input.Title }
	card.Description = input.Description
	card.DueDate = input.DueDate
	card.LabelColor = input.LabelColor

	err = h.CardService.UpdateCard(card)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card updated successfully"})
}

func (h *APIHandler) MoveCard(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		NewListID   int `json:"new_list_id" binding:"required"`
		NewPosition int `json:"new_position"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	err := h.CardService.MoveCard(cardID, input.NewListID, input.NewPosition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card moved successfully"})
}

func (h *APIHandler) DeleteCard(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))

	err := h.CardService.DeleteCard(cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deleted successfully"})
}

func (h *APIHandler) SearchCards(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	userID := c.MustGet("userID").(int)
	
	var boardIDPtr *int
	if bID := c.Query("board_id"); bID != "" {
		parsed, _ := strconv.Atoi(bID)
		boardIDPtr = &parsed
	}

	results, err := h.SearchService.Search(query, boardIDPtr, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, results)
}
