package jwtutils

import (
	"time"

	"github.com/Egot3/supel/backend/identity/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string, role types.UserRole) (string, error) {
	now := time.Now()
	expirationDate := now.Add(3600 * time.Second)

	claims := &Claims{
		UserID: userID, //неожидано
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

