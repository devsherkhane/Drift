package mocks

import (
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockListService struct {
	mock.Mock
}

func (m *MockListService) CreateList(boardID int, title string) (*models.List, error) {
	args := m.Called(boardID, title)
	if args.Get(0) != nil {
		return args.Get(0).(*models.List), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockListService) GetListsByBoard(boardID int) ([]models.List, error) {
	args := m.Called(boardID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.List), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockListService) UpdateListTitle(listID int, title string) error {
	args := m.Called(listID, title)
	return args.Error(0)
}

func (m *MockListService) DeleteList(listID int) error {
	args := m.Called(listID)
	return args.Error(0)
}
