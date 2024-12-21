package api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type contextKey string

type extendedContext struct {
	authInfo contextKey
	host     contextKey
}

type authInfo struct {
	ID int32
}

var (
	extendedCtx = extendedContext{
		authInfo: contextKey("auth_info"),
		host:     contextKey("host"),
	}
)

// Extend request context with additional info.
func extendRequestContext(req *http.Request, key any, value any) {
	authCtx := context.WithValue(req.Context(), key, value)
	*req = *req.WithContext(authCtx)
}

// Get auth info from request context if authenticated.
// Exit if auth info not found.
func AuthInfoOfRequest(ctx context.Context) authInfo {
	authInfo, ok := ctx.Value(extendedCtx.authInfo).(authInfo)
	if !ok {
		log.Fatal().Msg("auth info not found in request context")
	}

	return authInfo
}

// Get host info from request context.
// Exit if host info not found.
func HostUrlOfRequest(ctx context.Context) string {
	hostInfo, ok := ctx.Value(extendedCtx.host).(string)
	if !ok {
		log.Fatal().Msg("host url not found in request context")
	}

	return hostInfo
}
