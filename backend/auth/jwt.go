package auth

import (
	"errors"
	"hbd/env"
	"hbd/structs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(env.MK)

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(720 * time.Hour)
	claims := &structs.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*structs.Claims, error) {
	claims := &structs.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
