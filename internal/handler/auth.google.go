package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"study-planner-api/internal/api"
	"study-planner-api/internal/auth"
	"study-planner-api/internal/auth/provider"

	"github.com/rs/zerolog/log"
)

func (s *Handler) GetAuthGoogleAuthorize(ctx context.Context, request api.GetAuthGoogleAuthorizeRequestObject) (api.GetAuthGoogleAuthorizeResponseObject, error) {
	hostUrl := api.HostUrlOfRequest(ctx)
	log.Debug().Msgf("Host URL: %s", hostUrl)

	stateToken, err := provider.CreateGoogleStateToken(auth.RequestApplication{Host: hostUrl})
	if err != nil {
		return nil, err
	}

	authUrl := provider.GoogleAuthEndpoint(stateToken.Value)

	cookie := http.Cookie{
		Name:  "google_auth",
		Value: stateToken.Value,

		// HttpOnly: true,
		// Secure: true,
	}

	return api.GetAuthGoogleAuthorize303Response{
		Headers: api.GetAuthGoogleAuthorize303ResponseHeaders{
			Location:  authUrl,
			SetCookie: cookie.String(),
		},
	}, nil

}

func (s *Handler) GetAuthGoogleCallback(ctx context.Context, request api.GetAuthGoogleCallbackRequestObject) (api.GetAuthGoogleCallbackResponseObject, error) {
	authCode := request.Params.Code
	stateToken := request.Params.State

	_, err := provider.ValidateGoogleStateToken(stateToken)
	if err != nil {
		return nil, err
	}

	googleAccessToken, _, err := provider.ExchangeWithGoogleForAuthTokens(authCode)
	if err != nil {
		return nil, err
	}

	googleInfo, err := provider.GetGoogleUserInfo(googleAccessToken)
	if err != nil {
		return nil, err
	}

	userID, err := provider.ValidateGoogleAccount(googleInfo)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to validate google account"))
	}

	accessToken, refreshToken, err := auth.CreateAuthTokens(auth.AuthInfo{UserID: userID})
	if err != nil {
		return nil, err
	}

	err = auth.CreateSession(userID, refreshToken)
	if err != nil {
		return nil, err
	}

	resp := `
		<!doctype html>
		<html>
		<title>Google Auth Successful</title>
		<script>
			window.opener.postMessage(
				{
					accessToken: "` + accessToken.Value + `",
					refreshToken: "` + refreshToken.Value + `"
				}, 
				'*'
			);
			window.close();
		</script>
		</html>
	`

	return api.GetAuthGoogleCallback200TexthtmlResponse{
		Body:          strings.NewReader(resp),
		ContentLength: int64(len(resp)),
	}, nil

	// return api.GetAuthGoogleCallback301Response{
	// 	Headers: api.GetAuthGoogleCallback301ResponseHeaders{
	// 		Location: fmt.Sprintf("%s#access_token=%s", hostUrl, token),
	// 	},
	// }, nil
}
