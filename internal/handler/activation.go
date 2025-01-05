package handler

import (
	"context"
	"errors"
	"study-planner-api/internal/api"
	"study-planner-api/internal/user"
	"study-planner-api/internal/utils"
)

// PostActivation implements api.StrictServerInterface.
func (s *Handler) PostActivation(
	ctx context.Context,
	request api.PostActivationRequestObject,
) (api.PostActivationResponseObject, error) {
	err := user.ActivateAccount(request.Body.UserId, request.Body.Token)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrExpiredToken):
			return api.PostActivation400JSONResponse{
				Type: utils.Ptr(api.ActivationErrorTypeExpiredToken),
			}, err
		case errors.Is(err, user.ErrInvalidToken):
			return api.PostActivation400JSONResponse{
				Type: utils.Ptr(api.ActivationErrorTypeInvalidToken),
			}, err
		default:
			return nil, err
		}
	}

	return api.PostActivation200Response{}, nil
}

// PostActivationEmail implements api.StrictServerInterface.
func (s *Handler) PostActivationEmail(
	ctx context.Context,
	request api.PostActivationEmailRequestObject,
) (api.PostActivationEmailResponseObject, error) {
	authInfo := api.AuthInfoOfRequest(ctx)

	err := user.SendActivationEmail(authInfo.ID)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyActivated) {
			return api.PostActivationEmail400Response{}, err
		}

		return nil, err
	}

	return api.PostActivationEmail200Response{}, nil
}
