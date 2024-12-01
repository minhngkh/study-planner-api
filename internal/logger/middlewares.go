package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func LogWithZerolog() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Info().
					Str("URI", v.URI).
					Int("status", v.Status).
					Msg("request")
			} else {
				log.Error().
					Str("URI", v.URI).
					Int("status", v.Status).
					Err(v.Error).
					Msg("request")
			}

			return nil
		},
	})
}
