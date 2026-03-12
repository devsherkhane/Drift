package repository

import (
	"database/sql"
	"github.com/devsherkhane/trello-clone/internal/models"
)

type ListRepository interface {
	GetByBoardID(boardID int) ([]models.List, error)
	Create(boardID int, title string) (int64, error)
	UpdateTitle(listID int, title string) error
	Delete(listID int) error
}

type listRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) ListRepository {
	return &listRepository{db: db}
}

func (r *listRepository) GetByBoardID(boardID int) ([]models.List, error) {
	query := "SELECT id, board_id, title, position FROM lists WHERE board_id = ? ORDER BY position ASC"
	rows, err := r.db.Query(query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []models.List
	for rows.Next() {
		var l models.List
		if err := rows.Scan(&l.ID, &l.BoardID, &l.Title, &l.Position); err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}
	return lists, nil
}

func (r *listRepository) Create(boardID int, title string) (int64, error) {
	// Simple position logic: max + 1
	var pos int
	r.db.QueryRow("SELECT COALESCE(MAX(position), 0) + 1 FROM lists WHERE board_id = ?", boardID).Scan(&pos)

	query := "INSERT INTO lists (board_id, title, position) VALUES (?, ?, ?)"
	res, err := r.db.Exec(query, boardID, title, pos)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *listRepository) UpdateTitle(listID int, title string) error {
	_, err := r.db.Exec("UPDATE lists SET title = ? WHERE id = ?", title, listID)
	return err
}

func (r *listRepository) Delete(listID int) error {
	_, err := r.db.Exec("DELETE FROM lists WHERE id = ?", listID)
	return err
}
