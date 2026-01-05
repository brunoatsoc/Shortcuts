package service

import (
	"intellisys-backend/internal/configuration"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(configuration.Config.JWTSecret)
}
