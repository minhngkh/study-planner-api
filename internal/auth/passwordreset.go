package auth

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
	"study-planner-api/internal/user"
	"study-planner-api/internal/utils"
	"study-planner-api/internal/utils/email"
	"time"

	"gorm.io/gorm"
)

func getPasswordResetTemplate() *template.Template {
	curDir := utils.CurrentFileDir()
	path := filepath.Join(curDir, "passwordreset.template.html")

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(err)
	}

	return tmpl
}

func getPasswordResetCallbackUrl() *url.URL {
	url, err := url.Parse(os.Getenv("PASSWORD_RESET_CALLBACK_URL"))
	if err != nil {
		panic(err)
	}

	return url
}

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("expired token")
	ErrCannotSendEmail = errors.New("cannot send email")
	ErrUnknownEmail    = errors.New("unknown email")

	passwordResetTemplate    = getPasswordResetTemplate()
	passwordResetCallbackUrl = getPasswordResetCallbackUrl()
)

type passwordResetEmailData struct {
	Url string
}

func SendPasswordResetEmail(userEmail string) error {
	var user model.User
	result := database.Instance().
		Model(&model.User{}).
		Where("email = ?", userEmail).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrUnknownEmail
		}
		return result.Error
	}

	token, err := token.CreateToken(user.ID, token.PasswordReset)
	if err != nil {
		return err
	}

	url := passwordResetCallbackUrl
	q := url.Query()
	q.Set("user_id", strconv.Itoa(int(user.ID)))
	q.Set("token", token)
	url.RawQuery = q.Encode()

	content, err := utils.CreateHtml(
		passwordResetTemplate,
		passwordResetEmailData{Url: url.String()},
	)
	if err != nil {
		return err
	}

	err = email.Send(userEmail, "Reset your password", content)
	if err != nil {
		return ErrCannotSendEmail
	}

	return nil
}

func ResetPassword(userId int32, resetToken string, newPassword string) error {
	var t model.Token
	result := database.Instance().
		Select("token.*").
		Model(&model.User{}).
		Joins("INNER JOIN token ON user.id = token.user_id").
		Where("user.id = ? AND token.purpose = ?", userId, token.PasswordReset.String()).
		First(&t)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return result.Error
	}

	if !token.VerifyHash(resetToken, t.TokenHash) {
		return ErrInvalidToken
	}

	if t.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	err := user.UpdatePassword(userId, newPassword)
	if err != nil {
		return errors.New("cannot update password, something went wrong")
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

func VerifyPasswordResetToken(userId int32, resetToken string) error {
	var t model.Token
	result := database.Instance().
		Select("token.*").
		Model(&model.User{}).
		Joins("INNER JOIN token ON user.id = token.user_id").
		Where("user.id = ? AND token.purpose = ?", userId, token.PasswordReset.String()).
		First(&t)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrInvalidToken
		}
		return result.Error
	}

	if !token.VerifyHash(resetToken, t.TokenHash) {
		return ErrInvalidToken
	}

	if t.ExpiresAt.Before(time.Now()) {
		return ErrExpiredToken
	}

	return nil
}
