package utils

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/devsherkhane/trello-clone/internal/auth"
)

// GenerateJWT acts as an adapter for the internal auth package generator
func GenerateJWT(userID int, email string) (string, error) {
	// The original auth.GenerateToken just takes userID
	return auth.GenerateToken(userID)
}

// GenerateResetToken safely generates cryptographically secure reset strings
func GenerateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
