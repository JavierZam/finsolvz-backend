package middleware

import (
	"context"
	"net/http"

	"finsolvz-backend/internal/utils"
	"finsolvz-backend/internal/utils/log"
)

type UserContext struct {
	UserID string
	Role   string
}

// AuthMiddleware validates JWT tokens and adds user context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Bearer token
		token, err := utils.ExtractBearerToken(r)
		if err != nil {
			log.Warnf(r.Context(), "Authentication failed: %v", err)
			utils.HandleHTTPError(w, err, r)
			return
		}

		// Validate JWT token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			log.Warnf(r.Context(), "Token validation failed: %v", err)
			utils.HandleHTTPError(w, err, r)
			return
		}

		// Add user context to request
		userCtx := &UserContext{
			UserID: claims.UserID,
			Role:   claims.Role,
		}

		ctx := context.WithValue(r.Context(), "user", userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts user context from request
func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	user, ok := ctx.Value("user").(*UserContext)
	return user, ok
}

// RequireRole creates middleware that requires specific roles
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				utils.HandleHTTPError(w, utils.ErrUnauthorized, r)
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				utils.HandleHTTPError(w, utils.ErrForbidden, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
