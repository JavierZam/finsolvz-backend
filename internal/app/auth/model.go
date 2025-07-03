package auth

import (
	"finsolvz-backend/internal/domain"
)

// Request DTOs
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=SUPER_ADMIN ADMIN CLIENT"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// Response DTOs
type AuthResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Role    string   `json:"role"`
	Company []string `json:"company"`
}

// Helper to convert domain.User to UserInfo
func ToUserInfo(user *domain.User) UserInfo {
	companyIDs := make([]string, len(user.Company))
	for i, id := range user.Company {
		companyIDs[i] = id.Hex()
	}

	return UserInfo{
		ID:      user.ID.Hex(),
		Name:    user.Name,
		Email:   user.Email,
		Role:    string(user.Role),
		Company: companyIDs,
	}
}
