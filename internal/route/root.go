package route

import (
	"net/http"
	"study-planner-api/internal/auth"
	"study-planner-api/internal/auth/provider"
	"study-planner-api/internal/misc"
	"study-planner-api/internal/user"

	"github.com/labstack/echo/v4"
)

func RegisterRootRoutes(e *echo.Group) {
	e.GET("/", misc.HelloWorldHandler)
	e.GET("/health", misc.DBHealthHandler)

	e.POST("/register", auth.RegisterHandler)
	e.POST("/login", auth.LoginHandler)

	e.GET("/profile", user.ProfileHandler, auth.JwtMiddleware())

	e.POST("/auth/google", provider.GoogleAuthHandler)
	e.POST("/auth/refresh-token", auth.RefreshJwtTokensHandler)
	e.GET("/auth/me", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}, auth.JwtMiddleware())
}
