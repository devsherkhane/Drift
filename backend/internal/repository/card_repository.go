package repository

import (
	"database/sql"
	"github.com/devsherkhane/trello-clone/internal/models"
)

type CardRepository interface {
	GetByListID(listID int) ([]models.Card, error)
	Create(listID int, title string) (int64, error)
	Update(card *models.Card) error
	Move(cardID, newListID, newPosition int) error
	Delete(cardID int) error
	GetByID(cardID int) (*models.Card, error)
}

type cardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) GetByListID(listID int) ([]models.Card, error) {
	query := "SELECT id, list_id, title, description, position, due_date, label_color FROM cards WHERE list_id = ? ORDER BY position ASC"
	rows, err := r.db.Query(query, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		// due_date can be null
		var dueDate sql.NullString
		var labelColor sql.NullString
		if err := rows.Scan(&c.ID, &c.ListID, &c.Title, &c.Description, &c.Position, &dueDate, &labelColor); err != nil {
			return nil, err
		}
		if dueDate.Valid {
			c.DueDate = dueDate.String
		}
		if labelColor.Valid {
			c.LabelColor = labelColor.String
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *cardRepository) Create(listID int, title string) (int64, error) {
	var pos int
	r.db.QueryRow("SELECT COALESCE(MAX(position), 0) + 1 FROM cards WHERE list_id = ?", listID).Scan(&pos)

	query := "INSERT INTO cards (list_id, title, position) VALUES (?, ?, ?)"
	res, err := r.db.Exec(query, listID, title, pos)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *cardRepository) Update(c *models.Card) error {
	query := `
		UPDATE cards 
		SET title = ?, description = ?, due_date = ?, label_color = ? 
		WHERE id = ?
	`
	_, err := r.db.Exec(query, c.Title, c.Description, c.DueDate, c.LabelColor, c.ID)
	return err
}

func (r *cardRepository) Move(cardID, newListID, newPosition int) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// 1. Get current list_id and position
	var oldListID, oldPosition int
	err = tx.QueryRow("SELECT list_id, position FROM cards WHERE id = ?", cardID).Scan(&oldListID, &oldPosition)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Shift other cards
	if oldListID == newListID {
		if oldPosition < newPosition {
			// moving down
			_, err = tx.Exec("UPDATE cards SET position = position - 1 WHERE list_id = ? AND position > ? AND position <= ?", oldListID, oldPosition, newPosition)
		} else {
			// moving up
			_, err = tx.Exec("UPDATE cards SET position = position + 1 WHERE list_id = ? AND position >= ? AND position < ?", oldListID, newPosition, oldPosition)
		}
	} else {
		// moving between lists
		_, err = tx.Exec("UPDATE cards SET position = position - 1 WHERE list_id = ? AND position > ?", oldListID, oldPosition)
		if err == nil {
			_, err = tx.Exec("UPDATE cards SET position = position + 1 WHERE list_id = ? AND position >= ?", newListID, newPosition)
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Update the target card
	_, err = tx.Exec("UPDATE cards SET list_id = ?, position = ? WHERE id = ?", newListID, newPosition, cardID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *cardRepository) Delete(cardID int) error {
	_, err := r.db.Exec("DELETE FROM cards WHERE id = ?", cardID)
	return err
}

func (r *cardRepository) GetByID(cardID int) (*models.Card, error) {
	query := "SELECT id, list_id, title, description, position, due_date, label_color FROM cards WHERE id = ?"
	var c models.Card
	var dueDate, labelColor sql.NullString

	err := r.db.QueryRow(query, cardID).Scan(&c.ID, &c.ListID, &c.Title, &c.Description, &c.Position, &dueDate, &labelColor)
	if err != nil {
		return nil, err
	}
	if dueDate.Valid {
		c.DueDate = dueDate.String
	}
	if labelColor.Valid {
		c.LabelColor = labelColor.String
	}
	return &c, nil
}
