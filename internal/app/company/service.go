package company

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/platform/http/middleware"
	"finsolvz-backend/internal/utils/errors"
)

type Service interface {
	CreateCompany(ctx context.Context, req CreateCompanyRequest) (*CompanyResponse, error)
	GetCompanies(ctx context.Context) ([]*CompanyResponse, error)
	GetCompanyByID(ctx context.Context, id string) (*CompanyResponse, error)
	GetCompanyByName(ctx context.Context, name string) (*CompanyResponse, error)
	GetUserCompanies(ctx context.Context) ([]*CompanyResponse, error)
	UpdateCompany(ctx context.Context, id string, req UpdateCompanyRequest) (*CompanyResponse, error)
	DeleteCompany(ctx context.Context, id string) (*CompanyResponse, error)
}

type service struct {
	companyRepo domain.CompanyRepository
	userRepo    domain.UserRepository
}

func NewService(companyRepo domain.CompanyRepository, userRepo domain.UserRepository) Service {
	return &service{
		companyRepo: companyRepo,
		userRepo:    userRepo,
	}
}

func (s *service) CreateCompany(ctx context.Context, req CreateCompanyRequest) (*CompanyResponse, error) {
	// Trim and validate name
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidCompanyName
	}

	// Check if company already exists
	existingCompany, err := s.companyRepo.GetByName(ctx, name)
	if err == nil && existingCompany != nil {
		return nil, ErrCompanyAlreadyExists
	}

	// Process user IDs
	var userIDs []primitive.ObjectID
	for _, userIDStr := range req.User {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return nil, ErrInvalidUserID
		}

		// Verify user exists
		_, err = s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, ErrUserNotFound
		}

		userIDs = append(userIDs, userID)
	}

	// Create company
	company := &domain.Company{
		Name:           name,
		ProfilePicture: req.ProfilePicture,
		User:           userIDs,
	}

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}

	// Get populated user data for response
	users, err := s.getUsersByIDs(ctx, userIDs)
	if err != nil {
		// If user fetching fails, return basic response without user details
		response := ToCompanyResponse(company)
		return &response, nil
	}

	response := ToCompanyResponseWithUsers(company, users)
	return &response, nil
}

func (s *service) GetCompanies(ctx context.Context) ([]*CompanyResponse, error) {
	companies, err := s.companyRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*CompanyResponse, len(companies))
	for i, company := range companies {
		// Get populated user data
		users, err := s.getUsersByIDs(ctx, company.User)
		if err != nil {
			// If user fetching fails, return basic response
			response := ToCompanyResponse(company)
			responses[i] = &response
		} else {
			response := ToCompanyResponseWithUsers(company, users)
			responses[i] = &response
		}
	}

	return responses, nil
}

func (s *service) GetCompanyByID(ctx context.Context, id string) (*CompanyResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_COMPANY_ID", "Invalid company ID format", 400, err, nil)
	}

	company, err := s.companyRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Handle profile picture URL processing (legacy compatibility)
	if company.ProfilePicture != nil && !strings.HasPrefix(*company.ProfilePicture, "http") {
		// Add server URL prefix for legacy compatibility
		fullURL := "http://152.42.172.219:8787" + *company.ProfilePicture
		company.ProfilePicture = &fullURL
	}

	users, err := s.getUsersByIDs(ctx, company.User)
	if err != nil {
		response := ToCompanyResponse(company)
		return &response, nil
	}

	response := ToCompanyResponseWithUsers(company, users)
	return &response, nil
}

func (s *service) GetCompanyByName(ctx context.Context, name string) (*CompanyResponse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidCompanyName
	}

	company, err := s.companyRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Handle profile picture URL processing (legacy compatibility)
	if company.ProfilePicture != nil && !strings.HasPrefix(*company.ProfilePicture, "http") {
		fullURL := "http://152.42.172.219:8787" + *company.ProfilePicture
		company.ProfilePicture = &fullURL
	}

	users, err := s.getUsersByIDs(ctx, company.User)
	if err != nil {
		response := ToCompanyResponse(company)
		return &response, nil
	}

	response := ToCompanyResponseWithUsers(company, users)
	return &response, nil
}

func (s *service) GetUserCompanies(ctx context.Context) ([]*CompanyResponse, error) {
	// Get current user from context
	userCtx, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, errors.New("USER_CONTEXT_MISSING", "User context not found", 401, nil, nil)
	}

	userID, err := primitive.ObjectIDFromHex(userCtx.UserID)
	if err != nil {
		return nil, errors.New("INVALID_USER_ID", "Invalid user ID in context", 400, err, nil)
	}

	companies, err := s.companyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*CompanyResponse, len(companies))
	for i, company := range companies {
		response := ToCompanyResponse(company)
		responses[i] = &response
	}

	return responses, nil
}

func (s *service) UpdateCompany(ctx context.Context, id string, req UpdateCompanyRequest) (*CompanyResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_COMPANY_ID", "Invalid company ID format", 400, err, nil)
	}

	// Get existing company
	company, err := s.companyRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Update name if provided
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrInvalidCompanyName
		}

		// Check if new name conflicts with existing company (exclude current one)
		if name != company.Name {
			existingCompany, err := s.companyRepo.GetByName(ctx, name)
			if err == nil && existingCompany != nil {
				return nil, ErrCompanyAlreadyExists
			}
		}
		company.Name = name
	}

	// Update profile picture if provided
	if req.ProfilePicture != nil {
		company.ProfilePicture = req.ProfilePicture
	}

	// Process user IDs if provided
	if req.User != nil {
		var userIDs []primitive.ObjectID
		for _, userIDStr := range req.User {
			userID, err := primitive.ObjectIDFromHex(userIDStr)
			if err != nil {
				return nil, ErrInvalidUserID
			}

			// Verify user exists
			_, err = s.userRepo.GetByID(ctx, userID)
			if err != nil {
				return nil, ErrUserNotFound
			}

			userIDs = append(userIDs, userID)
		}
		company.User = userIDs
	}

	if err := s.companyRepo.Update(ctx, objectID, company); err != nil {
		return nil, err
	}

	users, err := s.getUsersByIDs(ctx, company.User)
	if err != nil {
		response := ToCompanyResponse(company)
		return &response, nil
	}

	response := ToCompanyResponseWithUsers(company, users)
	return &response, nil
}

func (s *service) DeleteCompany(ctx context.Context, id string) (*CompanyResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("INVALID_COMPANY_ID", "Invalid company ID format", 400, err, nil)
	}

	// Get company data before deletion
	company, err := s.companyRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if err := s.companyRepo.Delete(ctx, objectID); err != nil {
		return nil, err
	}

	response := ToCompanyResponse(company)
	return &response, nil
}

// Helper function to get users by IDs
func (s *service) getUsersByIDs(ctx context.Context, userIDs []primitive.ObjectID) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(userIDs))
	for _, userID := range userIDs {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err == nil {
			users = append(users, user)
		}
		// Continue even if some users are not found (soft error handling)
	}
	return users, nil
}