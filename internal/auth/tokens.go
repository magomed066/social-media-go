package auth

import (
	"fmt"
	"social-media-app/internal/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAuthToken(userID int64, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userID":   userID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := env.GetString("JWT_ACCESS_SECRET", "")
	if secret == "" {
		return "", fmt.Errorf("JWT_ACCESS_SECRET is not configured")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	secret := env.GetString("JWT_ACCESS_SECRET", "")
	if secret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is not configured")
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
