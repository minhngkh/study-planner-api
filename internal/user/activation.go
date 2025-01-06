package user

import (
	"errors"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"study-planner-api/internal/auth/token"
	"study-planner-api/internal/database"
	"study-planner-api/internal/model"
	"study-planner-api/internal/utils"
	"study-planner-api/internal/utils/email"
	"time"

	"gorm.io/gorm"
)

func getActivationTemplate() *template.Template {
	// curDir := utils.CurrentFileDir()
	path := filepath.Join("templates", "account-activation.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}

	return tmpl
}

func getActivationCallbackUrl() *url.URL {
	url, err := url.Parse(os.Getenv("ACTIVATION_CALLBACK_URL"))
	if err != nil {
		panic(err)
	}

	return url
}

var (
	ErrUserAlreadyActivated = errors.New("user already activated")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("expired token")
	ErrCannotSendEmail      = errors.New("cannot send email")

	activationTemplate    = getActivationTemplate()
	activationCallbackUrl = getActivationCallbackUrl()
)

type activationEmailData struct {
	Url string
}

func SendActivationEmail(userId int32) error {
	var user model.User
	result := database.Instance().
		Model(&model.User{}).
		Where("id = ?", userId).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return result.Error
	}

	if user.IsActivated {
		return ErrUserAlreadyActivated
	}

	token, err := token.CreateToken(userId, token.Activation)
	if err != nil {
		return err
	}

	url := activationCallbackUrl
	q := url.Query()
	q.Set("user_id", strconv.Itoa(int(userId)))
	q.Set("token", token)
	url.RawQuery = q.Encode()

	content, err := utils.CreateHtml(
		activationTemplate,
		activationEmailData{Url: url.String()},
	)
	if err != nil {
		return err
	}

	err = email.Send(*user.Email, "Confirm your email", content)
	if err != nil {
		return ErrCannotSendEmail
	}

	return nil
}

func ActivateAccount(userId int32, activationCode string) error {
	var t model.Token
	result := database.Instance().
		Select("token.*").
		Model(&model.User{}).
		Joins("INNER JOIN token ON user.id = token.user_id").
		Where("user.id = ? AND token.purpose = ?", userId, token.Activation.String()).
		First(&t)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrInvalidToken
		}
		return result.Error
	}

	if !token.VerifyHash(activationCode, t.TokenHash) {
		return ErrInvalidToken
	}

	if t.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	result = database.Instance().
		Model(&model.User{}).
		Where("id = ?", userId).
		Update("is_activated", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cannot update, something went wrong")
	}

	result = database.Instance().
		Model(&model.Token{}).
		Where("id = ?", t.ID).
		Delete(&model.Token{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("cannot update, something went wrong")
	}

	return nil
}
