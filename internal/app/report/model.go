package report

import (
	"time"

	"finsolvz-backend/internal/domain"
)

// ✅ FIXED: Request DTOs - exact field names sesuai dengan legacy Node.js
type CreateReportRequest struct {
	ReportName string      `json:"reportName" validate:"required,min=1,max=200"`
	ReportType string      `json:"reportType" validate:"required"`
	Year       string      `json:"year" validate:"required"`
	Company    string      `json:"company" validate:"required"`
	Currency   *string     `json:"currency,omitempty"`
	CreateBy   string      `json:"createBy" validate:"required"` // ✅ FIXED: "createBy" bukan "createdBy"
	UserAccess []string    `json:"userAccess,omitempty"`
	ReportData interface{} `json:"reportData,omitempty"`
}

type UpdateReportRequest struct {
	ReportName *string     `json:"reportName,omitempty" validate:"omitempty,min=1,max=200"`
	ReportType *string     `json:"reportType,omitempty"`
	Year       *string     `json:"year,omitempty"`
	Company    *string     `json:"company,omitempty"`
	Currency   *string     `json:"currency,omitempty"`
	UserAccess []string    `json:"userAccess,omitempty"`
	ReportData interface{} `json:"reportData,omitempty"`
}

type GetReportsByCompaniesRequest struct {
	CompanyIds []string `json:"companyIds" validate:"required,min=2"` // ✅ Legacy expects "companyIds"
}

// ✅ Response DTOs - EXACT format seperti legacy Node.js dengan populate
type ReportResponse struct {
	ID         string          `json:"_id"`
	ReportName string          `json:"reportName"`
	ReportType *ReportTypeInfo `json:"reportType"`
	Year       string          `json:"year"` // ✅ Always string
	Company    *CompanyInfo    `json:"company"`
	Currency   *string         `json:"currency"`
	CreatedBy  *UserInfo       `json:"createdBy"` // ✅ Response uses "createdBy"
	UserAccess []*UserInfo     `json:"userAccess"`
	ReportData interface{}     `json:"reportData"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

// Nested response types untuk populated data (exact legacy format)
type ReportTypeInfo struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type CompanyInfo struct {
	ID             string    `json:"_id"`
	Name           string    `json:"name"`
	ProfilePicture *string   `json:"profilePicture"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type UserInfo struct {
	ID        string    `json:"_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ✅ ENHANCED: Helper functions untuk konversi domain ke response
func ToReportResponse(report *domain.PopulatedReport) *ReportResponse {
	response := &ReportResponse{
		ID:         report.ID.Hex(),
		ReportName: report.ReportName,
		Year:       report.Year, // Already guaranteed to be string from repo fix
		Currency:   report.Currency,
		ReportData: report.ReportData,
		CreatedAt:  report.CreatedAt,
		UpdatedAt:  report.UpdatedAt,
	}

	// ✅ Handle nil case untuk reportData seperti legacy
	if response.ReportData == nil {
		response.ReportData = []interface{}{} // Default empty array like legacy
	}

	// Convert ReportType
	if report.ReportType != nil {
		response.ReportType = &ReportTypeInfo{
			ID:   report.ReportType.ID.Hex(),
			Name: report.ReportType.Name,
		}
	}

	// Convert Company
	if report.Company != nil {
		response.Company = &CompanyInfo{
			ID:             report.Company.ID.Hex(),
			Name:           report.Company.Name,
			ProfilePicture: report.Company.ProfilePicture,
			CreatedAt:      report.Company.CreatedAt,
			UpdatedAt:      report.Company.UpdatedAt,
		}
	}

	// Convert CreatedBy
	if report.CreatedBy != nil {
		response.CreatedBy = &UserInfo{
			ID:        report.CreatedBy.ID.Hex(),
			Name:      report.CreatedBy.Name,
			Email:     report.CreatedBy.Email,
			Role:      string(report.CreatedBy.Role),
			CreatedAt: report.CreatedBy.CreatedAt,
			UpdatedAt: report.CreatedBy.UpdatedAt,
		}
	}

	// Convert UserAccess array - handle nil case
	if report.UserAccess != nil {
		response.UserAccess = make([]*UserInfo, len(report.UserAccess))
		for i, user := range report.UserAccess {
			response.UserAccess[i] = &UserInfo{
				ID:        user.ID.Hex(),
				Name:      user.Name,
				Email:     user.Email,
				Role:      string(user.Role),
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
		}
	} else {
		response.UserAccess = []*UserInfo{} // Empty array like legacy
	}

	return response
}

// ToReportResponseArray converts array of domain reports to response array
func ToReportResponseArray(reports []*domain.PopulatedReport) []*ReportResponse {
	responses := make([]*ReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = ToReportResponse(report)
	}
	return responses
}
