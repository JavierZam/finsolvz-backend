package user

import (
	"finsolvz-backend/internal/utils/errors"
	"net/http"
)

var (
	ErrUserNotFound       = errors.New("USER_NOT_FOUND", "User not found", http.StatusNotFound, nil, nil)
	ErrEmailAlreadyExists = errors.New("EMAIL_ALREADY_EXISTS", "Email already used by another user", http.StatusConflict, nil, nil)
	ErrPasswordMismatch   = errors.New("PASSWORD_MISMATCH", "Passwords do not match", http.StatusBadRequest, nil, nil)
	ErrUnauthorizedAccess = errors.New("UNAUTHORIZED_ACCESS", "You are not authorized to perform this action", http.StatusForbidden, nil, nil)
)
