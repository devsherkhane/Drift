package service_test

import (
	"errors"
	"testing"

	"github.com/devsherkhane/trello-clone/internal/mocks"
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/devsherkhane/trello-clone/internal/repository"
	"github.com/devsherkhane/trello-clone/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestBoardService_CreateBoard(t *testing.T) {
	mockRepo := new(mocks.MockBoardRepository)
	boardSvc := service.NewBoardService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Create", "Personal Project", 42).Return(int64(10), nil).Once()

		board, err := boardSvc.CreateBoard("Personal Project", 42)

		assert.NoError(t, err)
		assert.NotNil(t, board)
		assert.Equal(t, 10, board.ID)
		assert.Equal(t, "Personal Project", board.Title)
		assert.Equal(t, 42, board.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.On("Create", "Error Project", 42).Return(int64(0), errors.New("db disconnect")).Once()

		board, err := boardSvc.CreateBoard("Error Project", 42)

		assert.Error(t, err)
		assert.Nil(t, board)
		assert.Equal(t, "db disconnect", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestBoardService_GetBoards(t *testing.T) {
	mockRepo := new(mocks.MockBoardRepository)
	boardSvc := service.NewBoardService(mockRepo)

	t.Run("Success Returns List", func(t *testing.T) {
		expectedBoards := []models.Board{
			{ID: 1, Title: "Board A", UserID: 42},
			{ID: 2, Title: "Board B", UserID: 42},
		}
		mockRepo.On("GetByUserID", 42).Return(expectedBoards, nil).Once()

		boards, err := boardSvc.GetBoards(42)

		assert.NoError(t, err)
		assert.Len(t, boards, 2)
		assert.Equal(t, "Board A", boards[0].Title)
		mockRepo.AssertExpectations(t)
	})
}

func TestBoardService_AddCollaborator(t *testing.T) {
	mockRepo := new(mocks.MockBoardRepository)
	boardSvc := service.NewBoardService(mockRepo)

	t.Run("Success Authorized Owner", func(t *testing.T) {
		// Mock checking ownership
		board := &models.Board{ID: 1, UserID: 42}
		mockRepo.On("GetByID", 1, 42).Return(board, nil).Once()
		// Mock adding collb
		mockRepo.On("AddCollaborator", 1, "friend@example.com", "member").Return(nil).Once()

		err := boardSvc.AddCollaborator(1, 42, "friend@example.com", "member")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized User", func(t *testing.T) {
		// Board is owned by 99, but user 42 is trying to add collaborator
		board := &models.Board{ID: 1, UserID: 99}
		mockRepo.On("GetByID", 1, 42).Return(board, nil).Once()

		err := boardSvc.AddCollaborator(1, 42, "friend@example.com", "member")
		
		assert.Error(t, err)
		assert.Equal(t, repository.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
