package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/devsherkhane/trello-clone/internal/database"
	"github.com/gin-gonic/gin"
)

func (h *APIHandler) UpdateTheme(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	var input struct {
		Theme string `json:"theme" binding:"required,oneof=light dark auto"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := database.DB.Exec("UPDATE users SET theme = ? WHERE id = ?", input.Theme, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theme updated"})
}

func (h *APIHandler) UploadAvatar(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	
	uploadDir := "uploads/avatars"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	filename := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	_, err = database.DB.Exec("UPDATE users SET avatar_url = ? WHERE id = ?", filename, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar updated", "avatar_url": filename})
}
