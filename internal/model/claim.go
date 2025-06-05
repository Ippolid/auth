package model

import "github.com/dgrijalva/jwt-go"

// UserClaims структура для хранения информации о пользователе в JWT-токене
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
