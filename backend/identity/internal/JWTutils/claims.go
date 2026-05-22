package jwtutils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct { //Метаморфоз https://purpleschool.ru/knowledge-base/article/work-with-jwt
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}
