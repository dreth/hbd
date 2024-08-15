package auth

import (
	"errors"
	"fmt"
	"hbd/env"
	"hbd/structs"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(env.MK)

// GenerateJWT generates a JWT token with the given email and duration
func GenerateJWT(email string, duration int) (string, error) {
	// If the duration is 0, default to 720
	if duration == 0 {
		duration = 720
	}

	// Set the expiration time for the token
	expirationTime := time.Now().Add(time.Duration(duration) * time.Hour)

	// Create the JWT claims, which includes the email and expiry time
	claims := &structs.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenStr string) (*structs.Claims, error) {
	// Parse the JWT token
	claims := &structs.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Check if the token is valid
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetJWTDurationFromHeader gets the JWT duration from the header
func GetJWTDurationFromHeader(c *gin.Context, defaultDuration int) (int, error) {
	// Get the JWT duration from the header
	jwtDurationStr := c.GetHeader("X-Jwt-Token-Duration")

	// If the duration is not provided, return the default duration
	if jwtDurationStr != "" {
		jwtDuration, err := strconv.Atoi(jwtDurationStr)
		if err != nil {
			return 0, fmt.Errorf("invalid JWT duration")
		}
		return jwtDuration, nil
	}

	return defaultDuration, nil
}
