package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/utils"
	"finsolvz-backend/internal/utils/errors"
)

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
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      domain.UserRole(req.Role),
		Company:   []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

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
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

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
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("USER_NOT_FOUND", "User not found", 404, err, nil)
	}

	newPassword, err := utils.GenerateRandomPassword()
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user.ID, user); err != nil {
		return err
	}

	if err := s.emailService.SendForgotPasswordEmail(user.Email, user.Name, newPassword); err != nil {
		return err
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	user, err := s.userRepo.GetByResetToken(ctx, req.Token)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Clear reset token after successful password change
	user.Password = hashedPassword
	user.ResetPasswordToken = nil
	user.ResetPasswordExpires = nil

	if err := s.userRepo.Update(ctx, user.ID, user); err != nil {
		return err
	}

	return nil
}
