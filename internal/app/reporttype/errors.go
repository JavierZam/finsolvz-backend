package reporttype

import (
	"finsolvz-backend/internal/utils/errors"
	"net/http"
)

var (
	ErrReportTypeNotFound      = errors.New("REPORT_TYPE_NOT_FOUND", "Report type not found", http.StatusNotFound, nil, nil)
	ErrReportTypeAlreadyExists = errors.New("REPORT_TYPE_ALREADY_EXISTS", "Report type name already exists", http.StatusConflict, nil, nil)
	ErrInvalidReportTypeName   = errors.New("INVALID_REPORT_TYPE_NAME", "Report type name is invalid", http.StatusBadRequest, nil, nil)
)