package service_test

import (
	"errors"
	"testing"

	"github.com/devsherkhane/trello-clone/internal/mocks"
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/devsherkhane/trello-clone/internal/service"
	"github.com/devsherkhane/trello-clone/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authSvc := service.NewAuthService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByEmail", "new@example.com").Return(nil, errors.New("not found")).Once()
		mockRepo.On("Create", "newuser", "new@example.com", mock.Anything).Return(int64(1), nil).Once()

		user, err := authSvc.Register("newuser", "new@example.com", "password123")
		
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "newuser", user.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		existingUser := &models.User{Email: "taken@example.com"}
		mockRepo.On("GetByEmail", "taken@example.com").Return(existingUser, nil).Once()

		user, err := authSvc.Register("newuser", "taken@example.com", "password123")
		
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "email already registered", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authSvc := service.NewAuthService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		hashedPass, _ := utils.HashPassword("password123")
		dbUser := &models.User{ID: 1, Email: "test@example.com", PasswordHash: hashedPass}
		
		mockRepo.On("GetByEmail", "test@example.com").Return(dbUser, nil).Once()

		token, user, err := authSvc.Login("test@example.com", "password123")
		
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		// Ensure password hash is cleared before returning
		assert.Empty(t, user.PasswordHash) 
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		hashedPass, _ := utils.HashPassword("password123")
		dbUser := &models.User{ID: 1, Email: "test@example.com", PasswordHash: hashedPass}
		
		mockRepo.On("GetByEmail", "test@example.com").Return(dbUser, nil).Once()

		token, user, err := authSvc.Login("test@example.com", "wrongpass")
		
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
