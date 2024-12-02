package auth

import (
	"net/http"
	"study-planner-api/internal/validator"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type registerRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required,min=6"`
}

func RegisterHandler(c echo.Context) error {
	req := new(registerRequest)
	if httpErr := validator.BindAndValidateRequest(c, req); httpErr != nil {
		return httpErr
	}

	userId, err := CreateUser(req.Email, req.Password)
	if err != nil {
		log.Debug().Msg(err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user")
	}

	tokens, err := NewJwtAuthTokens(AccessTokenCustomClaims{
		UserID: userId,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = CreateSession(userId, tokens.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	SetJwtTokenInCookie(c, tokens)

	return c.JSON(http.StatusCreated, map[string]string{
		"message":       "User successfully registered",
		"access_token":  tokens.AccessToken.Value,
		"refresh_token": tokens.RefreshToken.Value,
	})
}

type loginRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required"`
}

func LoginHandler(c echo.Context) error {
	var req loginRequest
	if httpErr := validator.BindAndValidateRequest(c, &req); httpErr != nil {
		return httpErr
	}

	loginInfo := &LoginInfo{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := VerifyLoginInfo(loginInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	tokens, err := NewJwtAuthTokens(AccessTokenCustomClaims{
		UserID: user.ID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = CreateSession(user.ID, tokens.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	SetJwtTokenInCookie(c, tokens)

	return c.JSON(http.StatusOK, map[string]string{
		"message":       "Login successful",
		"access_token":  tokens.AccessToken.Value,
		"refresh_token": tokens.RefreshToken.Value,
	})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshJwtTokensHandler(c echo.Context) error {
	var refreshToken string

	var req refreshTokenRequest
	httpErr := validator.BindAndValidateRequest(c, &req)
	if httpErr == nil {
		refreshToken = req.RefreshToken
	} else {
		tokenInCookie, err := c.Cookie("refresh-token")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		refreshToken = tokenInCookie.Value
	}

	claims, err := ParseRefreshToken(refreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	err = verifySession(claims.UserID, refreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	newTokens, err := NewJwtAuthTokens(AccessTokenCustomClaims{
		UserID: claims.UserID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = UpdateSession(claims.UserID, refreshToken, newTokens.RefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	SetJwtTokenInCookie(c, newTokens)

	return c.JSON(http.StatusOK, map[string]string{
		"message":       "Tokens refreshed",
		"access_token":  newTokens.AccessToken.Value,
		"refresh_token": newTokens.RefreshToken.Value,
	})
}

type RemoveSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RemoveSessionHandler(c echo.Context) error {
	var refreshToken string

	var req RemoveSessionRequest
	httpErr := validator.BindAndValidateRequest(c, &req)
	if httpErr == nil {
		refreshToken = req.RefreshToken
	} else {
		tokenInCookie, err := c.Cookie("refresh-token")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		refreshToken = tokenInCookie.Value
	}

	claims, err := ParseRefreshToken(refreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	if err := RemoveSession(claims.UserID, refreshToken); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, "Session removed")
}

func SetJwtTokenInCookie(c echo.Context, tokens JwtAuthTokens) {
	c.SetCookie(&http.Cookie{
		Name:    "access-token",
		Value:   tokens.AccessToken.Value,
		Expires: tokens.AccessToken.ExpiresAt,
	})
	c.SetCookie(&http.Cookie{
		Name:    "refresh-token",
		Value:   tokens.RefreshToken.Value,
		Expires: tokens.RefreshToken.ExpiresAt,
	})
}
