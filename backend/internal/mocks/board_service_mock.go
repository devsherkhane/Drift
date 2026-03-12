package mocks

import (
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockBoardService struct {
	mock.Mock
}

func (m *MockBoardService) CreateBoard(title string, userID int) (*models.Board, error) {
	args := m.Called(title, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Board), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBoardService) GetBoards(userID int) ([]models.Board, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Board), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBoardService) GetBoardByID(boardID, userID int) (*models.Board, error) {
	args := m.Called(boardID, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Board), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBoardService) UpdateBoardTitle(boardID, userID int, title string) error {
	args := m.Called(boardID, userID, title)
	return args.Error(0)
}

func (m *MockBoardService) DeleteBoard(boardID, userID int) error {
	args := m.Called(boardID, userID)
	return args.Error(0)
}

func (m *MockBoardService) ArchiveBoard(boardID, userID int) error {
	args := m.Called(boardID, userID)
	return args.Error(0)
}

func (m *MockBoardService) AddCollaborator(boardID, ownerID int, email, role string) error {
	args := m.Called(boardID, ownerID, email, role)
	return args.Error(0)
}

func (m *MockBoardService) RemoveCollaborator(boardID, ownerID, targetUserID int) error {
	args := m.Called(boardID, ownerID, targetUserID)
	return args.Error(0)
}

func (m *MockBoardService) GetCollaborators(boardID int) ([]models.User, error) {
	args := m.Called(boardID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBoardService) UpdateCollaboratorRole(boardID, userID, collaboratorID int, role string) error {
	args := m.Called(boardID, userID, collaboratorID, role)
	return args.Error(0)
}

func (m *MockBoardService) GetActivityLogs(boardID int) ([]models.ActivityLog, error) {
	args := m.Called(boardID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.ActivityLog), args.Error(1)
	}
	return nil, args.Error(1)
}
