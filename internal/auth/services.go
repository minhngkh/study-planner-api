package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	db "study-planner-api/internal/database"
	"study-planner-api/internal/model"
)

type LoginInfo struct {
	Email    string
	Password string
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserHasNoPassword = errors.New("user has no password")
)

func VerifyLoginInfo(info LoginInfo) (model.User, error) {

	var user model.User
	result := db.Instance().Where("email = ?", info.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, result.Error
	}

	if user.Password == nil {
		return model.User{}, ErrUserHasNoPassword
	}

	err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(info.Password))
	if err != nil {
		return model.User{}, ErrIncorrectPassword
	}

	return user, nil
}
