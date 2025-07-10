package auth

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"finsolvz-backend/internal/domain"
	"finsolvz-backend/internal/utils"
)

// Mock repository untuk testing
type mockUserRepository struct {
	users           []domain.User
	lastCreatedUser *domain.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users = append(m.users, *user)
	m.lastCreatedUser = user
	return nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	for i := range m.users {
		if m.users[i].Email == email {
			return &m.users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	for i := range m.users {
		if m.users[i].ID == id {
			return &m.users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	var result []*domain.User
	for i := range m.users {
		result = append(result, &m.users[i])
	}
	return result, nil
}

func (m *mockUserRepository) Update(ctx context.Context, id primitive.ObjectID, user *domain.User) error {
	for i := range m.users {
		if m.users[i].ID == id {
			user.UpdatedAt = time.Now()
			m.users[i] = *user
			return nil
		}
	}
	return domain.ErrUserNotFound
}

func (m *mockUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	for i := range m.users {
		if m.users[i].ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return domain.ErrUserNotFound
}

func (m *mockUserRepository) SetResetToken(ctx context.Context, email, token string, expires time.Time) error {
	for i := range m.users {
		if m.users[i].Email == email {
			m.users[i].ResetPasswordToken = &token
			m.users[i].ResetPasswordExpires = &expires
			return nil
		}
	}
	return domain.ErrUserNotFound
}

func (m *mockUserRepository) GetByResetToken(ctx context.Context, token string) (*domain.User, error) {
	for i := range m.users {
		if m.users[i].ResetPasswordToken != nil && *m.users[i].ResetPasswordToken == token {
			if m.users[i].ResetPasswordExpires != nil && time.Now().Before(*m.users[i].ResetPasswordExpires) {
				return &m.users[i], nil
			}
		}
	}
	return nil, domain.ErrUserNotFound
}

// Mock email service
type mockEmailService struct {
	lastEmailTo   string
	lastEmailName string
	shouldFail    bool
}

func (m *mockEmailService) SendForgotPasswordEmail(to, name, newPassword string) error {
	m.lastEmailTo = to
	m.lastEmailName = name
	if m.shouldFail {
		return domain.ErrEmailSendFailed
	}
	return nil
}

// Test functions
func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name        string
		request     RegisterRequest
		expectError bool
		errorType   string
	}{
		{
			name: "Valid registration",
			request: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Role:     "CLIENT",
			},
			expectError: false,
		},
		{
			name: "Empty name",
			request: RegisterRequest{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
				Role:     "CLIENT",
			},
			expectError: true,
			errorType:   "VALIDATION_ERROR",
		},
		{
			name: "Invalid email",
			request: RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
				Role:     "CLIENT",
			},
			expectError: true,
			errorType:   "VALIDATION_ERROR",
		},
		{
			name: "Short password",
			request: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "123",
				Role:     "CLIENT",
			},
			expectError: true,
			errorType:   "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockUserRepository{}
			mockEmail := &mockEmailService{}
			service := NewService(mockRepo, mockEmail)

			// Execute
			response, err := service.Register(context.Background(), tt.request)

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
				if response != nil {
					if response.Name != tt.request.Name {
						t.Errorf("Expected name %s, got %s", tt.request.Name, response.Name)
					}
					if response.Email != tt.request.Email {
						t.Errorf("Expected email %s, got %s", tt.request.Email, response.Email)
					}
				}
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Setup
	mockRepo := &mockUserRepository{}
	mockEmail := &mockEmailService{}
	service := NewService(mockRepo, mockEmail)

	// Create test user
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hashedPassword,
		Role:     "CLIENT",
	}
	mockRepo.users = append(mockRepo.users, testUser)

	tests := []struct {
		name        string
		request     LoginRequest
		expectError bool
	}{
		{
			name: "Valid login",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectError: false,
		},
		{
			name: "Invalid email",
			request: LoginRequest{
				Email:    "wrong@example.com",
				Password: "password123",
			},
			expectError: true,
		},
		{
			name: "Invalid password",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			expectError: true,
		},
		{
			name: "Empty email",
			request: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			response, err := service.Login(context.Background(), tt.request)

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if response != nil && response.AccessToken == "" {
					t.Errorf("Expected access token but got empty string")
				}
			}
		})
	}
}

func TestAuthService_ForgotPassword(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		userExists  bool
		emailFails  bool
		expectError bool
	}{
		{
			name:        "Valid forgot password",
			email:       "test@example.com",
			userExists:  true,
			emailFails:  false,
			expectError: false,
		},
		{
			name:        "User not found",
			email:       "notfound@example.com",
			userExists:  false,
			emailFails:  false,
			expectError: true,
		},
		{
			name:        "Email service fails",
			email:       "test@example.com",
			userExists:  true,
			emailFails:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mockUserRepository{}
			mockEmail := &mockEmailService{shouldFail: tt.emailFails}
			service := NewService(mockRepo, mockEmail)

			if tt.userExists {
				testUser := domain.User{
					ID:    primitive.NewObjectID(),
					Name:  "Test User",
					Email: tt.email,
					Role:  "CLIENT",
				}
				mockRepo.users = append(mockRepo.users, testUser)
			}

			// Execute
			err := service.ForgotPassword(context.Background(), ForgotPasswordRequest{Email: tt.email})

			// Assert
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				// Check if email was sent
				if mockEmail.lastEmail == nil {
					t.Errorf("Expected email to be sent")
				}
			}
		})
	}
}

// Performance test
func TestAuthService_LoginPerformance(t *testing.T) {
	// Setup
	mockRepo := &mockUserRepository{}
	mockEmail := &mockEmailService{}
	service := NewService(mockRepo, mockEmail)

	// Create test user
	hashedPassword, _ := utils.HashPassword("password123")
	testUser := domain.User{
		ID:       primitive.NewObjectID(),
		Name:     "Test User",
		Email:    "perf@example.com",
		Password: hashedPassword,
		Role:     "CLIENT",
	}
	mockRepo.users = append(mockRepo.users, testUser)

	// Performance test
	start := time.Now()

	for i := 0; i < 100; i++ {
		_, err := service.Login(context.Background(), LoginRequest{
			Email:    "perf@example.com",
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("Login failed in performance test: %v", err)
		}
	}

	duration := time.Since(start)
	avgPerRequest := duration / 100

	// Should complete 100 logins in reasonable time
	if avgPerRequest > 10*time.Millisecond {
		t.Errorf("Login performance too slow: %v per request", avgPerRequest)
	}

	t.Logf("Login performance: %v per request (100 requests in %v)", avgPerRequest, duration)
}
