package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"study-planner-api/internal/db"
	"study-planner-api/internal/model"
)

type LoginInfo struct {
	Email    string
	Password string
}

var (
	ErrUserExists        = errors.New("User already exists")
	ErrUserNotFound      = errors.New("User not found")
	ErrIncorrectPassword = errors.New("Incorrect password")
)

func CreateUser(email string, password string) (id int32, err error) {
	// Hash password with bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	user := model.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	// Create user in database
	result := db.Get().Select("Email", "Password").Create(&user)
	if result.RowsAffected == 0 {
		return -1, ErrUserExists
	}
	if result.Error != nil {
		return -1, result.Error
	}

	return user.ID, nil
}

func VerifyLoginInfo(info *LoginInfo) (*model.User, error) {
	var user model.User
	result := db.Get().Where("email = ?", info.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password))
	if err != nil {
		return nil, ErrIncorrectPassword
	}

	return &user, nil
}
