package main

import (
	"log"

	"github.com/devsherkhane/trello-clone/internal/auth"
	"github.com/devsherkhane/trello-clone/internal/database"
	"github.com/devsherkhane/trello-clone/internal/handlers"
	"github.com/devsherkhane/trello-clone/internal/logger"
	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	logger.SetupLogging()
	// 2. Initialize the database connection
	database.InitDB()

	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	// 3. CORS Middleware to allow requests from your frontend
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public Routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected Routes (Require JWT)
	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	{
		// Board Routes
		api.POST("/boards", handlers.CreateBoard)
		api.GET("/boards", handlers.GetBoards)
		api.PUT("/boards/:id", handlers.UpdateBoard)
		api.DELETE("/boards/:id", handlers.DeleteBoard)
		api.GET("/boards/:id", handlers.GetBoard)

		// List Routes
		api.POST("/lists", handlers.CreateList)
		api.GET("/boards/:id/lists", handlers.GetListsByBoard)
		api.PUT("/lists/:id", handlers.UpdateList)
		api.DELETE("/lists/:id", handlers.DeleteList)

		api.GET("/ping", func(c *gin.Context) {
			userID := c.MustGet("userID")
			c.JSON(200, gin.H{"user_id": userID})
		})

		api.POST("/cards", handlers.CreateCard)
		api.GET("/lists/:id/cards", handlers.GetCardsByList)
		api.PUT("/cards/:id", handlers.UpdateCard)
		api.DELETE("/cards/:id", handlers.DeleteCard)
		api.PATCH("/cards/:id/move", handlers.MoveCard)

	}

	// 4. Run the server on port 8080
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
