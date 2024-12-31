package provider

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	db "study-planner-api/internal/database"
	"study-planner-api/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	googleConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "postmessage",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
)

type GoogleLoginRequest struct {
	AuthCode string `json:"code" validate:"required"`
}

type GoogleUserInfo struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
}

// func GoogleAuthHandler(c echo.Context) error {
// 	var req googleLoginRequest
// 	if err := validator.BindAndValidateRequest(c, &req); err != nil {
// 		return err
// 	}

// 	AccessToken, err := googleConfig.Exchange(c.Request().Context(), req.AuthCode)
// 	if err != nil {
// 		log.Debug().Msg(err.Error())
// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid authorization code")
// 	}

// 	client := googleConfig.Client(c.Request().Context(), AccessToken)
// 	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "Failed to get user info from Google")
// 	}
// 	defer res.Body.Close()

// 	var info googleUserInfo
// 	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode user info")
// 	}

// 	var userID int32
// 	user, err := user.GetUserInfoByGoogleID(info.ID)
// 	if err == nil {
// 		userID = user.ID
// 	} else {
// 		userID, err = linkGoogleAccountToUser(&info)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError)
// 		}
// 	}

// 	// TODO: duplicated code with login handler in auth
// 	tokens, err := auth.NewJwtAuthTokens(auth.AccessTokenCustomClaims{
// 		UserID: userID,
// 	})
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create JWT tokens")
// 	}

// 	err = auth.CreateSession(userID, tokens.RefreshToken)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError)
// 	}

// 	return c.JSON(http.StatusOK, map[string]string{
// 		"access_token":  tokens.AccessToken.Value,
// 		"refresh_token": tokens.RefreshToken.Value,
// 	})
// }

// Link Google account to existing user, if not, creating an empty one
func linkGoogleAccountToUser(info *GoogleUserInfo) (int32, error) {
	// Attempt to create or update the user with GoogleID
	result := db.Instance().
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "email"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"google_id": gorm.Expr("CASE WHEN google_id IS NULL OR google_id = '' THEN ? ELSE google_id END", info.ID),
			}),
		}).
		Create(&model.User{
			Email:    info.Email,
			GoogleID: info.ID,
		})

	if result.Error != nil {
		return -1, result.Error
	}

	log.Info().Msgf("Start")

	// TODO: Migrate to sqlc soon, cause this is absurd
	// Fetch the user to ensure the ID is populated
	var user2 model.User
	result2 := db.Instance().Where("email = ?", info.Email).First(&user2)
	if result2.Error != nil {
		return -1, result2.Error
	}

	log.Info().Msgf("Linked Google ID: %s to user ID: %d", info.ID, user2.ID)

	return user2.ID, nil
}
