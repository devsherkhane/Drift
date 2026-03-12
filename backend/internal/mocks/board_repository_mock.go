package mocks

import (
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockBoardRepository struct {
	mock.Mock
}

func (m *MockBoardRepository) Create(title string, userID int) (int64, error) {
	args := m.Called(title, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockBoardRepository) GetByUserID(userID int) ([]models.Board, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Board), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBoardRepository) GetByID(boardID, userID int) (*models.Board, error) {
	args := m.Called(boardID, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Board), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBoardRepository) UpdateTitle(boardID, userID int, title string) error {
	args := m.Called(boardID, userID, title)
	return args.Error(0)
}

func (m *MockBoardRepository) Delete(boardID, userID int) error {
	args := m.Called(boardID, userID)
	return args.Error(0)
}

func (m *MockBoardRepository) Archive(boardID, userID int) error {
	args := m.Called(boardID, userID)
	return args.Error(0)
}

func (m *MockBoardRepository) AddCollaborator(boardID int, email, role string) error {
	args := m.Called(boardID, email, role)
	return args.Error(0)
}

func (m *MockBoardRepository) GetCollaborators(boardID int) ([]models.User, error) {
	args := m.Called(boardID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBoardRepository) UpdateCollaboratorRole(boardID, collaboratorID int, role string) error {
	args := m.Called(boardID, collaboratorID, role)
	return args.Error(0)
}

func (m *MockBoardRepository) RemoveCollaborator(boardID, collaboratorID int) error {
	args := m.Called(boardID, collaboratorID)
	return args.Error(0)
}

func (m *MockBoardRepository) GetActivityLogs(boardID int) ([]models.ActivityLog, error) {
	args := m.Called(boardID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.ActivityLog), args.Error(1)
	}
	return nil, args.Error(1)
}
