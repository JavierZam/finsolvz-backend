package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"finsolvz-backend/internal/app/auth"
	"finsolvz-backend/internal/platform/http/middleware"
	"finsolvz-backend/internal/utils"
)

type Handler struct {
	service     Service
	authService auth.Service
	validator   *validator.Validate
}

func NewHandler(service Service, authService auth.Service) *Handler {
	return &Handler{
		service:     service,
		authService: authService,
		validator:   validator.New(),
	}
}

// ✅ FIXED: RegisterRoutes - No role-based middleware, auth check in controller like legacy
func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	// Protected routes - require authentication only
	protected := router.PathPrefix("").Subrouter()
	protected.Use(authMiddleware)

	// ✅ All routes use same pattern as legacy - auth check in controller if needed
	protected.HandleFunc("/api/users", h.GetUsers).Methods("GET")
	protected.HandleFunc("/api/users/{id}", h.GetUserByID).Methods("GET")
	protected.HandleFunc("/api/loginUser", h.GetLoginUser).Methods("GET")
	protected.HandleFunc("/api/users/{id}", h.UpdateUser).Methods("PUT")
	protected.HandleFunc("/api/users/{id}", h.DeleteUser).Methods("DELETE")

	// ✅ FIXED: Register route - no role middleware, check in controller like legacy
	protected.HandleFunc("/api/register", h.Register).Methods("POST")
	protected.HandleFunc("/api/updateRole", h.UpdateRole).Methods("PUT")

	// Password change
	protected.HandleFunc("/api/change-password", h.ChangePassword).Methods("PATCH")
}

// ✅ FIXED: Register method - authorization check in controller like legacy
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	// ✅ EXACT legacy authorization check pattern
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userCtx.Role != "SUPER_ADMIN" {
		utils.HandleHTTPError(w, utils.ErrForbidden, r)
		return
	}

	// Use auth service to register (includes token generation)
	response, err := h.authService.Register(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// ✅ EXACT legacy format
	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Success",
		"newUser": response.User, // Only return user info, not token
	})
}

// ✅ FIXED: GetUsers - authorization check in controller like legacy
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// ✅ EXACT legacy authorization check pattern
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || (userCtx.Role != "SUPER_ADMIN" && userCtx.Role != "ADMIN") {
		utils.HandleHTTPError(w, utils.ErrForbidden, r)
		return
	}

	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *Handler) GetLoginUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetLoginUser(r.Context())
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

// ✅ FIXED: UpdateUser - authorization check in controller like legacy
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateUserRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	// ✅ EXACT legacy authorization check pattern - only SUPER_ADMIN can update
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userCtx.Role != "SUPER_ADMIN" {
		utils.HandleHTTPError(w, utils.ErrForbidden, r)
		return
	}

	response, err := h.service.UpdateUser(r.Context(), id, req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "User updated",
		"updatedUser": response,
	})
}

// ✅ FIXED: DeleteUser - authorization check in controller like legacy
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// ✅ EXACT legacy authorization check pattern
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userCtx.Role != "SUPER_ADMIN" {
		utils.HandleHTTPError(w, utils.ErrForbidden, r)
		return
	}

	deletedUser, err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// ✅ EXACT legacy format
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Success",
		"user":    deletedUser,
	})
}

// ✅ FIXED: UpdateRole - authorization check in controller like legacy
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var req UpdateRoleRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	// ✅ EXACT legacy authorization check pattern
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userCtx.Role != "SUPER_ADMIN" {
		utils.HandleHTTPError(w, utils.ErrForbidden, r)
		return
	}

	response, err := h.service.UpdateRole(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Success",
		"user":    response,
	})
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	if err := h.service.ChangePassword(r.Context(), req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Password successfully changed",
	})
}
