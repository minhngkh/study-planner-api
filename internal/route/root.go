package route

import (
	"study-planner-api/internal/auth"
	"study-planner-api/internal/misc"
	"study-planner-api/internal/user"

	"github.com/labstack/echo/v4"
)

func RegisterRootRoutes(e *echo.Group) {
	e.GET("/", misc.HelloWorldHandler)
	e.GET("/health", misc.DBHealthHandler)

	e.POST("/register", auth.RegisterHandler)
	e.POST("/login", auth.LoginHandler)
	e.POST("/refresh", auth.RefreshJwtTokensHandler)

	e.GET("/profile", user.ProfileHandler, auth.JwtMiddleware())
}
