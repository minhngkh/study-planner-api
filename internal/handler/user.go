package handler

import (
	"context"
	"errors"
	"study-planner-api/internal/api"
	"study-planner-api/internal/user"
)

func (s *Handler) GetProfile(
	ctx context.Context,
	request api.GetProfileRequestObject,
) (api.GetProfileResponseObject, error) {
	authInfo, ok := ctx.Value(api.ContextKey("auth")).(api.AuthInfo)
	if !ok {
		return nil, errors.New("auth info not found")
	}

	userInfo, err := user.GetUserInfo(authInfo.ID)
	if err != nil {
		return nil, err
	}

	return api.GetProfile200JSONResponse{
		CreatedAt: &userInfo.CreatedAt,
		Email:     &userInfo.Email,
		ID:        &userInfo.ID,
	}, nil
}
