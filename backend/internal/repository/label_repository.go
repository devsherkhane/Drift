package repository

import (
	"database/sql"
)

type LabelRepository interface {
	AddLabelToCard(cardID, labelID int) error
}

type labelRepository struct {
	db *sql.DB
}

func NewLabelRepository(db *sql.DB) LabelRepository {
	return &labelRepository{db: db}
}

func (r *labelRepository) AddLabelToCard(cardID, labelID int) error {
	_, err := r.db.Exec("INSERT INTO card_labels (card_id, label_id) VALUES (?, ?)", cardID, labelID)
	return err
}
