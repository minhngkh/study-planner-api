package user

import (
	"errors"
	db "study-planner-api/internal/database"
	"study-planner-api/internal/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	ID        int32
	Email     string
	CreatedAt time.Time
}

var (
	ErrUserExists = errors.New("user already exists")
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
	result := db.Instance().Select("Email", "Password").Create(&user)
	if result.RowsAffected == 0 {
		return -1, ErrUserExists
	}
	if result.Error != nil {
		return -1, result.Error
	}

	return user.ID, nil
}

func GetUserInfo(id int32) (UserInfo, error) {
	var info UserInfo
	result := db.Instance().
		Model(&model.User{ID: id}).
		Where("id = ?", id).
		First(&info)
	if result.Error != nil {
		return UserInfo{}, result.Error
	}

	return info, nil
}

func GetUserInfoByGoogleID(googleID string) (UserInfo, error) {
	var info UserInfo
	result := db.Instance().
		Model(&model.User{}).
		Where("google_id = ?", googleID).
		First(&info)
	if result.Error != nil {
		return UserInfo{}, result.Error
	}

	return info, nil
}
