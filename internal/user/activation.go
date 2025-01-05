package user

import (
	"errors"
	"html/template"
	"study-planner-api/internal/auth/token"
	"study-planner-api/internal/database"
	"study-planner-api/internal/model"
	"study-planner-api/internal/utils"
	"study-planner-api/internal/utils/email"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserAlreadyActivated = errors.New("user already activated")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidToken         = errors.New("invalid token")
	ErrExpiredToken         = errors.New("expired token")
	ErrCannotSendEmail      = errors.New("cannot send email")

	activationTemplate = template.Must(template.ParseFiles("activation.template.html"))
)

type activationEmailData struct {
	url string
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

	url := utils.GetServerHost()
	q := url.Query()
	q.Set("user_id", string(userId))
	q.Set("token", token)
	url.RawQuery = q.Encode()

	content, err := utils.CreateHtml(
		activationTemplate,
		activationEmailData{url: url.String()},
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
		Where("user.id = ?", userId).
		First(&t)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return result.Error
	}

	if token.VerifyHash(activationCode, t.TokenHash) {
		return ErrInvalidToken
	}

	if t.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	return nil
}
