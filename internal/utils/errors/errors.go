package errors

import (
	"fmt"
	"net/http"
)

// AppError interface yang harus diimplementasikan oleh semua custom error aplikasi.
type AppError interface {
	error
	Code() string
	Message() string
	Status() int
	Details() map[string]interface{}
	Unwrap() error
}

// baseError implementasi dasar AppError.
type baseError struct {
	err     error
	code    string
	message string
	status  int
	details map[string]interface{}
}

func (e *baseError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.message, e.err)
	}
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

func (e *baseError) Code() string                    { return e.code }
func (e *baseError) Message() string                 { return e.message }
func (e *baseError) Status() int                     { return e.status }
func (e *baseError) Details() map[string]interface{} { return e.details }
func (e *baseError) Unwrap() error                   { return e.err }

// New adalah konstruktor untuk membuat AppError baru.
func New(code, message string, status int, originalErr error, details map[string]interface{}) AppError {
	return &baseError{
		err:     originalErr,
		code:    code,
		message: message,
		status:  status,
		details: details,
	}
}

// Variabel error global yang sudah didefinisikan
var (
	ErrBadRequest     = New("BAD_REQUEST", "Invalid request payload or parameters", http.StatusBadRequest, nil, nil)
	ErrUnauthorized   = New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized, nil, nil)
	ErrForbidden      = New("FORBIDDEN", "Access denied", http.StatusForbidden, nil, nil)
	ErrNotFound       = New("NOT_FOUND", "Resource not found", http.StatusNotFound, nil, nil)
	ErrInternalServer = New("INTERNAL_SERVER_ERROR", "An unexpected internal server error occurred", http.StatusInternalServerError, nil, nil)
	ErrConflict       = New("CONFLICT", "Resource conflict", http.StatusConflict, nil, nil)
)
