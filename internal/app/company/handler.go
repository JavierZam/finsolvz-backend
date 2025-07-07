package company

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"finsolvz-backend/internal/platform/http/middleware"
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

// RegisterRoutes registers company routes
func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	// Protected routes - require authentication
	protected := router.PathPrefix("").Subrouter()
	protected.Use(authMiddleware)

	// Company routes
	protected.HandleFunc("/api/company", h.GetCompanies).Methods("GET")
	protected.HandleFunc("/api/company", h.CreateCompany).Methods("POST")
	protected.HandleFunc("/api/user/companies", h.GetUserCompanies).Methods("GET")
	protected.HandleFunc("/api/company/{idOrName}", h.GetCompanyByIDOrName).Methods("GET")
	
	// Admin-only routes
	adminOnly := protected.PathPrefix("").Subrouter()
	adminOnly.Use(middleware.RequireRole("SUPER_ADMIN"))
	adminOnly.HandleFunc("/api/company/{id}", h.UpdateCompany).Methods("PUT")
	adminOnly.HandleFunc("/api/company/{id}", h.DeleteCompany).Methods("DELETE")
}

func (h *Handler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.GetCompanies(r.Context())
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// ✅ EXACT legacy format: return array directly
	utils.RespondJSON(w, http.StatusOK, companies)
}

func (h *Handler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var req CreateCompanyRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	company, err := h.service.CreateCompany(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// ✅ EXACT legacy format
	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Company created successfully",
		"company": company,
	})
}

func (h *Handler) GetUserCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.GetUserCompanies(r.Context())
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, companies)
}

func (h *Handler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateCompanyRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	company, err := h.service.UpdateCompany(r.Context(), id, req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// ✅ EXACT legacy format
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Success",
		"company": company,
	})
}

func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	deletedCompany, err := h.service.DeleteCompany(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// ✅ EXACT legacy format
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Company deleted successfully",
		"company": deletedCompany,
	})
}

func (h *Handler) GetCompanyByIDOrName(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idOrName := vars["idOrName"]

    var company *CompanyResponse
    var err error

    // Logic detection: 24 hex chars = ObjectID, else = Name
    if len(idOrName) == 24 && isHexString(idOrName) {
        company, err = h.service.GetCompanyByID(r.Context(), idOrName)
    } else {
        company, err = h.service.GetCompanyByName(r.Context(), idOrName)
    }

    if err != nil {
        utils.HandleHTTPError(w, err, r)
        return
    }

    utils.RespondJSON(w, http.StatusOK, company)
}

func isHexString(s string) bool {
    for _, char := range s {
        if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
            return false
        }
    }
    return true
}