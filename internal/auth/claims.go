package auth

import "github.com/golang-jwt/jwt/v5"

// Claims — структура для работы с JWT
type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
