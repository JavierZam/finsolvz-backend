package user

import (
	"time"                               // ✅ Added missing import
	"finsolvz-backend/internal/domain"
)

// Request DTOs
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=SUPER_ADMIN ADMIN CLIENT"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=SUPER_ADMIN ADMIN CLIENT"`
}

type UpdateRoleRequest struct {
	UserID  string `json:"userId" validate:"required"`
	NewRole string `json:"newRole" validate:"required,oneof=SUPER_ADMIN ADMIN CLIENT"`
}

type ChangePasswordRequest struct {
	NewPassword     string `json:"newPassword" validate:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=6"`
}

// Response DTOs
type UserResponse struct {
	ID        string    `json:"_id"`        // ✅ Changed to "_id" like legacy
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Company   []string  `json:"company"`
	CreatedAt time.Time `json:"createdAt"`  // ✅ Added missing field
	UpdatedAt time.Time `json:"updatedAt"`  // ✅ Added missing field
}

// Helper to convert domain.User to UserResponse
func ToUserResponse(user *domain.User) UserResponse {
	companyIDs := make([]string, len(user.Company))
	for i, id := range user.Company {
		companyIDs[i] = id.Hex()
	}

	return UserResponse{
		ID:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		Company:   companyIDs,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}