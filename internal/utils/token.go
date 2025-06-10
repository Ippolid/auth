package utils

import (
	"time"

	"github.com/Ippolid/auth/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// GenerateToken создает JWT-токен для пользователя с заданной информацией, секретным ключом и временем действия
func GenerateToken(info model.UserInfoJwt, secretKey []byte, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		Username: info.Username,
		Role:     map[bool]string{true: "admin", false: "user"}[info.Role],
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

// VerifyToken проверяет JWT-токен и возвращает информацию о пользователе, если токен действителен
func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.Errorf("invalid token claims")
	}

	return claims, nil
}
