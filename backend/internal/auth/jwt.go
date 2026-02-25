package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var jwtSecret []byte

// init runs automatically when the package is imported
func init() {
	// Load .env if it hasn't been loaded yet
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found in auth package, relying on environment variables")
	}
	
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in environment variables")
	}
	jwtSecret = []byte(secret)
}

func getSecretKey() []byte {
	return jwtSecret
}

func GenerateToken(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecretKey())
}