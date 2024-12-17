package auth

import (
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type AccessTokenCustomClaims struct {
	UserID int32 `json:"id"`
}

type RefreshTokenCustomClaims struct {
	UserID int32 `json:"id"`
}

type AccessTokenClaims struct {
	AccessTokenCustomClaims
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	RefreshTokenCustomClaims
	jwt.RegisteredClaims
}

type JwtToken struct {
	Value     string
	ExpiresAt time.Time
}

type JwtAuthTokens struct {
	AccessToken  JwtToken
	RefreshToken JwtToken
}

var (
	jwtSecret = os.Getenv("JWT_SECRET")

	AccessTokenDuration  = time.Minute * 15
	RefreshTokenDuration = time.Hour * 24 * 7

	// AccessTokenDuration  = time.Second * 1
	// RefreshTokenDuration = time.Second * 1
)

func NewJwtAuthTokens(info AccessTokenCustomClaims) (JwtAuthTokens, error) {
	var accessToken, refreshToken JwtToken
	var accessErr, refreshErr error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		accessToken, accessErr = newAccessToken(info)
	}()

	go func() {
		defer wg.Done()
		refreshToken, refreshErr = newRefreshToken(RefreshTokenCustomClaims(info))
	}()

	wg.Wait()

	if accessErr != nil {
		return JwtAuthTokens{}, accessErr
	}
	if refreshErr != nil {
		return JwtAuthTokens{}, refreshErr
	}

	return JwtAuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func newAccessToken(info AccessTokenCustomClaims) (JwtToken, error) {
	expireTime := time.Now().Add(AccessTokenDuration)
	claims := &AccessTokenClaims{
		info,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{
		Value:     signedToken,
		ExpiresAt: expireTime,
	}, nil
}

func newRefreshToken(info RefreshTokenCustomClaims) (JwtToken, error) {
	expireTime := time.Now().Add(RefreshTokenDuration)
	claims := &RefreshTokenClaims{
		info,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{
		Value:     signedToken,
		ExpiresAt: expireTime,
	}, nil
}

func ParseRefreshToken(token string) (*RefreshTokenClaims, error) {
	claims := new(RefreshTokenClaims)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return new(RefreshTokenClaims), err
	}

	return claims, nil
}

func ParseAccessToken(token string) (*AccessTokenClaims, error) {
	claims := new(AccessTokenClaims)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return new(AccessTokenClaims), err
	}

	return claims, nil
}

func ParseJwtToken[T jwt.Claims](token string) (T, error) {
	claims := new(T)
	_, err := jwt.ParseWithClaims(token, *claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return *new(T), err
	}

	return *claims, nil
}

func ValidateJwtToken(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	return err
}
