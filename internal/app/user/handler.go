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
	authService auth.Service  // ✅ Add auth service for register
	validator   *validator.Validate
}

// ✅ Updated constructor to include auth service
func NewHandler(service Service, authService auth.Service) *Handler {
	return &Handler{
		service:     service,
		authService: authService,
		validator:   validator.New(),
	}
}

// RegisterRoutes registers user routes
func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	// Protected routes - require authentication
	protected := router.PathPrefix("").Subrouter()
	protected.Use(authMiddleware)

	// User management routes
	protected.HandleFunc("/api/users", h.GetUsers).Methods("GET")
	protected.HandleFunc("/api/users/{id}", h.GetUserByID).Methods("GET")
	protected.HandleFunc("/api/loginUser", h.GetLoginUser).Methods("GET")
	protected.HandleFunc("/api/users/{id}", h.UpdateUser).Methods("PUT")
	protected.HandleFunc("/api/users/{id}", h.DeleteUser).Methods("DELETE")

	// Role management - SUPER_ADMIN only
	superAdminOnly := protected.PathPrefix("").Subrouter()
	superAdminOnly.Use(middleware.RequireRole("SUPER_ADMIN"))
	superAdminOnly.HandleFunc("/api/register", h.Register).Methods("POST")  // ✅ ADD REGISTER HERE
	superAdminOnly.HandleFunc("/api/updateRole", h.UpdateRole).Methods("PUT")

	// Password change
	protected.HandleFunc("/api/change-password", h.ChangePassword).Methods("PATCH")
}

// ✅ NEW: Register method using auth service (SUPER_ADMIN only)
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest  // Use auth.RegisterRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	// Use auth service to register (includes token generation)
	response, err := h.authService.Register(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// Return same format as legacy Node.js
	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Success",
		"newUser": response.User,  // Only return user info, not token
	})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	response, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Success",
		"newUser": response,
	})
}

// Rest of handler methods remain the same...
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Check if user has permission (SUPER_ADMIN or ADMIN)
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

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if user can update (SUPER_ADMIN only for role changes)
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.HandleHTTPError(w, utils.ErrUnauthorized, r)
		return
	}

	var req UpdateUserRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.HandleValidationError(w, err, r)
		return
	}

	// Only SUPER_ADMIN can change roles
	if req.Role != nil && userCtx.Role != "SUPER_ADMIN" {
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

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Only SUPER_ADMIN can delete users
	userCtx, ok := middleware.GetUserFromContext(r.Context())
	if !ok || userCtx.Role != "SUPER_ADMIN" {
		utils.HandleHTTPError(w, utils.ErrForbidden, r)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	// Get deleted user data from service
	deletedUser, err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		utils.HandleHTTPError(w, err, r)
		return
	}

	// Return same format as legacy Node.js
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Success",
		"user":    deletedUser,
	})
}

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