package route

import (
	"study-planner-api/internal/auth"
	"study-planner-api/internal/auth/provider"
	"study-planner-api/internal/misc"

	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(e *echo.Group) {
	e.POST("/google", provider.GoogleAuthHandler)

	e.POST("/refresh-token", auth.RefreshJwtTokensHandler)

	e.GET("/check", misc.EmptyHandler, auth.JwtMiddleware())

	e.POST("/logout", auth.RemoveSessionHandler, auth.JwtMiddleware())
}
