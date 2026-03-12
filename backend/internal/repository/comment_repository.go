package repository

import (
	"database/sql"
	"github.com/devsherkhane/trello-clone/internal/models"
)

type CommentRepository interface {
	Create(cardID, userID int, text string) (int64, error)
	GetByCardID(cardID int) ([]models.Comment, error)
}

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(cardID, userID int, text string) (int64, error) {
	query := "INSERT INTO comments (card_id, user_id, text) VALUES (?, ?, ?)"
	res, err := r.db.Exec(query, cardID, userID, text)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *commentRepository) GetByCardID(cardID int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.card_id, c.user_id, c.text, c.created_at, u.username
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.card_id = ?
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.Query(query, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		var username string
		if err := rows.Scan(&c.ID, &c.CardID, &c.UserID, &c.Text, &c.CreatedAt, &username); err != nil {
			return nil, err
		}
		// Storing username in memory struct, in a real app would use a DTO
		c.UpdatedAt = username // Using UpdatedAt string field for username temporarily in frontend
		comments = append(comments, c)
	}
	return comments, nil
}
