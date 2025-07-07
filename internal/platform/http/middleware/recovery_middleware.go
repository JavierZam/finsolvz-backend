package middleware

import (
	"net/http"
	"runtime/debug"

	"finsolvz-backend/internal/utils"
	"finsolvz-backend/internal/utils/log"
)

// RecoveryMiddleware recovers from panics and returns a 500 error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf(r.Context(), "Panic recovered: %v\nStack trace:\n%s", err, debug.Stack())

				// Return internal server error
				utils.HandleHTTPError(w, utils.ErrInternalServer, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
