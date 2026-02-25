package logger

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupLogging() {
	// Check the environment variable
	env := os.Getenv("APP_ENV")

	if env == "production" {
		// In production, log to a file
		f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Could not create log file")
		}
		// Set Gin to write logs to the file
		gin.DefaultWriter = io.MultiWriter(f)
		// Set standard log to write to the file
		log.SetOutput(f)
	} else {
		// In development, log to the console (default behavior)
		gin.DefaultWriter = os.Stdout
		log.SetOutput(os.Stdout)
	}
}