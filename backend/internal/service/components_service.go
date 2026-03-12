package service

import (
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/devsherkhane/trello-clone/internal/repository"
)

type ComponentService interface {
    GetLists(boardID int) ([]models.List, error)
    GetCards(listID int) ([]models.Card, error)
}

type ListService interface {
	CreateList(boardID int, title string) (*models.List, error)
	GetListsByBoard(boardID int) ([]models.List, error)
	UpdateListTitle(listID int, title string) error
	DeleteList(listID int) error
}

type listService struct {
	listRepo repository.ListRepository
}

func NewListService(repo repository.ListRepository) ListService {
	return &listService{listRepo: repo}
}

func (s *listService) CreateList(boardID int, title string) (*models.List, error) {
	id, err := s.listRepo.Create(boardID, title)
	if err != nil {
		return nil, err
	}
	return &models.List{ID: int(id), BoardID: boardID, Title: title}, nil
}

func (s *listService) GetListsByBoard(boardID int) ([]models.List, error) {
	return s.listRepo.GetByBoardID(boardID)
}

func (s *listService) UpdateListTitle(listID int, title string) error {
	return s.listRepo.UpdateTitle(listID, title)
}

func (s *listService) DeleteList(listID int) error {
	return s.listRepo.Delete(listID)
}

type CardService interface {
	CreateCard(listID int, title string) (*models.Card, error)
	GetCardsByList(listID int) ([]models.Card, error)
	GetCardByID(cardID int) (*models.Card, error)
	UpdateCard(card *models.Card) error
	MoveCard(cardID, newListID, newPosition int) error
	DeleteCard(cardID int) error
}

type cardService struct {
	cardRepo repository.CardRepository
}

func NewCardService(repo repository.CardRepository) CardService {
	return &cardService{cardRepo: repo}
}

func (s *cardService) CreateCard(listID int, title string) (*models.Card, error) {
	id, err := s.cardRepo.Create(listID, title)
	if err != nil {
		return nil, err
	}
	return s.cardRepo.GetByID(int(id))
}

func (s *cardService) GetCardsByList(listID int) ([]models.Card, error) {
	return s.cardRepo.GetByListID(listID)
}

func (s *cardService) GetCardByID(cardID int) (*models.Card, error) {
	return s.cardRepo.GetByID(cardID)
}

func (s *cardService) UpdateCard(card *models.Card) error {
	return s.cardRepo.Update(card)
}

func (s *cardService) MoveCard(cardID, newListID, newPosition int) error {
	return s.cardRepo.Move(cardID, newListID, newPosition)
}

func (s *cardService) DeleteCard(cardID int) error {
	return s.cardRepo.Delete(cardID)
}
