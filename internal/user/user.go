package user

import (
	"study-planner-api/internal/db"
	"study-planner-api/internal/model"
	"time"
)

type UserInfo struct {
	ID        int32
	Email     string
	CreatedAt time.Time
}

func GetUserInfo(id int32) (UserInfo, error) {
	var info UserInfo
	result := db.Get().
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
	result := db.Get().
		Model(&model.User{}).
		Where("google_id = ?", googleID).
		First(&info)
	if result.Error != nil {
		return UserInfo{}, result.Error
	}

	return info, nil
}
