package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	fmt.Printf("Connecting to DB: %s at %s\n", dbName, host)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true", username, password, host, dbName)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	// Run migrations
	driver, err := mysql.WithInstance(DB, &mysql.Config{})
	if err != nil {
		log.Fatal("Could not instantiate mysql driver for migrations:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("Could not instantiate migration tool:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error running migrations:", err)
	} else if err == migrate.ErrNoChange {
		log.Println("Database already up to date (no migrations to run).")
	} else {
		log.Println("Database migrations applied successfully!")
	}
}
