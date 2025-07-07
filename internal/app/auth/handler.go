package auth

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

// RegisterRoutes registers auth routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Public routes (no authentication required)
	router.HandleFunc("/api/login", h.Login).Methods("POST")
	router.HandleFunc("/api/forgot-password", h.ForgotPassword).Methods("POST")
	router.HandleFunc("/api/reset-password", h.ResetPassword).Methods("POST")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, utils.ErrBadRequest, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	response, err := h.service.Login(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": response.Token,
	})
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, utils.ErrBadRequest, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	if err := h.service.ForgotPassword(r.Context(), req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "New password has been sent to your email",
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, utils.ErrBadRequest, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	if err := h.service.ResetPassword(r.Context(), req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password successfully reset",
	})
}