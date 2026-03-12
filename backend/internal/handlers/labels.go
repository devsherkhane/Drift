package handlers

import (
	"net/http"
	"strconv"

	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (h *APIHandler) AddLabelToCard(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		LabelID int `json:"label_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": middleware.FormatValidationErrors(err)})
		return
	}

	err := h.LabelService.AddLabel(cardID, input.LabelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add label"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Label added to card"})
}
