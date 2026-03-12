package handlers

import (
	"net/http"
	"time"

	"github.com/devsherkhane/trello-clone/internal/database"
	"github.com/devsherkhane/trello-clone/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthCheck verifies the application and its dependencies are running
func HealthCheck(c *gin.Context) {
	// Check Database Connection
	dbStatus := "connected"
	if database.DB == nil {
		dbStatus = "disconnected"
	} else if err := database.DB.Ping(); err != nil {
		dbStatus = "disconnected"
		logger.Log.Error("Healthcheck failed to ping database", zap.Error(err))
	}

	// Overall System Status
	systemStatus := "UP"
	if dbStatus != "connected" {
		systemStatus = "DOWN"
	}

	// Payload
	response := gin.H{
		"status":   systemStatus,
		"database": dbStatus,
		"time":     time.Now().Format(time.RFC3339),
	}

	if systemStatus == "UP" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}
