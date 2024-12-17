package auth

import (
	"github.com/golang-jwt/jwt/v5"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware() echo.MiddlewareFunc {
	return echoJwt.WithConfig(echoJwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(AccessTokenClaims)
		},
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "header:Authorization:Bearer ,cookie:jwt",
	})
}
