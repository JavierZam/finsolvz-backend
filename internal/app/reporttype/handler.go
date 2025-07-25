package reporttype

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"finsolvz-backend/internal/utils"
)

type Handler struct {
	service   Service
	validator *validator.Validate
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service:   service,
		validator: validator.New(),
	}
}

// RegisterRoutes registers report type routes
func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	protected := router.PathPrefix("").Subrouter()
	protected.Use(authMiddleware)

	protected.HandleFunc("/api/reportTypes", h.GetReportTypes).Methods("GET")
	protected.HandleFunc("/api/reportTypes", h.CreateReportType).Methods("POST")
	protected.HandleFunc("/api/reportTypes/{id}", h.UpdateReportType).Methods("PUT")
	protected.HandleFunc("/api/reportTypes/{id}", h.DeleteReportType).Methods("DELETE")
	protected.HandleFunc("/api/reportTypes/{idOrName}", h.GetReportTypeByIDOrName).Methods("GET")
}

func (h *Handler) GetReportTypes(w http.ResponseWriter, r *http.Request) {
	reportTypes, err := h.service.GetReportTypes(r.Context())
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reportTypes)
}

// GetReportTypeByIDOrName retrieves a report type by ID or name
func (h *Handler) GetReportTypeByIDOrName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idOrName := vars["idOrName"]

	var reportType *ReportTypeResponse
	var err error

	// Check if parameter is ObjectID format (24 hex characters) or name
	if len(idOrName) == 24 {
		reportType, err = h.service.GetReportTypeByID(r.Context(), idOrName)
	} else {
		reportType, err = h.service.GetReportTypeByName(r.Context(), idOrName)
	}

	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reportType)
}

func (h *Handler) CreateReportType(w http.ResponseWriter, r *http.Request) {
	var req CreateReportTypeRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	reportType, err := h.service.CreateReportType(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "Report type added successfully",
		"reportType": reportType,
	})
}

func (h *Handler) UpdateReportType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateReportTypeRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	reportType, err := h.service.UpdateReportType(r.Context(), id, req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message":    "Report Type updated successfully",
		"reportType": reportType,
	})
}

func (h *Handler) DeleteReportType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteReportType(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
