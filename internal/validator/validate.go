package validator

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BindAndValidateRequest[T any](c echo.Context, req *T) *echo.HTTPError {
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload")
	}

	return nil
}
