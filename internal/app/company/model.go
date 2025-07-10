package company

import (
	"time"

	"finsolvz-backend/internal/domain"
)

// Request DTOs
type CreateCompanyRequest struct {
	Name           string   `json:"name" validate:"required,min=2,max=100"`
	ProfilePicture *string  `json:"profilePicture,omitempty"`
	User           []string `json:"user,omitempty"` // Array of user IDs as strings
}

type UpdateCompanyRequest struct {
	Name           *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	ProfilePicture *string  `json:"profilePicture,omitempty"` // Simple URL string
	User           []string `json:"user,omitempty"`           // Array of user IDs as strings
}

// Response DTOs - exact legacy format
type CompanyResponse struct {
	ID             string     `json:"_id"` // ✅ Changed to "_id" exactly like legacy
	Name           string     `json:"name"`
	ProfilePicture *string    `json:"profilePicture"`
	User           []UserInfo `json:"user"` // Populated user data
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type UserInfo struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

// Helper to convert domain.Company to CompanyResponse
func ToCompanyResponse(company *domain.Company) CompanyResponse {
	return CompanyResponse{
		ID:             company.ID.Hex(),
		Name:           company.Name,
		ProfilePicture: company.ProfilePicture,
		User:           []UserInfo{}, // Will be populated by service layer
		CreatedAt:      company.CreatedAt,
		UpdatedAt:      company.UpdatedAt,
	}
}

// Helper to convert domain.Company to CompanyResponse with populated users
func ToCompanyResponseWithUsers(company *domain.Company, users []*domain.User) CompanyResponse {
	userInfos := make([]UserInfo, len(users))
	for i, user := range users {
		userInfos[i] = UserInfo{
			ID:   user.ID.Hex(),
			Name: user.Name,
		}
	}

	return CompanyResponse{
		ID:             company.ID.Hex(),
		Name:           company.Name,
		ProfilePicture: company.ProfilePicture,
		User:           userInfos,
		CreatedAt:      company.CreatedAt,
		UpdatedAt:      company.UpdatedAt,
	}
}
