package route

import (
	"study-planner-api/internal/auth"
	"study-planner-api/internal/misc"
	"study-planner-api/internal/user"

	"github.com/labstack/echo/v4"
)

func RegisterRootRoutes(e *echo.Group) {
	RegisterAuthRoutes(e.Group("/auth"))

	e.GET("/", misc.HelloWorldHandler)
	e.GET("/health", misc.DBHealthHandler)

	e.POST("/register", auth.RegisterHandler)
	e.POST("/login", auth.LoginHandler)

	e.GET("/profile", user.ProfileHandler, auth.JwtMiddleware())
}
