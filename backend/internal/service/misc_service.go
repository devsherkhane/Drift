package service

import (
	"github.com/devsherkhane/trello-clone/internal/models"
	"github.com/devsherkhane/trello-clone/internal/repository"
)

// Comment Service
type CommentService interface {
	CreateComment(cardID, userID int, text string) (*models.Comment, error)
	GetCommentsByCard(cardID int) ([]models.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(r repository.CommentRepository) CommentService {
	return &commentService{repo: r}
}

func (s *commentService) CreateComment(cardID, userID int, text string) (*models.Comment, error) {
	id, err := s.repo.Create(cardID, userID, text)
	if err != nil {
		return nil, err
	}
	return &models.Comment{ID: int(id), CardID: cardID, UserID: userID, Text: text}, nil
}

func (s *commentService) GetCommentsByCard(cardID int) ([]models.Comment, error) {
	return s.repo.GetByCardID(cardID)
}

// Attachment Service
type AttachmentService interface {
	AddAttachment(cardID, userID int, filename, filePath string) (*models.Attachment, error)
	GetAttachmentsByCard(cardID int) ([]models.Attachment, error)
}

type attachmentService struct {
	repo repository.AttachmentRepository
}

func NewAttachmentService(r repository.AttachmentRepository) AttachmentService {
	return &attachmentService{repo: r}
}

func (s *attachmentService) AddAttachment(cardID, userID int, filename, filePath string) (*models.Attachment, error) {
	id, err := s.repo.Create(cardID, userID, filename, filePath)
	if err != nil {
		return nil, err
	}
	return &models.Attachment{ID: int(id), CardID: cardID, UserID: userID, Filename: filename, FilePath: filePath}, nil
}

func (s *attachmentService) GetAttachmentsByCard(cardID int) ([]models.Attachment, error) {
	return s.repo.GetByCardID(cardID)
}

// Label Service
type LabelService interface {
	AddLabel(cardID, labelID int) error
}

type labelService struct {
	repo repository.LabelRepository
}

func NewLabelService(r repository.LabelRepository) LabelService {
	return &labelService{repo: r}
}

func (s *labelService) AddLabel(cardID, labelID int) error {
	return s.repo.AddLabelToCard(cardID, labelID)
}

// Search Service
type SearchService interface {
	Search(query string, boardID *int, userID int) ([]map[string]interface{}, error)
}

type searchService struct {
	repo repository.SearchRepository
}

func NewSearchService(r repository.SearchRepository) SearchService {
	return &searchService{repo: r}
}

func (s *searchService) Search(query string, boardID *int, userID int) ([]map[string]interface{}, error) {
	return s.repo.AdvancedSearch(query, boardID, userID)
}
