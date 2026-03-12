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

func setupAuthRouter() (*gin.Engine, *mocks.MockAuthService) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockAuthService := new(mocks.MockAuthService)
	// Inject the mock
	apiHandler := handlers.NewAPIHandler(mockAuthService, nil, nil, nil, nil, nil, nil, nil)
	
	r.POST("/login", apiHandler.Login)
	r.POST("/register", apiHandler.Register)

	return r, mockAuthService
}

func TestLoginHandler(t *testing.T) {
	r, mockAuth := setupAuthRouter()

	t.Run("Empty Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("{}")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "required")
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "not-an-email",
			"password": "password123",
		}
		body, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Valid Login", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(loginData)

		expectedUser := &models.User{ID: 1, Email: "test@example.com"}
		mockAuth.On("Login", "test@example.com", "password123").Return("mockToken123", expectedUser, nil)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "mockToken123")
		mockAuth.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "wrong@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(loginData)

		mockAuth.On("Login", "wrong@example.com", "password123").Return("", nil, errors.New("invalid credentials"))

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})
}

func TestRegisterHandler(t *testing.T) {
	r, mockAuth := setupAuthRouter()

	t.Run("Valid Registration", func(t *testing.T) {
		regData := map[string]string{
			"username": "newuser",
			"email":    "new@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(regData)

		expectedUser := &models.User{ID: 1, Username: "newuser", Email: "new@example.com"}
		mockAuth.On("Register", "newuser", "new@example.com", "password123").Return(expectedUser, nil)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "newuser")
		mockAuth.AssertExpectations(t)
	})

	t.Run("Email Already Taken", func(t *testing.T) {
		regData := map[string]string{
			"username": "takenuser",
			"email":    "taken@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(regData)

		mockAuth.On("Register", "takenuser", "taken@example.com", "password123").Return(nil, errors.New("email already registered"))

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "email already registered")
	})
}
