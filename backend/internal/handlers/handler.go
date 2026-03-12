package handlers

import (
	"github.com/devsherkhane/trello-clone/internal/service"
)

type APIHandler struct {
	AuthService       service.AuthService
	BoardService      service.BoardService
	ListService       service.ListService
	CardService       service.CardService
	CommentService    service.CommentService
	AttachmentService service.AttachmentService
	LabelService      service.LabelService
	SearchService     service.SearchService
}

func NewAPIHandler(
	auth service.AuthService,
	board service.BoardService,
	list service.ListService,
	card service.CardService,
	comment service.CommentService,
	attachment service.AttachmentService,
	label service.LabelService,
	search service.SearchService,
) *APIHandler {
	return &APIHandler{
		AuthService:       auth,
		BoardService:      board,
		ListService:       list,
		CardService:       card,
		CommentService:    comment,
		AttachmentService: attachment,
		LabelService:      label,
		SearchService:     search,
	}
}
