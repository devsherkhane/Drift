package repository

import (
	"database/sql"
	"github.com/devsherkhane/trello-clone/internal/models"
)

type AttachmentRepository interface {
	Create(cardID, userID int, filename, filePath string) (int64, error)
	GetByCardID(cardID int) ([]models.Attachment, error)
}

type attachmentRepository struct {
	db *sql.DB
}

func NewAttachmentRepository(db *sql.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(cardID, userID int, filename, filePath string) (int64, error) {
	query := "INSERT INTO attachments (card_id, user_id, filename, file_path) VALUES (?, ?, ?, ?)"
	res, err := r.db.Exec(query, cardID, userID, filename, filePath)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *attachmentRepository) GetByCardID(cardID int) ([]models.Attachment, error) {
	query := "SELECT id, card_id, user_id, filename, file_path, created_at FROM attachments WHERE card_id = ?"
	rows, err := r.db.Query(query, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.CardID, &a.UserID, &a.Filename, &a.FilePath, &a.CreatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}
