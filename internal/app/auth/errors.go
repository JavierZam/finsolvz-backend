package auth

import (
	"finsolvz-backend/internal/utils/errors"
	"net/http"
)

var (
	ErrUserAlreadyExists  = errors.New("USER_ALREADY_EXISTS", "Email already registered", http.StatusConflict, nil, nil)
	ErrInvalidCredentials = errors.New("INVALID_CREDENTIALS", "Invalid email or password", http.StatusUnauthorized, nil, nil)
	ErrTokenExpired       = errors.New("TOKEN_EXPIRED", "Token has expired", http.StatusUnauthorized, nil, nil)
	ErrInvalidToken       = errors.New("INVALID_TOKEN", "Invalid token", http.StatusUnauthorized, nil, nil)
)
