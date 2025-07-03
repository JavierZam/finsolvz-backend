package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"finsolvz-backend/internal/utils/errors"

	"github.com/go-playground/validator/v10"
)

// DecodeJSON decodes JSON request body
func DecodeJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return errors.New("INVALID_JSON", "Invalid JSON format", 400, err, nil)
	}

	return nil
}

// HandleValidationError handles validation errors from go-playground/validator
func HandleValidationError(w http.ResponseWriter, err error, r *http.Request) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]interface{})
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			switch fieldError.Tag() {
			case "required":
				details[field] = "This field is required"
			case "email":
				details[field] = "Please provide a valid email address"
			case "min":
				details[field] = "This field is too short"
			case "max":
				details[field] = "This field is too long"
			case "oneof":
				details[field] = "Invalid value provided"
			default:
				details[field] = "Invalid value"
			}
		}

		validationErr := errors.New("VALIDATION_ERROR", "Invalid input data", 400, err, details)
		HandleHTTPError(w, validationErr, r)
		return
	}

	// Fallback for other validation errors
	HandleHTTPError(w, ErrBadRequest, r)
}

// ExtractBearerToken extracts Bearer token from Authorization header
func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("MISSING_AUTH_HEADER", "Authorization header is required", 401, nil, nil)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("INVALID_AUTH_FORMAT", "Authorization header must be in Bearer format", 401, nil, nil)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("MISSING_TOKEN", "Token is missing from Authorization header", 401, nil, nil)
	}

	return token, nil
}
