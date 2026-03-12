package repository

import (
	"database/sql"
	"github.com/devsherkhane/trello-clone/internal/models"
)

type UserRepository interface {
	Create(username, email, passwordHash string) (int64, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetByResetToken(token string) (*models.User, error)
	Update(user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(username, email, passwordHash string) (int64, error) {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	res, err := r.db.Exec(query, username, email, passwordHash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := "SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = ?"
	var u models.User
	err := r.db.QueryRow(query, email).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := "SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = ?"
	var u models.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByResetToken(token string) (*models.User, error) {
	query := "SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE reset_token = ? AND reset_expires > NOW()"
	var u models.User
	err := r.db.QueryRow(query, token).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(u *models.User) error {
	query := "UPDATE users SET username = ?, password_hash = ?, reset_token = ?, reset_expires = ? WHERE id = ?"
	_, err := r.db.Exec(query, u.Username, u.PasswordHash, u.ResetToken, u.ResetExpires, u.ID)
	return err
}
