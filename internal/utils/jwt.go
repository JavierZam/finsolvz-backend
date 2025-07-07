package utils

import (
	"os"
	"time"

	"finsolvz-backend/internal/utils/errors"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET_MISSING", "JWT secret not configured", 500, nil, nil)
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("JWT_GENERATION_ERROR", "Failed to generate JWT token", 500, err, nil)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET_MISSING", "JWT secret not configured", 500, nil, nil)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("JWT_INVALID", "Invalid JWT token", 401, err, nil)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("JWT_INVALID", "Invalid JWT token claims", 401, nil, nil)
}
