package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"finsolvz-backend/internal/utils/errors"
	"finsolvz-backend/internal/utils/log"
)

// ErrorResponse struct untuk respons error yang konsisten ke klien.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// RespondJSON menulis respons JSON ke klien dengan status code dan data yang diberikan.
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Errorf(context.Background(), "Failed to write JSON response: %v", err)
			http.Error(w, `{"code":"INTERNAL_SERVER_ERROR","message":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// HandleHTTPError memetakan AppError ke respons HTTP yang sesuai.
func HandleHTTPError(w http.ResponseWriter, err error, r *http.Request) {
	appErr, ok := err.(errors.AppError)
	if !ok {
		log.Errorf(r.Context(), "Unhandled error: %v", err)
		RespondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Code:    errors.ErrInternalServer.Code(),
			Message: errors.ErrInternalServer.Message(),
			Details: err.Error(),
		})
		return
	}

	if appErr.Status() >= http.StatusInternalServerError {
		log.Errorf(r.Context(), "Server error occurred: %v", appErr)
		detailsMessage := appErr.Message()
		if os.Getenv("APP_ENV") == "development" {
			if unwrappedErr := appErr.Unwrap(); unwrappedErr != nil {
				detailsMessage = unwrappedErr.Error()
			} else {
				detailsMessage = appErr.Error()
			}
		}
		RespondJSON(w, appErr.Status(), ErrorResponse{
			Code:    appErr.Code(),
			Message: appErr.Message(),
			Details: detailsMessage,
		})
	} else {
		log.Warnf(r.Context(), "Client-side error: %v", appErr)
		RespondJSON(w, appErr.Status(), ErrorResponse{
			Code:    appErr.Code(),
			Message: appErr.Message(),
			Details: formatErrorDetails(appErr.Details()),
		})
	}
}

func formatErrorDetails(details map[string]interface{}) string {
	if details == nil {
		return ""
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return ""
	}
	return string(detailsJSON)
}
