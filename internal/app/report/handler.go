package report

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

// RegisterRoutes registers report routes
func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	protected := router.PathPrefix("").Subrouter()
	protected.Use(authMiddleware)

	protected.HandleFunc("/api/reports", h.CreateReport).Methods("POST")
	protected.HandleFunc("/api/reports/{id}", h.UpdateReport).Methods("PUT")
	protected.HandleFunc("/api/reports/{id}", h.DeleteReport).Methods("DELETE")

	protected.HandleFunc("/api/reports", h.GetReports).Methods("GET")
	protected.HandleFunc("/api/reports/{id}", h.GetReportByID).Methods("GET")
	protected.HandleFunc("/api/reports/name/{name}", h.GetReportByName).Methods("GET")
	protected.HandleFunc("/api/reports/company/{companyId}", h.GetReportsByCompany).Methods("GET")
	protected.HandleFunc("/api/reports/companies", h.GetReportsByCompanies).Methods("POST")

	protected.HandleFunc("/api/reports/reportType/{reportType}", h.GetReportsByReportType).Methods("GET")
	protected.HandleFunc("/api/reports/userAccess/{id}", h.GetReportsByUserAccess).Methods("GET")
	protected.HandleFunc("/api/reports/createdBy/{id}", h.GetReportsByCreatedBy).Methods("GET")
}

func (h *Handler) CreateReport(w http.ResponseWriter, r *http.Request) {
	var req CreateReportRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	report, err := h.service.CreateReport(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, report)
}

func (h *Handler) UpdateReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateReportRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	report, err := h.service.UpdateReport(r.Context(), id, req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, report)
}

func (h *Handler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteReport(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Report deleted successfully",
	})
}

func (h *Handler) GetReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.service.GetReports(r.Context())
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reports)
}

func (h *Handler) GetReportByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	report, err := h.service.GetReportByID(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, report)
}

func (h *Handler) GetReportByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	report, err := h.service.GetReportByName(r.Context(), name)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, report)
}

func (h *Handler) GetReportsByCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyId := vars["companyId"]

	reports, err := h.service.GetReportsByCompany(r.Context(), companyId)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reports)
}

func (h *Handler) GetReportsByCompanies(w http.ResponseWriter, r *http.Request) {
	var req GetReportsByCompaniesRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	reports, err := h.service.GetReportsByCompanies(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reports)
}

func (h *Handler) GetReportsByReportType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportType := vars["reportType"]

	reports, err := h.service.GetReportsByReportType(r.Context(), reportType)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reports)
}

func (h *Handler) GetReportsByUserAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reports, err := h.service.GetReportsByUserAccess(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reports)
}

func (h *Handler) GetReportsByCreatedBy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	reports, err := h.service.GetReportsByCreatedBy(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, reports)
}
