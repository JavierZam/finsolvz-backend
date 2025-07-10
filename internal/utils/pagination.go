package utils

import (
	"net/http"
	"strconv"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page    int `json:"page"`
	Limit   int `json:"limit"`
	Skip    int `json:"skip"`
	Total   int `json:"total,omitempty"`
}

// PaginatedResponse wraps data with pagination info
type PaginatedResponse struct {
	Data       interface{}       `json:"data"`
	Pagination PaginationParams  `json:"pagination"`
}

// GetPaginationParams extracts pagination parameters from request
func GetPaginationParams(r *http.Request) PaginationParams {
	page := 1
	limit := 10 // Default limit to reduce data transfer

	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	skip := (page - 1) * limit

	return PaginationParams{
		Page:  page,
		Limit: limit,
		Skip:  skip,
	}
}

// CreatePaginatedResponse creates a paginated response
func CreatePaginatedResponse(data interface{}, pagination PaginationParams) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	}
}