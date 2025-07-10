package company

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
)

// Mock repositories
type mockCompanyRepository struct {
	companies []domain.Company
}

func (m *mockCompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	company.ID = primitive.NewObjectID()
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	m.companies = append(m.companies, *company)
	return nil
}

func (m *mockCompanyRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Company, error) {
	for i := range m.companies {
		if m.companies[i].ID == id {
			return &m.companies[i], nil
		}
	}
	return nil, domain.ErrCompanyNotFound
}

func (m *mockCompanyRepository) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	for i := range m.companies {
		if m.companies[i].Name == name {
			return &m.companies[i], nil
		}
	}
	return nil, domain.ErrCompanyNotFound
}

func (m *mockCompanyRepository) GetAll(ctx context.Context) ([]*domain.Company, error) {
	var result []*domain.Company
	for i := range m.companies {
		result = append(result, &m.companies[i])
	}
	return result, nil
}

func (m *mockCompanyRepository) GetByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Company, error) {
	var result []*domain.Company
	for i := range m.companies {
		for _, uid := range m.companies[i].User {
			if uid == userID {
				result = append(result, &m.companies[i])
				break
			}
		}
	}
	return result, nil
}

func (m *mockCompanyRepository) Update(ctx context.Context, id primitive.ObjectID, company *domain.Company) error {
	for i := range m.companies {
		if m.companies[i].ID == id {
			company.UpdatedAt = time.Now()
			m.companies[i] = *company
			return nil
		}
	}
	return domain.ErrCompanyNotFound
}

func (m *mockCompanyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	for i := range m.companies {
		if m.companies[i].ID == id {
			m.companies = append(m.companies[:i], m.companies[i+1:]...)
			return nil
		}
	}
	return domain.ErrCompanyNotFound
}

func (m *mockCompanyRepository) SearchByName(ctx context.Context, name string) ([]*domain.Company, error) {
	var result []*domain.Company
	for i := range m.companies {
		if m.companies[i].Name == name {
			result = append(result, &m.companies[i])
		}
	}
	if len(result) == 0 {
		return nil, domain.ErrCompanyNotFound
	}
	return result, nil
}

type mockUserRepository struct {
	users []domain.User
}

