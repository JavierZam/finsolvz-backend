package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/platform/http/middleware"
	"finsolvz-backend/internal/utils"
	"finsolvz-backend/internal/utils/errors"
)

type Service interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error)
	GetUsers(ctx context.Context) ([]*UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*UserResponse, error)
	GetLoginUser(ctx context.Context) (*UserResponse, error)
	UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error)
	DeleteUser(ctx context.Context, id string) (*UserResponse, error)  // ✅ Updated return type
	UpdateRole(ctx context.Context, req UpdateRoleRequest) (*UserResponse, error)
	ChangePassword(ctx context.Context, req ChangePasswordRequest) error
}

type service struct {
	userRepo domain.UserRepository
}

func NewService(userRepo domain.UserRepository) Service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("USER_ALREADY_EXISTS", "Email already registered", 409, nil, nil)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     domain.UserRole(req.Role),
		Company:  []primitive.ObjectID{},
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	response := ToUserResponse(user)
	return &response, nil
}

func (s *service) GetUsers(ctx context.Context) ([]*UserResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		response := ToUserResponse(user)
		responses[i] = &response
	}

	return responses, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID format", 400, err, nil)
	}

	user, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	response := ToUserResponse(user)
	return &response, nil
}

func (s *service) GetLoginUser(ctx context.Context) (*UserResponse, error) {
	// Get user from context
	userCtx, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, errors.New("USER_CONTEXT_MISSING", "User context not found", 401, nil, nil)
	}

	objectID, err := primitive.ObjectIDFromHex(userCtx.UserID)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID in context", 400, err, nil)
	}

	user, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	response := ToUserResponse(user)
	return &response, nil
}

func (s *service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID format", 400, err, nil)
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and already exists
	if req.Email != nil && *req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err == nil && existingUser != nil {
			return nil, ErrEmailAlreadyExists
		}
	}

	// Update fields
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = domain.UserRole(*req.Role)
	}
	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := s.userRepo.Update(ctx, objectID, user); err != nil {
		return nil, err
	}

	response := ToUserResponse(user)
	return &response, nil
}

func (s *service) DeleteUser(ctx context.Context, id string) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID format", 400, err, nil)
	}

	user, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Delete(ctx, objectID); err != nil {
		return nil, err
	}

	response := ToUserResponse(user)
	return &response, nil
}

func (s *service) UpdateRole(ctx context.Context, req UpdateRoleRequest) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID format", 400, err, nil)
	}

	user, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	user.Role = domain.UserRole(req.NewRole)

	if err := s.userRepo.Update(ctx, objectID, user); err != nil {
		return nil, err
	}

	response := ToUserResponse(user)
	return &response, nil
}

func (s *service) ChangePassword(ctx context.Context, req ChangePasswordRequest) error {
	// Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		return ErrPasswordMismatch
	}

	// Get user from context
	userCtx, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return errors.New("USER_CONTEXT_MISSING", "User context not found", 401, nil, nil)
	}

	objectID, err := primitive.ObjectIDFromHex(userCtx.UserID)
	if err != nil {
		return errors.New("INVALID_USER_ID", "Invalid user ID in context", 400, err, nil)
	}

	user, err := s.userRepo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, objectID, user)
}
