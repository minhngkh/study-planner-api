package handler

import (
	"context"
	"errors"
	"study-planner-api/internal/api"
	"study-planner-api/internal/auth"
	"study-planner-api/internal/user"
	"study-planner-api/internal/utils"

	"github.com/rs/zerolog/log"
)

func (s *Handler) PostRegister(
	ctx context.Context,
	request api.PostRegisterRequestObject,
) (api.PostRegisterResponseObject, error) {
	email, password := request.Body.Email, request.Body.Password

	userId, err := user.CreateUser(email, password)
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			return api.PostRegister400JSONResponse{
				Type:    utils.Ptr(api.DuplicateEmail),
				Message: utils.Ptr("User already exists"),
			}, nil
		}

		return nil, err
	}

	accessToken, refreshToken, err := auth.CreateAuthTokens(auth.AuthInfo{UserID: userId})
	if err != nil {
		return nil, err
	}

	err = auth.CreateSession(userId, refreshToken)
	if err != nil {
		return nil, err
	}

	return api.PostRegister201JSONResponse{
		AccessToken:  &accessToken.Value,
		RefreshToken: &refreshToken.Value,
	}, nil
}

func (s *Handler) PostLogin(
	ctx context.Context,
	request api.PostLoginRequestObject,
) (api.PostLoginResponseObject, error) {
	loginInfo := auth.LoginInfo{
		Email:    *request.Body.Email,
		Password: *request.Body.Password,
	}

	user, err := auth.VerifyLoginInfo(loginInfo)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			fallthrough
		case errors.Is(err, auth.ErrIncorrectPassword):
			return api.PostLogin400Response{}, nil
		default:
			log.Debug().Err(err).Msg("incorrect password")
			return nil, err
		}
	}

	accessToken, refreshToken, err := auth.CreateAuthTokens(auth.AuthInfo{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	err = auth.CreateSession(user.ID, refreshToken)
	if err != nil {
		return nil, err
	}

	return api.PostLogin200JSONResponse{
		AccessToken:  &accessToken.Value,
		RefreshToken: &refreshToken.Value,
	}, nil
}

func (s *Handler) PostAuthRefreshToken(
	ctx context.Context,
	request api.PostAuthRefreshTokenRequestObject,
) (api.PostAuthRefreshTokenResponseObject, error) {
	var refreshToken string
	if request.Body.RefreshToken != nil {
		log.Debug().Str("refreshToken", refreshToken).Msg("refresh token")
		refreshToken = *request.Body.RefreshToken
	} else if request.Params.RefreshToken != nil {
		log.Debug().Str("refreshToken", refreshToken).Msg("refresh token")

		refreshToken = *request.Params.RefreshToken
	} else {
		log.Debug().Str("refreshToken", refreshToken).Msg("refresh token")

		return api.PostAuthRefreshToken401Response{}, nil
	}

	info, _, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		return api.PostAuthRefreshToken401Response{}, nil
	}

	newAccessToken, newRefreshToken, err := auth.CreateAuthTokens(auth.RefreshInfo{
		UserID: info.UserID,
	})
	if err != nil {
		return nil, err
	}

	err = auth.UpdateSession(info.UserID, refreshToken, newRefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrMaliciousRefreshToken) {
			return api.PostAuthRefreshToken403Response{}, nil
		}

		return nil, err
	}

	return api.PostAuthRefreshToken200JSONResponse{
		AccessToken:  &newAccessToken.Value,
		RefreshToken: &newRefreshToken.Value,
	}, nil
}

func (s *Handler) PostLogout(
	ctx context.Context,
	request api.PostLogoutRequestObject,
) (api.PostLogoutResponseObject, error) {
	var refreshToken string
	if request.Body.RefreshToken != nil {
		refreshToken = *request.Body.RefreshToken
	} else {
		refreshToken = *request.Params.RefreshToken
	}

	info, _, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		return api.PostLogout401Response{}, nil
	}

	err = auth.RemoveSession(info.UserID, refreshToken)
	if err != nil {
		return api.PostLogout401Response{}, nil
	}

	return api.PostLogout200Response{}, nil
}
