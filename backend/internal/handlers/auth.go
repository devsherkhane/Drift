package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log" // Ensure log is imported
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/devsherkhane/trello-clone/internal/auth"
	"github.com/devsherkhane/trello-clone/internal/database"
	"github.com/devsherkhane/trello-clone/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)
// LoginUser godoc
// @Summary Authenticate a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body models.LoginInput true "User Credentials"
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Create a new struct specifically for Login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	// Gin's ShouldBindJSON already checks for required fields and email format
	// based on the struct tags you defined.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: ensure email is valid and password is at least 6 characters"})
		return
	}

	// 1. Check if the user already exists
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? OR username = ?)"
	err := database.DB.QueryRow(checkQuery, req.Email, req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database check failed"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or Email already taken"})
		return
	}

	// 2. Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Insert the new user
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	result, err := database.DB.Exec(query, req.Username, req.Email, string(hashed))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": id,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	var id int
	var hashedPwd string
	// Ensure the column names (id, password_hash) match your MySQL table exactly
	err := database.DB.QueryRow("SELECT user_id, password_hash FROM users WHERE email = ?", req.Email).Scan(&id, &hashedPwd)
	if err != nil {
		log.Printf("Login Error: %v", err) // CHECK YOUR TERMINAL FOR THIS LOG
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := auth.GenerateToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// UpdateProfile allows a user to change their username or password
func UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	var input struct {
		Username string `json:"username"`
		Password string `json:"password" binding:"omitempty,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Update Username if provided
	if input.Username != "" {
		_, err := database.DB.Exec("UPDATE users SET username = ? WHERE user_id = ?", input.Username, userID)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
	}

	// Update Password if provided (hashed)
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		_, err := database.DB.Exec("UPDATE users SET password_hash = ? WHERE user_id = ?", string(hashed), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func UpdateTheme(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	var input struct {
		Theme string `json:"theme" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Theme is required"})
		return
	}

	_, err := database.DB.Exec("UPDATE users SET theme_preference = ? WHERE user_id = ?", input.Theme, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save preference"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Theme updated"})
}

func Setup2FA(c *gin.Context) {
	userID := c.MustGet("userID").(int)
	var email string
	database.DB.QueryRow("SELECT email FROM users WHERE user_id = ?", userID).Scan(&email)

	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "TrelloClone",
		AccountName: email,
	})

	// Store secret temporarily/permanently depending on your flow
	_, err := database.DB.Exec("UPDATE users SET tfa_secret = ? WHERE user_id = ?", key.Secret(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret": key.Secret(),
		"url":    key.URL(), // Used to generate a QR code on the frontend
	})
}

func UploadAvatar(c *gin.Context) {
	userID := c.MustGet("userID").(int)

	// 1. Get the file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image uploaded"})
		return
	}

	// 2. Validate file type (Images only)
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, JPEG, and PNG are allowed"})
		return
	}

	// 3. Create unique filename
	newFileName := fmt.Sprintf("avatar-%d-%d%s", userID, time.Now().Unix(), ext)
	uploadPath := filepath.Join("uploads", "avatars", newFileName)

	// Ensure directory exists
	os.MkdirAll(filepath.Join("uploads", "avatars"), os.ModePerm)

	// 4. Fetch old avatar to delete it later
	var oldAvatar string
	database.DB.QueryRow("SELECT avatar_url FROM users WHERE user_id = ?", userID).Scan(&oldAvatar)

	// 5. Save the new file
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// 6. Update database
	_, err = database.DB.Exec("UPDATE users SET avatar_url = ? WHERE user_id = ?", uploadPath, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed"})
		return
	}

	// 7. Cleanup: Delete old file if it exists
	if oldAvatar != "" {
		os.Remove(oldAvatar)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Avatar updated successfully",
		"avatar_url": uploadPath,
	})
}

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email is required"})
		return
	}

	// 1. Check if user exists
	var userID int
	err := database.DB.QueryRow("SELECT user_id FROM users WHERE email = ?", input.Email).Scan(&userID)
	if err != nil {
		// Security tip: Don't reveal if the email exists.
		// Just say "If an account exists, an email has been sent."
		c.JSON(http.StatusOK, gin.H{"message": "If this email is registered, a reset link has been sent."})
		return
	}

	// 2. Generate a secure random token
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	expiry := time.Now().Add(1 * time.Hour) // Token valid for 1 hour

	// 3. Store token in DB
	_, err = database.DB.Exec("INSERT INTO password_resets (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset link"})
		return
	}

	// 4. Send Email
	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click the link below to reset your password. It expires in 1 hour.\n\n%s", resetLink)

	go utils.SendEmail(input.Email, subject, body)

	c.JSON(http.StatusOK, gin.H{"message": "Reset link sent."})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token and valid password are required"})
		return
	}

	// 1. Validate Token and Expiry
	var userID int
	err := database.DB.QueryRow("SELECT user_id FROM password_resets WHERE token = ? AND expires_at > NOW()", input.Token).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// 2. Hash New Password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	// 3. Update User Password and Delete the token
	tx, _ := database.DB.Begin()
	tx.Exec("UPDATE users SET password_hash = ? WHERE user_id = ?", string(hashed), userID)
	tx.Exec("DELETE FROM password_resets WHERE token = ?", input.Token)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully."})
}
