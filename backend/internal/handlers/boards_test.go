package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devsherkhane/trello-clone/internal/handlers"
	"github.com/devsherkhane/trello-clone/internal/mocks"
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Fake auth middleware just to attach a userID to the context for protected routes
func fakeAuthMiddleware(userID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

func setupBoardsRouter(userID int) (*gin.Engine, *mocks.MockBoardService, *mocks.MockListService) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockBoardService := new(mocks.MockBoardService)
	// We might need list service for GetBoard
	mockListService := new(mocks.MockListService)

	apiHandler := handlers.NewAPIHandler(nil, mockBoardService, mockListService, nil, nil, nil, nil, nil)
	
	// Protected routes group
	api := r.Group("/api")
	api.Use(fakeAuthMiddleware(userID))
	{
		api.POST("/boards", apiHandler.CreateBoard)
		api.GET("/boards", apiHandler.GetBoards)
	}

	return r, mockBoardService, mockListService
}

// We need a dummy mock for ListService since GetBoard uses it. Let's create it inline or just map it in our setup.
// Wait, for GetBoards and CreateBoard, we don't need ListService. It's safe to use a dummy struct here.
func TestCreateBoardHandler(t *testing.T) {
	userID := 42
	r, mockBoard, _ := setupBoardsRouter(userID)

	t.Run("Valid Board Creation", func(t *testing.T) {
		inputData := map[string]string{
			"title": "New Project",
		}
		body, _ := json.Marshal(inputData)

		expectedBoard := &models.Board{ID: 100, Title: "New Project", UserID: userID}
		mockBoard.On("CreateBoard", "New Project", userID).Return(expectedBoard, nil)

		req, _ := http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "New Project")
		mockBoard.AssertExpectations(t)
	})

	t.Run("Missing Title", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/boards", bytes.NewBuffer([]byte("{}")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "required")
	})

	t.Run("Database Error", func(t *testing.T) {
		inputData := map[string]string{
			"title": "Error Project",
		}
		body, _ := json.Marshal(inputData)

		mockBoard.On("CreateBoard", "Error Project", userID).Return(nil, errors.New("db error"))

		req, _ := http.NewRequest("POST", "/api/boards", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockBoard.AssertExpectations(t)
	})
}
