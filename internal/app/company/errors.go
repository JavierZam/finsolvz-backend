package company

import (
	"finsolvz-backend/internal/utils/errors"
	"net/http"
)

var (
	ErrCompanyNotFound      = errors.New("COMPANY_NOT_FOUND", "Company not found", http.StatusNotFound, nil, nil)
	ErrCompanyAlreadyExists = errors.New("COMPANY_ALREADY_EXISTS", "Company name already exists", http.StatusConflict, nil, nil)
	ErrInvalidCompanyName   = errors.New("INVALID_COMPANY_NAME", "Company name is invalid", http.StatusBadRequest, nil, nil)
	ErrInvalidUserID        = errors.New("INVALID_USER_ID", "Invalid user ID format", http.StatusBadRequest, nil, nil)
	ErrUserNotFound         = errors.New("USER_NOT_FOUND", "User not found", http.StatusNotFound, nil, nil)
)