package user

import (
	"net/http"
	"study-planner-api/internal/auth"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func ProfileHandler(c echo.Context) error {
	user := auth.GetUserInfoFromJwtToken(c)

	log.Info().Int32("User requested profile", user.UserID)

	info, err := GetUserInfo(user.UserID)
	if err != nil {
		return echo.NewHTTPError(echo.ErrInternalServerError.Code, "Error getting user info")
	}

	return c.JSON(http.StatusOK, info)
}
