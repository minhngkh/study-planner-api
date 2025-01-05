package handler

import (
	"context"
	"errors"
	"study-planner-api/internal/api"
	"study-planner-api/internal/auth"
	"study-planner-api/internal/utils"

	"github.com/rs/zerolog/log"
)

func (s *Handler) PostAuthPasswordReset(
	ctx context.Context,
	request api.PostAuthPasswordResetRequestObject,
) (api.PostAuthPasswordResetResponseObject, error) {
	email := request.Body.Email
	if err := s.Validate.Var(email, "required,email"); err != nil {
		return api.PostAuthPasswordReset400Response{}, err
	}

	err := auth.SendPasswordResetEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUnknownEmail):
			// Still return 200 to prevent email enumeration
			return api.PostAuthPasswordReset200Response{}, nil
		default:
			return nil, err
		}
	}

	return api.PostAuthPasswordReset200Response{}, nil
}

func (s *Handler) PostAuthPasswordResetConfirm(
	ctx context.Context,
	request api.PostAuthPasswordResetConfirmRequestObject,
) (api.PostAuthPasswordResetConfirmResponseObject, error) {
	userId := request.Body.UserId
	token := request.Body.Token
	password := request.Body.NewPassword

	if err := s.Validate.Var(password, "required,default-password"); err != nil {
		log.Error().Err(err).Msg("Password validation failed")
		return api.PostAuthPasswordResetConfirm400Response{}, nil
	}

	err := auth.ResetPassword(userId, token, password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidToken):
			return api.PostAuthPasswordResetConfirm403JSONResponse{
				TokenErrorJSONResponse: api.TokenErrorJSONResponse{
					Type: utils.Ptr(api.InvalidToken),
				},
			}, nil
		case errors.Is(err, auth.ErrExpiredToken):
			return api.PostAuthPasswordResetConfirm403JSONResponse{
				TokenErrorJSONResponse: api.TokenErrorJSONResponse{
					Type: utils.Ptr(api.ExpiredToken),
				},
			}, nil
		default:
			return nil, err
		}
	}

	return api.PostAuthPasswordResetConfirm200Response{}, nil
}
