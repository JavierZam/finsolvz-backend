// internal/app/auth/service.go - VERIFY THIS FILE HAS CORRECT INTERFACE
package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/utils"
	"finsolvz-backend/internal/utils/errors"
)

// ✅ Make sure this interface uses the correct types from model.go
type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req ResetPasswordRequest) error
}

type service struct {
	userRepo     domain.UserRepository
	emailService utils.EmailService
}

func NewService(userRepo domain.UserRepository, emailService utils.EmailService) Service {
	return &service{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      domain.UserRole(req.Role),
		Company:   []primitive.ObjectID{}, // Empty array like in your data
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  ToUserInfo(user),
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  ToUserInfo(user),
	}, nil
}

func (s *service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("USER_NOT_FOUND", "User not found", 404, err, nil)
	}

	// Generate new random password
	newPassword, err := utils.GenerateRandomPassword()
	if err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user.ID, user); err != nil {
		return err
	}

	// Send email with new password
	if err := s.emailService.SendForgotPasswordEmail(user.Email, user.Name, newPassword); err != nil {
		return err
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	// Get user by reset token
	user, err := s.userRepo.GetByResetToken(ctx, req.Token)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update user password and clear reset token
	user.Password = hashedPassword
	user.ResetPasswordToken = nil
	user.ResetPasswordExpires = nil

	if err := s.userRepo.Update(ctx, user.ID, user); err != nil {
		return err
	}

	return nil
}