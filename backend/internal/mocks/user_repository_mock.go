package mocks

import (
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(username, email, passwordHash string) (int64, error) {
	args := m.Called(username, email, passwordHash)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByResetToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}
