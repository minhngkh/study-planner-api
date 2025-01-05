package handler

import (
	"context"
	"study-planner-api/internal/api"
	"study-planner-api/internal/database"
)

// Check if Application fully implements StrictServerInterface
var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	Test string
	DB   *database.Database
}

// PostAuthPasswordReset implements api.StrictServerInterface.
func (s *Handler) PostAuthPasswordReset(ctx context.Context, request api.PostAuthPasswordResetRequestObject) (api.PostAuthPasswordResetResponseObject, error) {
	panic("unimplemented")
}

// PostAuthPasswordResetConfirm implements api.StrictServerInterface.
func (s *Handler) PostAuthPasswordResetConfirm(ctx context.Context, request api.PostAuthPasswordResetConfirmRequestObject) (api.PostAuthPasswordResetConfirmResponseObject, error) {
	panic("unimplemented")
}

func NewHandler() *Handler {
	return &Handler{
		Test: "Hello World",
		DB:   database.Instance(),
	}
}
