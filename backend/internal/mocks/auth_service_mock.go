package mocks

import (
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(username, email, password string) (*models.User, error) {
	args := m.Called(username, email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) Login(email, password string) (string, *models.User, error) {
	args := m.Called(email, password)
	if args.Get(1) != nil {
		return args.String(0), args.Get(1).(*models.User), args.Error(2)
	}
	return args.String(0), nil, args.Error(2)
}

func (m *MockAuthService) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(token, newPassword string) error {
	args := m.Called(token, newPassword)
	return args.Error(0)
}

func (m *MockAuthService) UpdateProfile(userID int, username, newEmail string) (*models.User, error) {
	args := m.Called(userID, username, newEmail)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}
