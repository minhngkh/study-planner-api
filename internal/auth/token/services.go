package token

import (
	"errors"
	"study-planner-api/internal/database"
	"study-planner-api/internal/model"
	"time"

	"gorm.io/gorm"
)

const (
	TokenLength = 32
)

var (
	ErrNoTokenFound = errors.New("no token found")
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type TokenPurpose struct {
	Duration time.Duration
	alias    string
}

var (
	Activation = TokenPurpose{
		Duration: time.Hour * 24,
		alias:    "activation",
	}
	PasswordReset = TokenPurpose{
		Duration: time.Hour * 1,
		alias:    "password_reset",
	}
)

func (tp TokenPurpose) String() string {
	return tp.alias
}

func TokenPurposeFromString(str string) (TokenPurpose, error) {
	switch str {
	case Activation.alias:
		return Activation, nil
	case PasswordReset.alias:
		return PasswordReset, nil
	default:
		return *new(TokenPurpose), errors.New("invalid token purpose")
	}
}

func CreateToken(userId int32, purpose TokenPurpose) (string, error) {
	token, err := GenerateRandomToken(TokenLength)
	if err != nil {
		return "", err
	}

	hash := HashToken(token)
	expirationTime := time.Now().Add(purpose.Duration)

	result := database.Instance().
		Model(&model.Token{}).
		Create(&model.Token{
			UserID:    userId,
			TokenHash: hash,
			Purpose:   purpose.String(),
			ExpiresAt: expirationTime,
		})
	if result.Error != nil {
		return "", result.Error
	}

	return token, nil
}

func VerifyToken(userId int32, token string, purpose TokenPurpose) error {
	var tokenModel model.Token
	result := database.Instance().
		Model(&model.Token{}).
		Where("user_id = ? AND purpose = ?", userId, purpose.String()).
		First(&tokenModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ErrNoTokenFound
		}
		return result.Error
	}

	if !VerifyHash(token, tokenModel.TokenHash) {
		return ErrInvalidToken
	}

	if tokenModel.ExpiresAt.Before(time.Now()) {
		return ErrTokenExpired
	}

	return nil
}
