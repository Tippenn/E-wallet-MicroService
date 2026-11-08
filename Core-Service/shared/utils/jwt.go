package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my-bootcamp-secret-key-123")

func GenerateToken(Id string) (string, error) {
	claims := jwt.MapClaims{
		"user_Id": Id,
		"exp":     time.Now().Add(1 * time.Hour).Unix(), // expired at
		"iat":     time.Now().Unix(),                    // issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid or expired token: %w", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("Invalid claims")
	}

	id, ok := claims["user_Id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid id in token")
	}

	return id, nil
}
