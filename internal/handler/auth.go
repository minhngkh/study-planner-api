package handler

import (
	"context"
	"errors"
	"net/http"
	"study-planner-api/internal/api"
	"study-planner-api/internal/auth"
	"study-planner-api/internal/auth/token"
	"study-planner-api/internal/user"
	"study-planner-api/internal/utils"
	"time"

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

	accessToken, refreshToken, err := token.CreateAuthTokens(token.AuthInfo{UserID: userId})
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

	accessToken, refreshToken, err := token.CreateAuthTokens(token.AuthInfo{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	err = auth.CreateSession(user.ID, refreshToken)
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Value,
		HttpOnly: true,
		Expires:  refreshToken.Expiry.Time(),
	}

	return api.PostLogin200JSONResponse{
		Headers: api.PostLogin200ResponseHeaders{
			SetCookie: cookie.String(),
		},
		Body: api.AuthTokens{
			AccessToken:  &accessToken.Value,
			RefreshToken: &refreshToken.Value,
		},
	}, nil
}

func (s *Handler) PostAuthRefreshToken(
	ctx context.Context,
	request api.PostAuthRefreshTokenRequestObject,
) (api.PostAuthRefreshTokenResponseObject, error) {
	var refreshToken string
	if request.Body.RefreshToken != nil {
		refreshToken = *request.Body.RefreshToken
	} else if request.Params.RefreshToken != nil {
		refreshToken = *request.Params.RefreshToken
	} else {
		log.Debug().Msg("no refresh token provided")
		return api.PostAuthRefreshToken403Response{}, nil
	}

	info, _, err := token.ValidateRefreshToken(refreshToken)
	if err != nil {
		return api.PostAuthRefreshToken403Response{}, nil
	}

	newAccessToken, newRefreshToken, err := token.CreateAuthTokens(token.AuthInfo{
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

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken.Value,
		HttpOnly: true,
	}

	return api.PostAuthRefreshToken200JSONResponse{
		Headers: api.PostAuthRefreshToken200ResponseHeaders{
			SetCookie: cookie.String(),
		},
		Body: api.AuthTokens{
			AccessToken:  &newAccessToken.Value,
			RefreshToken: &newRefreshToken.Value,
		},
	}, nil
}

func (s *Handler) PostLogout(
	ctx context.Context,
	request api.PostLogoutRequestObject,
) (api.PostLogoutResponseObject, error) {
	deletedCookie := http.Cookie{
		Name:    "refresh_token",
		Value:   "deleted",
		Expires: time.Unix(0, 0),
	}

	var refreshToken string
	if request.Body.RefreshToken != nil {
		refreshToken = *request.Body.RefreshToken
	} else if request.Params.RefreshToken != nil {
		refreshToken = *request.Params.RefreshToken
	} else {
		log.Debug().Msg("no refresh token provided")
		return api.PostLogout403Response{
			Headers: api.PostLogout403ResponseHeaders{
				SetCookie: deletedCookie.String(),
			},
		}, nil
	}

	info, _, err := token.ValidateRefreshToken(refreshToken)
	if err != nil {
		return api.PostLogout403Response{
			Headers: api.PostLogout403ResponseHeaders{
				SetCookie: deletedCookie.String(),
			},
		}, nil
	}

	err = auth.RemoveSession(info.UserID, refreshToken)
	if err != nil {
		return api.PostLogout403Response{
			Headers: api.PostLogout403ResponseHeaders{
				SetCookie: deletedCookie.String(),
			},
		}, nil
	}

	return api.PostLogout200Response{
		Headers: api.PostLogout200ResponseHeaders{
			SetCookie: deletedCookie.String(),
		},
	}, nil
}
