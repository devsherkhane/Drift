package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) UploadAttachment(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	cardIDStr := c.PostForm("card_id")
	cardID, err := strconv.Atoi(cardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card_id"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	Filename := file.Filename
	
	// Ensure uploads directory exists
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	filePath := filepath.Join(uploadDir, Filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	attachment, err := h.AttachmentService.AddAttachment(cardID, userID, Filename, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *APIHandler) GetAttachments(c *gin.Context) {
	cardID, _ := strconv.Atoi(c.Param("id"))

	attachments, err := h.AttachmentService.GetAttachmentsByCard(cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attachments"})
		return
	}

	if attachments == nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	c.JSON(http.StatusOK, attachments)
}
