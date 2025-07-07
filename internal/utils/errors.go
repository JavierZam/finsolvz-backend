package utils

import "finsolvz-backend/internal/utils/errors"

// Export error globals untuk kemudahan akses
var (
	ErrBadRequest     = errors.ErrBadRequest
	ErrUnauthorized   = errors.ErrUnauthorized
	ErrForbidden      = errors.ErrForbidden
	ErrNotFound       = errors.ErrNotFound
	ErrInternalServer = errors.ErrInternalServer
	ErrConflict       = errors.ErrConflict
)
