package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"finsolvz-backend/internal/app/auth"
	"finsolvz-backend/internal/app/company"
	"finsolvz-backend/internal/app/user"
	"finsolvz-backend/internal/config"
	"finsolvz-backend/internal/platform/http/middleware"
	"finsolvz-backend/internal/repository"
	"finsolvz-backend/internal/utils"
)

// Test configuration
const (
	testDBName = "finsolvz_test"
	testPort   = "8788"
)

// Test server setup
type TestServer struct {
	Server    *httptest.Server
	Router    *mux.Router
	DB        *mongo.Database
	AuthToken string
	AdminUser *auth.AuthResponse
}

// Setup test server
func setupTestServer(t *testing.T) *TestServer {
	// Setup test database
	ctx := context.Background()

	// Use environment variable or default to localhost
	mongoURI := os.Getenv("TEST_MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/" + testDBName
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Skipf("Skipping integration tests: MongoDB not available (%v)", err)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("Skipping integration tests: Cannot ping MongoDB (%v)", err)
	}

	db := client.Database(testDBName)

	// Clean test database
	if err := db.Drop(ctx); err != nil {
		t.Logf("Warning: Could not drop test database: %v", err)
	}

	// Create indexes for test database
	if err := config.CreateIndexes(db); err != nil {
		t.Logf("Warning: Could not create indexes: %v", err)
	}

	// Setup repositories
	userRepo := repository.NewUserMongoRepository(db)
	companyRepo := repository.NewCompanyMongoRepository(db)

	// Setup services
	emailService := utils.NewEmailService()
	authService := auth.NewService(userRepo, emailService)
	userService := user.NewService(userRepo)
	companyService := company.NewService(companyRepo, userRepo)

	// Setup handlers
	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService, authService)
	companyHandler := company.NewHandler(companyService)

	// Setup router
	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)

	// Register routes
	authHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router, middleware.AuthMiddleware)
	companyHandler.RegisterRoutes(router, middleware.AuthMiddleware)

	// Health check endpoint
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondJSON(w, http.StatusOK, map[string]string{
			"message": "Test API",
			"status":  "healthy",
		})
	}).Methods("GET")

	// Create test server
	server := httptest.NewServer(router)

	return &TestServer{
		Server: server,
		Router: router,
		DB:     db,
	}
}

// Cleanup test server
func (ts *TestServer) Cleanup(t *testing.T) {
	ts.Server.Close()
	if err := ts.DB.Drop(context.Background()); err != nil {
		t.Logf("Warning: Could not cleanup test database: %v", err)
	}
}

