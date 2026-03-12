package handlers

import (
	"net/http"
	"strconv"

	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (h *APIHandler) CreateComment(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	
	var input struct {
		CardID int    `json:"card_id" binding:"required"`
		Text   string `json:"text" binding:"required,max=5000"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	comment, err := h.CommentService.CreateComment(input.CardID, userID, input.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *APIHandler) GetComments(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))

	comments, err := h.CommentService.GetCommentsByCard(cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	if comments == nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	c.JSON(http.StatusOK, comments)
}
