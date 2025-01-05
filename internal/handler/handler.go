package handler

import (
	"study-planner-api/internal/api"
	"study-planner-api/internal/database"
	"study-planner-api/internal/validator"
)

// Check if Application fully implements StrictServerInterface
var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	Test     string
	DB       *database.Database
	Validate *validator.Validate
}

func NewHandler() *Handler {
	return &Handler{
		Test:     "Hello World",
		DB:       database.Instance(),
		Validate: validator.Instance(),
	}
}
