package models

import "time"

// User represents the person using the app
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Scrambled password, hidden from JSON
	ResetToken   string    `json:"-"`
	ResetExpires time.Time `json:"-"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

// Board is the top-level container (e.g., "Project Alpha")
type Board struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	UserID     int    `json:"user_id"`
	IsArchived bool   `json:"is_archived"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// List represents a column (e.g., "To Do", "Doing")
type List struct {
	ID       int    `json:"id"`
	BoardID  int    `json:"board_id"`
	Title    string `json:"title"`
	Position int    `json:"position"` // Determines the left-to-right order
}

// Card is the individual task inside a list
type Card struct {
	ID          int    `json:"id"`
	ListID      int    `json:"list_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Position    int    `json:"position"`
	DueDate     string `json:"due_date"`
	LabelColor  string `json:"label_color"`
}

// Attachment represents a file linked to a card
type Attachment struct {
	ID        int    `json:"id"`
	CardID    int    `json:"card_id"`
	UserID    int    `json:"user_id"`
	FilePath  string `json:"file_path"`
	Filename  string `json:"filename"`
	CreatedAt string `json:"created_at"`
}

// Comment represents a discussion entry on a card
type Comment struct {
	ID        int    `json:"id"`
	CardID    int    `json:"card_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username,omitempty"` // Temporarily mapped via UpdatedAt hack in repo
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ActivityLog represents board action history
type ActivityLog struct {
	ID        int    `json:"id"`
	BoardID   int    `json:"board_id"`
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	CreatedAt string `json:"created_at"`
}

// Notification represents a real-time alert for a user
type Notification struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	FromUser string `json:"from_user"`
	CardID   int    `json:"card_id"`
}
