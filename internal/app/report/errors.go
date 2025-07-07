package report

import (
	"finsolvz-backend/internal/utils/errors"
	"net/http"
)

var (
	ErrReportNotFound        = errors.New("REPORT_NOT_FOUND", "Report not found", http.StatusNotFound, nil, nil)
	ErrReportAlreadyExists   = errors.New("REPORT_ALREADY_EXISTS", "Report with this name already exists", http.StatusConflict, nil, nil)
	ErrInvalidReportName     = errors.New("INVALID_REPORT_NAME", "Report name is invalid", http.StatusBadRequest, nil, nil)
	ErrInvalidReportID       = errors.New("INVALID_REPORT_ID", "Invalid report ID format", http.StatusBadRequest, nil, nil)
	ErrInvalidReportTypeID   = errors.New("INVALID_REPORT_TYPE_ID", "Invalid report type ID format", http.StatusBadRequest, nil, nil)
	ErrInvalidCompanyID      = errors.New("INVALID_COMPANY_ID", "Invalid company ID format", http.StatusBadRequest, nil, nil)
	ErrInvalidUserID         = errors.New("INVALID_USER_ID", "Invalid user ID format", http.StatusBadRequest, nil, nil)
	ErrInvalidYear           = errors.New("INVALID_YEAR", "Year format is invalid", http.StatusBadRequest, nil, nil)
	ErrInsufficientCompanies = errors.New("INSUFFICIENT_COMPANIES", "Need 2 or more companies", http.StatusBadRequest, nil, nil)
	ErrReportDataProcessing  = errors.New("REPORT_DATA_PROCESSING_ERROR", "Failed to process report data", http.StatusInternalServerError, nil, nil)
	ErrGeminiProcessing      = errors.New("GEMINI_PROCESSING_ERROR", "Failed to process data with AI", http.StatusInternalServerError, nil, nil)
)
