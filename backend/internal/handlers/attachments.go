package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/devsherkhane/trello-clone/internal/database"
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadAttachment(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	cardID := c.PostForm("card_id")

	// 1. Security Check: Verify user owns the board the card belongs to
	var ownerID int
	err := database.DB.QueryRow(`
		SELECT b.owner_id FROM boards b
		JOIN lists l ON b.id = l.board_id
		JOIN cards c ON l.id = c.list_id
		WHERE c.id = ?`, cardID).Scan(&ownerID)

	if err != nil || ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to add attachments to this card"})
		return
	}

	// 2. Process File
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Create unique filename to prevent overwriting
	extension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d-%s%s", time.Now().Unix(), uuid.New().String(), extension)
	uploadPath := filepath.Join("uploads", newFileName)

	// 3. Save to local filesystem
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 4. Save to Database
	query := "INSERT INTO card_attachments (card_id, file_path, file_name) VALUES (?, ?, ?)"
	_, err = database.DB.Exec(query, cardID, uploadPath, file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record attachment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "path": uploadPath})
}

// GetAttachmentsByCard retrieves all files for a specific card
func GetAttachmentsByCard(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	cardID := c.Param("id")

	// Security Check: Verify user owns the board this card belongs to
	var ownerID int
	err := database.DB.QueryRow(`
        SELECT b.owner_id FROM boards b
        JOIN lists l ON b.id = l.board_id
        JOIN cards c ON l.id = c.list_id
        WHERE c.id = ?`, cardID).Scan(&ownerID)

	if err != nil || ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access to these attachments"})
		return
	}

	var attachments []models.CardAttachment
	query := "SELECT id, card_id, file_path, file_name, uploaded_at FROM card_attachments WHERE card_id = ?"

	rows, err := database.DB.Query(query, cardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attachments"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var a models.CardAttachment
		if err := rows.Scan(&a.ID, &a.CardID, &a.FilePath, &a.FileName, &a.UploadedAt); err != nil {
			continue
		}
		attachments = append(attachments, a)
	}

	c.JSON(http.StatusOK, attachments)
}
