package auth

import (
	"errors"
	db "study-planner-api/internal/database"
	"study-planner-api/internal/model"
)

var (
	ErrMaliciousRefreshToken = errors.New("malicious refresh token")
)

func CreateSession(userID int32, refreshToken JwtToken) error {
	expirationTime := refreshToken.Expiry.Time() // ok to be nill

	result := db.Instance().Create(&model.UserSession{
		UserID:       userID,
		RefreshToken: string(refreshToken.Value),
		ExpiresAt:    expirationTime,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func UpdateSession(userID int32, oldRefreshTokenVal string, newRefreshToken JwtToken) error {
	expirationTime := newRefreshToken.Expiry.Time() // ok to be nill

	result := db.Instance().
		Model(&model.UserSession{}).
		Where("user_id = ? AND refresh_token = ?", userID, oldRefreshTokenVal).
		Updates(&model.UserSession{
			RefreshToken: newRefreshToken.Value,
			ExpiresAt:    expirationTime,
		}).
		Limit(1)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMaliciousRefreshToken
	}

	return nil
}

func VerifySession(userID int32, refreshToken JwtToken) error {
	var session model.UserSession
	result := db.Instance().
		Model(&model.UserSession{}).
		Where("user_id = ? AND refresh_token = ?", userID, refreshToken.Value).
		First(&session)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMaliciousRefreshToken
	}

	return nil
}

func RemoveSession(userID int32, refreshTokenVal string) error {
	result := db.Instance().
		Where("user_id = ? AND refresh_token = ?", userID, refreshTokenVal).
		Delete(&model.UserSession{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrMaliciousRefreshToken
	}

	return nil
}
