package user

import (
	"net/http"
	"study-planner-api/internal/auth"

	"github.com/labstack/echo/v4"
)

func ProfileHandler(c echo.Context) error {
	user := auth.GetUserInfoFromJwtToken(c)

	info, err := GetUserInfo(user.UserID)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Error getting user info")
	}

	return c.JSON(http.StatusOK, info)
}
