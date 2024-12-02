package auth

import (
	"study-planner-api/internal/db"
	"study-planner-api/internal/model"
)

func CreateSession(userID int32, refreshToken *JwtToken) error {
	result := db.Get().Create(&model.UserSession{
		UserID:       userID,
		RefreshToken: refreshToken.Value,
		ExpiresAt:    refreshToken.ExpiresAt,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateSession(userID int32, refreshToken string, newRefreshToken *JwtToken) error {
	result := db.Get().
		Model(&model.UserSession{}).
		Where("user_id = ? AND refresh_token = ?", userID, refreshToken).
		Updates(&model.UserSession{
			RefreshToken: newRefreshToken.Value,
			ExpiresAt:    newRefreshToken.ExpiresAt,
		}).
		Limit(1)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func verifySession(userID int32, refreshToken string) error {
	var session model.UserSession
	result := db.Get().
		Model(&model.UserSession{}).
		Where("user_id = ? AND refresh_token = ?", userID, refreshToken).
		First(&session)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func RemoveSession(userID int32, refreshToken string) error {
	result := db.Get().
		Where("user_id = ? AND refresh_token = ?", userID, refreshToken).
		Delete(&model.UserSession{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
