package misc

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"study-planner-api/internal/auth"
	"study-planner-api/internal/db"
)

func HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func DBHealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, db.Get().Health())
}

func TestHandler(c echo.Context) error {
	test := auth.GetUserInfoFromJwtToken(c)

	return c.JSON(http.StatusOK, test)
}
