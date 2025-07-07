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
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidCompanyName
	}

	existingCompany, err := s.companyRepo.GetByName(ctx, name)
	if err == nil && existingCompany != nil {
		return nil, ErrCompanyAlreadyExists
	}

	var userIDs []primitive.ObjectID
	for _, userIDStr := range req.User {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return nil, ErrInvalidUserID
		}

		_, err = s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, ErrUserNotFound
		}

		userIDs = append(userIDs, userID)
	}

	company := &domain.Company{
		Name:           name,
		ProfilePicture: req.ProfilePicture,
		User:           userIDs,
	}

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, err
	}

	users, err := s.getUsersByIDs(ctx, userIDs)
	if err != nil {
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
		users, err := s.getUsersByIDs(ctx, company.User)
		if err != nil {
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

	// Convert relative URLs to absolute URLs
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

	company, err := s.companyRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrInvalidCompanyName
		}

		// Check name uniqueness when being changed
		if name != company.Name {
			existingCompany, err := s.companyRepo.GetByName(ctx, name)
			if err == nil && existingCompany != nil {
				return nil, ErrCompanyAlreadyExists
			}
		}
		company.Name = name
	}

	if req.ProfilePicture != nil {
		company.ProfilePicture = req.ProfilePicture
	}

	if req.User != nil {
		var userIDs []primitive.ObjectID
		for _, userIDStr := range req.User {
			userID, err := primitive.ObjectIDFromHex(userIDStr)
			if err != nil {
				return nil, ErrInvalidUserID
			}

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

// getUsersByIDs retrieves users by their IDs, skipping any that are not found
func (s *service) getUsersByIDs(ctx context.Context, userIDs []primitive.ObjectID) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(userIDs))
	for _, userID := range userIDs {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

func (s *service) GetCompanyByName(ctx context.Context, name string) (*CompanyResponse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidCompanyName
	}

	// Try exact match first, then flexible search
	company, err := s.companyRepo.GetByName(ctx, name)
	if err == nil {
		return s.buildCompanyResponse(ctx, company)
	}

	// Fallback to flexible search if exact match fails
	companies, searchErr := s.companyRepo.SearchByName(ctx, name)
	if searchErr != nil || len(companies) == 0 {
		return nil, ErrCompanyNotFound
	}

	return s.buildCompanyResponse(ctx, companies[0])
}

// buildCompanyResponse creates a company response with populated users and processed URLs
func (s *service) buildCompanyResponse(ctx context.Context, company *domain.Company) (*CompanyResponse, error) {
	// Convert relative URLs to absolute URLs
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