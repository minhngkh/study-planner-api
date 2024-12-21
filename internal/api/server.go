package api

import (
	"context"
	"errors"
	"study-planner-api/internal/auth"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	oapiEchoMiddleware "github.com/oapi-codegen/echo-middleware"
	"github.com/rs/zerolog/log"
)

func NewEchoHandler() *echo.Echo {
	e := echo.New()

	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v echoMiddleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Info().
					Str("method", v.Method).
					Str("uri", v.URI).
					Int("status", v.Status).
					Msg("request")
			} else {
				log.Error().
					Str("method", v.Method).
					Str("uri", v.URI).
					Int("status", v.Status).
					Err(v.Error).
					Msg("request")
			}

			return nil
		},
	}))

	e.Use(echoMiddleware.Recover())

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderCookie,
		},
		AllowCredentials: true,
	}))

	specs, err := GetSwagger()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get swagger")
	}
	specs.Servers = nil

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			cookie, _ := c.Cookie("refresh_token")

			log.Debug().Interface("cookie", cookie).Msg("cookie")
			extendRequestContext(req, extendedCtx.host, req.Host)

			return next(c)
		}
	})

	e.Use(oapiEchoMiddleware.OapiRequestValidatorWithOptions(specs, &oapiEchoMiddleware.Options{
		ErrorHandler: func(c echo.Context, err *echo.HTTPError) error {
			return err
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
				log.Debug().Str("security_scheme_name", ai.SecuritySchemeName).Msg("validation")

				switch ai.SecuritySchemeName {
				case "bearerAuth":
					var token string
					req := ai.RequestValidationInput.Request
					authHeader := req.Header.Get("Authorization")

					switch {
					case authHeader == "":
						return errors.New("missing Authorization header")
					case len(authHeader) < 7 || authHeader[:7] != "Bearer ":
						return errors.New("invalid Authorization header")
					default:
						token = authHeader[7:]
					}

					info, _, err := auth.ValidateAccessToken(token)
					if err != nil {
						return errors.New("invalid token")
					}

					log.Debug().Interface("auth_info", info).Msg("validation")
					// log.Debug().Interface("cookie", req.Cookies()).Msg("cookie")
					// echoCtx := oapiEchoMiddleware.GetEchoContext(ctx)
					// ch := make(chan bool, 1)
					// go func() {
					// 	authCtx := context.WithValue(ctx, ContextKey("auth"), AuthInfo{ID: claims.UserID})
					// 	echoCtx.SetRequest(echoCtx.Request().WithContext(authCtx))
					// 	ch <- true
					// }()
					// <-ch

					extendRequestContext(req, extendedCtx.authInfo, authInfo{ID: info.UserID})

					return nil
				default:
					return errors.New("unimplemented security scheme")
				}
			},
		},
	}))

	return e
}
