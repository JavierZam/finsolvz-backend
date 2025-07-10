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
	ErrUserNotFound       = errors.New("USER_NOT_FOUND", "User not found", http.StatusNotFound, nil, nil)
	ErrEmailSendFailed    = errors.New("EMAIL_SEND_FAILED", "Failed to send email", http.StatusInternalServerError, nil, nil)
)