// Helper function to make HTTP requests
func (ts *TestServer) makeRequest(method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, ts.Server.URL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// Test health check
func TestIntegration_HealthCheck(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Cleanup(t)

	resp, err := ts.makeRequest("GET", "/", nil, nil)
	if err != nil {
		t.Fatalf("Health check request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}
}

// Test user registration and login flow
func TestIntegration_AuthFlow(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Cleanup(t)

	// Test registration
	registerReq := map[string]interface{}{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
		"role":     "CLIENT",
	}

	resp, err := ts.makeRequest("POST", "/api/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Registration request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var registerResponse auth.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResponse); err != nil {
		t.Fatalf("Failed to decode registration response: %v", err)
	}

	if registerResponse.User.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", registerResponse.User.Name)
	}

	if registerResponse.Token == "" {
		t.Errorf("Expected access token, got empty string")
	}

	// Test login
	loginReq := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	resp, err = ts.makeRequest("POST", "/api/login", loginReq, nil)
	if err != nil {
		t.Fatalf("Login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var loginResponse auth.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	if loginResponse.Token == "" {
		t.Errorf("Expected access token, got empty string")
	}

	// Store auth token for subsequent tests
	ts.AuthToken = loginResponse.Token
	ts.AdminUser = &loginResponse
}

// Test protected endpoints
func TestIntegration_ProtectedEndpoints(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Cleanup(t)

	// First, register and login to get auth token
	registerReq := map[string]interface{}{
		"name":     "Admin User",
		"email":    "admin@example.com",
		"password": "password123",
		"role":     "SUPER_ADMIN",
	}

	resp, err := ts.makeRequest("POST", "/api/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	defer resp.Body.Close()

	var authResponse auth.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatalf("Failed to decode auth response: %v", err)
	}

	authHeaders := map[string]string{
		"Authorization": "Bearer " + authResponse.Token,
	}

	// Test accessing protected endpoint without auth
	resp, err = ts.makeRequest("GET", "/api/users", nil, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Test accessing protected endpoint with auth
	resp, err = ts.makeRequest("GET", "/api/users", nil, authHeaders)
	if err != nil {
		t.Fatalf("Authenticated request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// Test company creation flow
func TestIntegration_CompanyFlow(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Cleanup(t)

	// Setup: Register admin user
	registerReq := map[string]interface{}{
		"name":     "Company Admin",
		"email":    "admin@company.com",
		"password": "password123",
		"role":     "ADMIN",
	}

	resp, err := ts.makeRequest("POST", "/api/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	defer resp.Body.Close()

	var authResponse auth.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatalf("Failed to decode auth response: %v", err)
	}

	authHeaders := map[string]string{
		"Authorization": "Bearer " + authResponse.Token,
	}

	// Test: Create company
	companyReq := map[string]interface{}{
		"name": "Test Company Ltd",
		"user": []string{authResponse.ID},
	}

	resp, err = ts.makeRequest("POST", "/api/company", companyReq, authHeaders)
	if err != nil {
		t.Fatalf("Company creation failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var companyResponse company.CompanyResponse
	if err := json.NewDecoder(resp.Body).Decode(&companyResponse); err != nil {
		t.Fatalf("Failed to decode company response: %v", err)
	}

	if companyResponse.Name != "Test Company Ltd" {
		t.Errorf("Expected company name 'Test Company Ltd', got %s", companyResponse.Name)
	}

	// Test: Get companies
	resp, err = ts.makeRequest("GET", "/api/company", nil, authHeaders)
	if err != nil {
		t.Fatalf("Get companies failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var companies []company.CompanyResponse
	if err := json.NewDecoder(resp.Body).Decode(&companies); err != nil {
		t.Fatalf("Failed to decode companies response: %v", err)
	}

	if len(companies) != 1 {
		t.Errorf("Expected 1 company, got %d", len(companies))
	}

	// Test: Get company by ID
	companyID := companyResponse.ID
	resp, err = ts.makeRequest("GET", fmt.Sprintf("/api/company/%s", companyID), nil, authHeaders)
	if err != nil {
		t.Fatalf("Get company by ID failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// Test error handling
func TestIntegration_ErrorHandling(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Cleanup(t)

	// Test invalid JSON
	resp, err := ts.makeRequest("POST", "/api/login", "invalid json", nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	// Test missing required fields
	loginReq := map[string]interface{}{
		"email": "test@example.com",
		// Missing password
	}

	resp, err = ts.makeRequest("POST", "/api/login", loginReq, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	// Test non-existent endpoint
	resp, err = ts.makeRequest("GET", "/api/nonexistent", nil, nil)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

// Performance test
func TestIntegration_Performance(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Cleanup(t)

	// Setup: Create admin user
	registerReq := map[string]interface{}{
		"name":     "Perf Admin",
		"email":    "perf@example.com",
		"password": "password123",
		"role":     "ADMIN",
	}

	resp, err := ts.makeRequest("POST", "/api/register", registerReq, nil)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	defer resp.Body.Close()

	var authResponse auth.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatalf("Failed to decode auth response: %v", err)
	}

	authHeaders := map[string]string{
		"Authorization": "Bearer " + authResponse.Token,
	}

	// Performance test: Multiple health checks
	start := time.Now()
	for i := 0; i < 10; i++ {
		resp, err := ts.makeRequest("GET", "/", nil, nil)
		if err != nil {
			t.Fatalf("Health check %d failed: %v", i, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Health check %d failed with status %d", i, resp.StatusCode)
		}
	}
	healthCheckDuration := time.Since(start)

	// Performance test: Multiple authenticated requests
	start = time.Now()
	for i := 0; i < 10; i++ {
		resp, err := ts.makeRequest("GET", "/api/users", nil, authHeaders)
		if err != nil {
			t.Fatalf("Authenticated request %d failed: %v", i, err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Authenticated request %d failed with status %d", i, resp.StatusCode)
		}
	}
	authRequestDuration := time.Since(start)

	// Performance assertions
	avgHealthCheck := healthCheckDuration / 10
	avgAuthRequest := authRequestDuration / 10

	t.Logf("Performance results:")
	t.Logf("  Average health check: %v", avgHealthCheck)
	t.Logf("  Average auth request: %v", avgAuthRequest)

	// Health checks should be very fast
	if avgHealthCheck > 50*time.Millisecond {
		t.Errorf("Health check too slow: %v (expected < 50ms)", avgHealthCheck)
	}

	// Authenticated requests should be reasonably fast
	if avgAuthRequest > 200*time.Millisecond {
		t.Errorf("Authenticated request too slow: %v (expected < 200ms)", avgAuthRequest)
	}
}
