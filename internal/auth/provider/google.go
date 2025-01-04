package provider

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"study-planner-api/internal/auth"
	db "study-planner-api/internal/database"
	"study-planner-api/internal/model"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm/clause"
)

var (
	config = oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	resourceUrl = "https://www.googleapis.com/oauth2/v2/userinfo"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	PictureUrl    string `json:"picture"`
}

type StateToken struct {
	Type  string `json:"token_type"`
	Value string `json:"token_value"`
}

func GoogleAuthEndpoint(stateToken string) string {
	return config.AuthCodeURL(stateToken)
}

func CreateGoogleStateToken(app auth.RequestApplication) (auth.JwtToken, error) {
	return auth.CreateOauth2StateToken(app, "google")
}

func ValidateGoogleStateToken(token string) (auth.RequestApplication, error) {
	return auth.ValidateOauth2StateToken(token, "google")
}

func ExchangeWithGoogleForAuthTokens(authCode string) (authToken string, refreshToken string, err error) {
	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return "", "", err
	}

	return token.AccessToken, token.RefreshToken, nil
}

func GetGoogleUserInfo(accessToken string) (GoogleUserInfo, error) {
	url, err := url.Parse(resourceUrl)
	if err != nil {
		return GoogleUserInfo{}, err
	}

	query := url.Query()
	query.Add("access_token", accessToken)
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return GoogleUserInfo{}, err
	}

	rawInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		return GoogleUserInfo{}, err
	}

	var userInfo GoogleUserInfo
	err = json.Unmarshal(rawInfo, &userInfo)
	if err != nil {
		return GoogleUserInfo{}, err
	}

	return userInfo, nil
}

var (
	ErrInvalidGoogleAccount = errors.New("invalid google account")
)

func ValidateGoogleAccount(googleInfo GoogleUserInfo) (int32, error) {
	if !googleInfo.VerifiedEmail {
		return -1, ErrInvalidGoogleAccount
	}

	var user model.User
	result := db.Instance().
		Model(&model.User{}).
		Where("google_id = ?", googleInfo.ID).
		First(&user)
	if result.Error != nil {
		return -1, result.Error
	}
	if result.RowsAffected == 0 {
		return linkGoogleAccount(googleInfo)
	}

	return user.ID, nil
}

func linkGoogleAccount(googleInfo GoogleUserInfo) (int32, error) {
	user := model.User{
		GoogleID: &googleInfo.ID,
	}

	var users []model.User
	result := db.Instance().
		Model(&users).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Where("email = ?", googleInfo.Email).
		Updates(&user).
		Limit(1)
	if result.Error != nil {
		return -1, result.Error
	}
	if result.RowsAffected == 0 {
		user.Email = &googleInfo.Email

		result = db.Instance().
			Select("email", "google_id").
			Create(&user)
		if result.Error != nil {
			return -1, result.Error
		}
		return user.ID, nil
	}

	return users[0].ID, nil
}
