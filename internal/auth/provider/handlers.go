package provider

import (
	"net/http"
	"study-planner-api/internal/validator"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/rs/zerolog/log"
)

func AuthCallbackHandler(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		log.Debug().Msg(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, user)
}

type AuthInitRequest struct {
	Provider string `param:"provider" validate:"required"`
}

func AuthInitHandler(c echo.Context) error {
	var req AuthInitRequest
	if httpErr := validator.BindAndValidateRequest(c, &req); httpErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unspecified provider")
	}

	queryString := c.Request().URL.Query()
	queryString.Add("provider", req.Provider)
	c.Request().URL.RawQuery = queryString.Encode()

	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		gothic.BeginAuthHandler(c.Response().Writer, c.Request())
		return nil
	}

	return c.JSON(http.StatusOK, user)
}
