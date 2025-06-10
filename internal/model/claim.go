package model

import "github.com/golang-jwt/jwt/v5"

// UserClaims структура для хранения информации о пользователе в JWT-токене
type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
