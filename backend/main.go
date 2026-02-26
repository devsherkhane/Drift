package main

import (
	"log"

	"github.com/devsherkhane/trello-clone/internal/auth"
	"github.com/devsherkhane/trello-clone/internal/database"
	"github.com/devsherkhane/trello-clone/internal/handlers"
	"github.com/devsherkhane/trello-clone/internal/logger"
	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/devsherkhane/trello-clone/internal/notifications"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Trello Clone API
// @version 1.0
// @description Backend API for Trello Clone with Boards, Lists, Cards, and Real-time updates.
// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	logger.SetupLogging()
	// 2. Initialize the database connection
	database.InitDB()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/uploads", "./uploads")

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
	apiPublic := r.Group("/api")
	{
		apiPublic.POST("/register", handlers.Register)
		apiPublic.POST("/login", handlers.Login)
	}

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
		api.GET("/search", handlers.SearchCards)
		api.GET("/cards/:id/attachments", handlers.GetAttachmentsByCard)

		api.PUT("/profile", handlers.UpdateProfile)

		api.POST("/attachments", handlers.UploadAttachment)
		api.GET("/boards/:id/activity", handlers.GetActivityLogs)

		api.GET("/ws", notifications.GlobalHub.HandleWS)
		api.GET("/boards/:id/export", handlers.ExportBoardCSV)
		api.PUT("/user/theme", handlers.UpdateTheme)

		api.POST("/boards/:id/collaborators", handlers.AddCollaborator)
		api.PATCH("/boards/:id/archive", handlers.ArchiveBoard)
		api.POST("/comments", handlers.CreateComment)
		api.GET("/cards/:id/comments", handlers.GetCommentsByCard)

		api.POST("/2fa/setup", handlers.Setup2FA)
		api.POST("/cards/labels", handlers.AddLabelToCard)

		api.GET("/search/advanced", handlers.AdvancedSearch)
		api.POST("/addlabel", handlers.AddLabelToCard)

		api.POST("/user/avatar", handlers.UploadAvatar)

		r.POST("/api/forgot-password", handlers.ForgotPassword)
		r.POST("/api/reset-password", handlers.ResetPassword)
	}

	// 4. Run the server on port 8080
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
