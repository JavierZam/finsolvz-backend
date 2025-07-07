package reporttype

import (
	"finsolvz-backend/internal/domain"
)

// Request DTOs
type CreateReportTypeRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateReportTypeRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// Response DTOs - exact legacy format
type ReportTypeResponse struct {
	ID   string `json:"id"`   // ✅ Changed to "id" exactly like legacy Mongoose
	Name string `json:"name"`
}

// Helper to convert domain.ReportType to ReportTypeResponse
func ToReportTypeResponse(reportType *domain.ReportType) ReportTypeResponse {
	return ReportTypeResponse{
		ID:   reportType.ID.Hex(),
		Name: reportType.Name,
	}
}