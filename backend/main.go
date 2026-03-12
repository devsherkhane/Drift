package main

import (
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/devsherkhane/trello-clone/internal/auth"
	"github.com/devsherkhane/trello-clone/internal/database"
	"github.com/devsherkhane/trello-clone/internal/handlers"
	"github.com/devsherkhane/trello-clone/internal/logger"
	"github.com/devsherkhane/trello-clone/internal/middleware"
	"github.com/devsherkhane/trello-clone/internal/notifications"
	"github.com/devsherkhane/trello-clone/internal/repository"
	"github.com/devsherkhane/trello-clone/internal/service"
	"github.com/gin-contrib/cors"
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
		log.Println("No .env file found (this is fine in production using injected env variables)")
	}
	logger.SetupLogging()
	defer logger.Sync()

	logger.Log.Info("Starting Trello Clone API server...")
	
	// Initialize database
	database.InitDB()
	db := database.DB

	// DI: Repositories
	userRepo := repository.NewUserRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	listRepo := repository.NewListRepository(db)
	cardRepo := repository.NewCardRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	labelRepo := repository.NewLabelRepository(db)
	searchRepo := repository.NewSearchRepository(db)

	// DI: Services
	authSvc := service.NewAuthService(userRepo)
	boardSvc := service.NewBoardService(boardRepo)
	listSvc := service.NewListService(listRepo)
	cardSvc := service.NewCardService(cardRepo)
	commentSvc := service.NewCommentService(commentRepo)
	attachmentSvc := service.NewAttachmentService(attachmentRepo)
	labelSvc := service.NewLabelService(labelRepo)
	searchSvc := service.NewSearchService(searchRepo)

	// DI: Handlers
	apiHandler := handlers.NewAPIHandler(
		authSvc, boardSvc, listSvc, cardSvc, commentSvc, attachmentSvc, labelSvc, searchSvc,
	)

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.StructuredLogger())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimiter(10, 20))
	r.Use(middleware.ErrorHandler())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Static("/uploads", "./uploads")

	// System Routes
	r.GET("/health", handlers.HealthCheck)

	// Public Routes
	apiPublic := r.Group("/api")
	apiPublic.Use(middleware.RateLimiter(3, 5))
	{
		apiPublic.POST("/register", apiHandler.Register)
		apiPublic.POST("/login", apiHandler.Login)
		apiPublic.POST("/forgot-password", apiHandler.ForgotPassword)
		apiPublic.POST("/reset-password", apiHandler.ResetPassword)
	}

	// Protected Routes (Require JWT)
	api := r.Group("/api")
	api.Use(auth.AuthMiddleware())
	{
		// Board Routes
		api.POST("/boards", apiHandler.CreateBoard)
		api.GET("/boards", apiHandler.GetBoards)
		api.PUT("/boards/:id", apiHandler.UpdateBoard)
		api.DELETE("/boards/:id", apiHandler.DeleteBoard)
		api.GET("/boards/:id", apiHandler.GetBoard)
		api.GET("/boards/:id/activity", apiHandler.GetActivityLogs)
		api.GET("/boards/:id/export", apiHandler.ExportBoardCSV)
		api.POST("/boards/:id/collaborators", apiHandler.AddCollaborator)
		api.PATCH("/boards/:id/invitation", apiHandler.RespondToInvitation)
		api.PATCH("/boards/:id/archive", apiHandler.ArchiveBoard)

		// List Routes
		api.POST("/lists", apiHandler.CreateList)
		api.GET("/boards/:id/lists", apiHandler.GetLists) // Was GetListsByBoard
		api.PUT("/lists/:id", apiHandler.UpdateList)
		api.DELETE("/lists/:id", apiHandler.DeleteList)

		// Card Routes
		api.POST("/cards", apiHandler.CreateCard)
		api.GET("/lists/:id/cards", apiHandler.GetCards) // Was GetCardsByList
		api.PUT("/cards/:id", apiHandler.UpdateCard)
		api.DELETE("/cards/:id", apiHandler.DeleteCard)
		api.PATCH("/cards/:id/move", apiHandler.MoveCard)
		
		api.POST("/cards/labels", apiHandler.AddLabelToCard)
		api.POST("/addlabel", apiHandler.AddLabelToCard)

		// Attachments & Comments
		api.POST("/comments", apiHandler.CreateComment)
		api.GET("/cards/:id/comments", apiHandler.GetComments)
		api.POST("/attachments", apiHandler.UploadAttachment)
		api.GET("/cards/:id/attachments", apiHandler.GetAttachments)

		// Search
		api.GET("/search", apiHandler.SearchCards)
		api.GET("/search/advanced", apiHandler.SearchCards)

		// User
		api.GET("/user/me", apiHandler.GetProfile)
		api.PUT("/profile", apiHandler.UpdateProfile)
		api.PUT("/user/theme", apiHandler.UpdateTheme)
		api.POST("/user/avatar", apiHandler.UploadAvatar)

		// WebSocket
		api.GET("/ws", notifications.GlobalHub.HandleWS)

		api.GET("/ping", func(c *gin.Context) {
			userID := c.MustGet("userID")
			c.JSON(200, gin.H{"user_id": userID})
		})
	}

	logger.Log.Info("Server listening on port 8080")
	if err := r.Run(":8080"); err != nil {
		logger.Log.Fatal("Failed to run server", zap.Error(err))
	}
}
