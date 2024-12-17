package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func GetUserInfoFromJwtToken(c echo.Context) *AccessTokenCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*AccessTokenClaims)

	log.Debug().Interface("info", claims.AccessTokenCustomClaims).Msg("User info from JWT token")

	return &claims.AccessTokenCustomClaims
}