func (m *mockUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	for i := range m.users {
		if m.users[i].ID == id {
			return &m.users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error { return nil }
func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) { return nil, nil }
func (m *mockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) { return nil, nil }
func (m *mockUserRepository) Update(ctx context.Context, id primitive.ObjectID, user *domain.User) error { return nil }
func (m *mockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error { return nil }
func (m *mockUserRepository) SetResetToken(ctx context.Context, email, token string, expires time.Time) error { return nil }
func (m *mockUserRepository) GetByResetToken(ctx context.Context, token string) (*domain.User, error) { return nil, nil }

func TestCompanyService_CreateCompany(t *testing.T) {
	tests := []struct {
		name        string
		request     CreateCompanyRequest
		expectError bool
		setupData   func(*mockCompanyRepository)
	}{
		{
			name: "Valid company creation",
			request: CreateCompanyRequest{
				Name: "Test Company",
				User: []string{primitive.NewObjectID().Hex()},
			},
			expectError: false,
			setupData:   func(repo *mockCompanyRepository) {},
		},
		{
			name: "Empty company name",
			request: CreateCompanyRequest{
				Name: "",
				User: []string{primitive.NewObjectID().Hex()},
			},
			expectError: true,
			setupData:   func(repo *mockCompanyRepository) {},
		},
		{
			name: "Duplicate company name",
			request: CreateCompanyRequest{
				Name: "Existing Company",
				User: []string{primitive.NewObjectID().Hex()},
			},
			expectError: true,
			setupData: func(repo *mockCompanyRepository) {
				existingCompany := domain.Company{
					ID:   primitive.NewObjectID(),
					Name: "Existing Company",
					User: []primitive.ObjectID{primitive.NewObjectID()},
				}
				repo.companies = append(repo.companies, existingCompany)
			},
		},
		{
			name: "Invalid user ID format",
			request: CreateCompanyRequest{
				Name: "Test Company",
				User: []string{"invalid-id"},
			},
			expectError: true,
			setupData:   func(repo *mockCompanyRepository) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockCompanyRepo := &mockCompanyRepository{}
			mockUserRepo := &mockUserRepository{}
			tt.setupData(mockCompanyRepo)
			
			service := NewService(mockCompanyRepo, mockUserRepo)

			// Execute
			response, err := service.CreateCompany(context.Background(), tt.request)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if response != nil {
					t.Errorf("Expected nil response when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if response != nil && response.Name != tt.request.Name {
					t.Errorf("Expected name %s, got %s", tt.request.Name, response.Name)
				}
			}
		})
	}
}

func TestCompanyService_GetCompanies(t *testing.T) {
	// Setup
	mockCompanyRepo := &mockCompanyRepository{}
	mockUserRepo := &mockUserRepository{}
	
	// Add test data
	userID := primitive.NewObjectID()
	testUser := domain.User{
		ID:   userID,
		Name: "Test User",
		Email: "test@example.com",
	}
	mockUserRepo.users = append(mockUserRepo.users, testUser)
	
	testCompany := domain.Company{
		ID:   primitive.NewObjectID(),
		Name: "Test Company",
		User: []primitive.ObjectID{userID},
	}
	mockCompanyRepo.companies = append(mockCompanyRepo.companies, testCompany)
	
	service := NewService(mockCompanyRepo, mockUserRepo)

	// Execute
	companies, err := service.GetCompanies(context.Background())

	// Assert
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	
	if len(companies) != 1 {
		t.Errorf("Expected 1 company, got %d", len(companies))
	}
	
	if len(companies) > 0 && companies[0].Name != "Test Company" {
		t.Errorf("Expected company name 'Test Company', got %s", companies[0].Name)
	}
}

func TestCompanyService_GetCompanyByID(t *testing.T) {
	// Setup
	mockCompanyRepo := &mockCompanyRepository{}
	mockUserRepo := &mockUserRepository{}
	
	companyID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	
	testUser := domain.User{
		ID:   userID,
		Name: "Test User",
		Email: "test@example.com",
	}
	mockUserRepo.users = append(mockUserRepo.users, testUser)
	
	testCompany := domain.Company{
		ID:   companyID,
		Name: "Test Company",
		User: []primitive.ObjectID{userID},
	}
	mockCompanyRepo.companies = append(mockCompanyRepo.companies, testCompany)
	
	service := NewService(mockCompanyRepo, mockUserRepo)

	tests := []struct {
		name        string
		companyID   string
		expectError bool
	}{
		{
			name:        "Valid company ID",
			companyID:   companyID.Hex(),
			expectError: false,
		},
		{
			name:        "Invalid company ID format",
			companyID:   "invalid-id",
			expectError: true,
		},
		{
			name:        "Non-existent company ID",
			companyID:   primitive.NewObjectID().Hex(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			company, err := service.GetCompanyByID(context.Background(), tt.companyID)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if company == nil {
					t.Errorf("Expected company but got nil")
				}
				if company != nil && company.Name != "Test Company" {
					t.Errorf("Expected company name 'Test Company', got %s", company.Name)
				}
			}
		})
	}
}

// Performance test for GetCompanies with caching
func TestCompanyService_GetCompaniesPerformance(t *testing.T) {
	// Setup
	mockCompanyRepo := &mockCompanyRepository{}
	mockUserRepo := &mockUserRepository{}
	
	// Add multiple companies for performance testing
	userID := primitive.NewObjectID()
	testUser := domain.User{
		ID:   userID,
		Name: "Test User",
		Email: "test@example.com",
	}
	mockUserRepo.users = append(mockUserRepo.users, testUser)
	
	// Add 50 companies
	for i := 0; i < 50; i++ {
		company := domain.Company{
			ID:   primitive.NewObjectID(),
			Name: "Test Company " + string(rune(i)),
			User: []primitive.ObjectID{userID},
		}
		mockCompanyRepo.companies = append(mockCompanyRepo.companies, company)
	}
	
	service := NewService(mockCompanyRepo, mockUserRepo)

	// First call (no cache)
	start := time.Now()
	companies1, err := service.GetCompanies(context.Background())
	firstCallDuration := time.Since(start)
	
	if err != nil {
		t.Fatalf("First call failed: %v", err)
	}

	// Second call (should use cache)
	start = time.Now()
	companies2, err := service.GetCompanies(context.Background())
	secondCallDuration := time.Since(start)
	
	if err != nil {
		t.Fatalf("Second call failed: %v", err)
	}

	// Assert
	if len(companies1) != len(companies2) {
		t.Errorf("Cache returned different number of companies")
	}
	
	if len(companies1) != 50 {
		t.Errorf("Expected 50 companies, got %d", len(companies1))
	}

	// Second call should be faster (cached)
	if secondCallDuration > firstCallDuration {
		t.Logf("Warning: Cached call took longer than first call. First: %v, Second: %v", 
			firstCallDuration, secondCallDuration)
	}
	
	t.Logf("Performance test - First call: %v, Cached call: %v", 
		firstCallDuration, secondCallDuration)

	// Both calls should be reasonably fast
	if firstCallDuration > 100*time.Millisecond {
		t.Errorf("First call too slow: %v", firstCallDuration)
	}
	
	if secondCallDuration > 50*time.Millisecond {
		t.Errorf("Cached call too slow: %v", secondCallDuration)
	}
}