package handler

import (
	"context"
	"fmt"
	"net/http"
	"study-planner-api/internal/api"
	"study-planner-api/internal/auth/provider"
	"study-planner-api/internal/database"
)

// Check if Application fully implements StrictServerInterface
var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	Test string
	DB   *database.Database
}

// GetAuthGoogleAuthorize implements api.StrictServerInterface.
func (s *Handler) GetAuthGoogleAuthorize(ctx context.Context, request api.GetAuthGoogleAuthorizeRequestObject) (api.GetAuthGoogleAuthorizeResponseObject, error) {

	reUrl := provider.GoogleAuthEndpoint("", "test-csrf")

	cookie := http.Cookie{
		Name:  "google_auth",
		Value: "test-csrf",
		// HttpOnly: true,
		// Secure: true,
	}

	return api.GetAuthGoogleAuthorize303Response{
		Headers: api.GetAuthGoogleAuthorize303ResponseHeaders{
			Location:  reUrl,
			SetCookie: cookie.String(),
		},
	}, nil

}

// PostAuthGoogleCallback implements api.StrictServerInterface.
func (s *Handler) GetAuthGoogleCallback(ctx context.Context, request api.GetAuthGoogleCallbackRequestObject) (api.GetAuthGoogleCallbackResponseObject, error) {
	hostUrl := api.HostUrlOfRequest(ctx)

	token := "test-token"

	// return api.GetAuthGoogleCallback200JSONResponse{}, nil

	return api.GetAuthGoogleCallback301Response{
		Headers: api.GetAuthGoogleCallback301ResponseHeaders{
			Location: fmt.Sprintf("%s#access_token=%s", hostUrl, token),
		},
	}, nil
}

func NewHandler() *Handler {
	return &Handler{
		Test: "Hello World",
		DB:   database.Instance(),
	}
}
