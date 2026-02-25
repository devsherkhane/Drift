package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found in database package")
	}
}
func InitDB() {
	var err error
	// Use os.Getenv to retrieve credentials from the environment
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	fmt.Printf("Connecting to DB: %s at %s\n", dbName, host)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, dbName)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
}
