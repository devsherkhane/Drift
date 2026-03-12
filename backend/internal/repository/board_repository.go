package repository

import (
	"database/sql"
	"errors"

	"github.com/devsherkhane/trello-clone/internal/models"
)

var ErrNotFound = errors.New("resource not found")

type BoardRepository interface {
	Create(title string, userID int) (int64, error)
	GetByUserID(userID int) ([]models.Board, error)
	GetByID(boardID, userID int) (*models.Board, error)
	UpdateTitle(boardID, userID int, title string) error
	Delete(boardID, userID int) error
	Archive(boardID, userID int) error
	GetActivityLogs(boardID int) ([]models.ActivityLog, error)

	// Collaborators
	AddCollaborator(boardID int, email, role string) error
	GetCollaborators(boardID int) ([]models.User, error)
	UpdateCollaboratorRole(boardID, collaboratorID int, role string) error
	RemoveCollaborator(boardID, collaboratorID int) error
}

type boardRepository struct {
	db *sql.DB
}

func NewBoardRepository(db *sql.DB) BoardRepository {
	return &boardRepository{db: db}
}

func (r *boardRepository) Create(title string, userID int) (int64, error) {
	query := "INSERT INTO boards (title, user_id) VALUES (?, ?)"
	result, err := r.db.Exec(query, title, userID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *boardRepository) GetByUserID(userID int) ([]models.Board, error) {
	query := `
		SELECT b.id, b.title, b.user_id, b.is_archived, b.created_at, b.updated_at,
		       CASE WHEN bc.board_id IS NOT NULL THEN true ELSE false END as is_collaborator
		FROM boards b
		LEFT JOIN board_collaborators bc ON b.id = bc.board_id AND bc.user_id = ?
		WHERE (b.user_id = ? OR bc.user_id = ?) AND b.is_archived = false
	`
	rows, err := r.db.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []models.Board
	for rows.Next() {
		var b models.Board
		var isCollaborator bool
		if err := rows.Scan(&b.ID, &b.Title, &b.UserID, &b.IsArchived, &b.CreatedAt, &b.UpdatedAt, &isCollaborator); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}
	return boards, nil
}

func (r *boardRepository) GetByID(boardID, userID int) (*models.Board, error) {
	query := `
		SELECT b.id, b.title, b.user_id, b.is_archived, b.created_at, b.updated_at
		FROM boards b
		LEFT JOIN board_collaborators bc ON b.id = bc.board_id
		WHERE b.id = ? AND (b.user_id = ? OR bc.user_id = ?)
	`
	var board models.Board
	err := r.db.QueryRow(query, boardID, userID, userID).Scan(&board.ID, &board.Title, &board.UserID, &board.IsArchived, &board.CreatedAt, &board.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &board, nil
}

func (r *boardRepository) UpdateTitle(boardID, userID int, title string) error {
	query := "UPDATE boards SET title = ? WHERE id = ? AND user_id = ?"
	res, err := r.db.Exec(query, title, boardID, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *boardRepository) Delete(boardID, userID int) error {
	query := "DELETE FROM boards WHERE id = ? AND user_id = ?"
	res, err := r.db.Exec(query, boardID, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *boardRepository) Archive(boardID, userID int) error {
	query := "UPDATE boards SET is_archived = true WHERE id = ? AND user_id = ?"
	res, err := r.db.Exec(query, boardID, userID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *boardRepository) GetActivityLogs(boardID int) ([]models.ActivityLog, error) {
	query := `
		SELECT id, board_id, user_id, action, created_at 
		FROM activity_logs 
		WHERE board_id = ? 
		ORDER BY created_at DESC LIMIT 50
	`
	rows, err := r.db.Query(query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var log models.ActivityLog
		if err := rows.Scan(&log.ID, &log.BoardID, &log.UserID, &log.Action, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// Collaborators implementation
func (r *boardRepository) AddCollaborator(boardID int, email, role string) error {
	// First get the user ID from the email
	var invitedUserID int
	err := r.db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&invitedUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return err
	}

	// Insert into board_collaborators (ignore duplicate if already exists)
	query := "INSERT INTO board_collaborators (board_id, user_id, role) VALUES (?, ?, ?)"
	_, err = r.db.Exec(query, boardID, invitedUserID, role)
	return err
}

func (r *boardRepository) GetCollaborators(boardID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, bc.role 
		FROM users u 
		JOIN board_collaborators bc ON u.id = bc.user_id 
		WHERE bc.board_id = ?
	`
	rows, err := r.db.Query(query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collaborators []models.User
	for rows.Next() {
		var u models.User
		var role string
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &role); err != nil {
			return nil, err
		}
		// Overriding password hash field with role temporarily just for output formatting,
		// though typically a specific dto/struct should be used.
		u.PasswordHash = role // using this field temporarily to hold role for json response
		collaborators = append(collaborators, u)
	}
	return collaborators, nil
}

func (r *boardRepository) UpdateCollaboratorRole(boardID, collaboratorID int, role string) error {
	query := "UPDATE board_collaborators SET role = ? WHERE board_id = ? AND user_id = ?"
	res, err := r.db.Exec(query, role, boardID, collaboratorID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *boardRepository) RemoveCollaborator(boardID, collaboratorID int) error {
	query := "DELETE FROM board_collaborators WHERE board_id = ? AND user_id = ?"
	res, err := r.db.Exec(query, boardID, collaboratorID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
