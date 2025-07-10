package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"finsolvz-backend/internal/app/auth"
	"finsolvz-backend/internal/app/company"
)

// E2E tests run against a live server
// Set FINSOLVZ_E2E_URL environment variable to test against live server
// Example: FINSOLVZ_E2E_URL=https://your-service.a.run.app go test ./tests -run TestE2E

const (
	defaultE2ETimeout = 30 * time.Second
)

// E2E test configuration
type E2EConfig struct {
	BaseURL string
	Timeout time.Duration
	Client  *http.Client
}

// Setup E2E test configuration
func setupE2E(t *testing.T) *E2EConfig {
	baseURL := os.Getenv("FINSOLVZ_E2E_URL")
	if baseURL == "" {
		t.Skip("Skipping E2E tests: FINSOLVZ_E2E_URL not set")
	}

	return &E2EConfig{
		BaseURL: baseURL,
		Timeout: defaultE2ETimeout,
		Client: &http.Client{
			Timeout: defaultE2ETimeout,
		},
	}
}

// Helper function to make E2E HTTP requests
func (cfg *E2EConfig) makeRequest(method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
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

	req, err := http.NewRequest(method, cfg.BaseURL+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return cfg.Client.Do(req)
}

// Test production health check
func TestE2E_HealthCheck(t *testing.T) {
	cfg := setupE2E(t)

	resp, err := cfg.makeRequest("GET", "/", nil, nil)
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
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

	t.Logf("✅ Health check passed: %v", response["message"])
}

// Test production performance
func TestE2E_Performance(t *testing.T) {
	cfg := setupE2E(t)

	// Define performance thresholds based on region
	thresholds := map[string]time.Duration{
		"health":    200 * time.Millisecond, // Health check should be very fast
		"auth":      500 * time.Millisecond, // Authentication should be reasonable
		"protected": 1 * time.Second,        // Protected endpoints can be slower
	}

	tests := []struct {
		name      string
		method    string
		path      string
		body      interface{}
		headers   map[string]string
		threshold time.Duration
		setup     func() map[string]string // Returns headers for subsequent requests
	}{
		{
			name:      "Health Check Performance",
			method:    "GET",
			path:      "/",
			threshold: thresholds["health"],
		},
		{
			name:   "Login Performance",
			method: "POST",
			path:   "/api/login",
			body: map[string]interface{}{
				"email":    "admin@finsolvz.com", // Assuming test admin exists
				"password": "admin123",
			},
			threshold: thresholds["auth"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup if needed
			var headers map[string]string
			if tt.setup != nil {
				headers = tt.setup()
			}
			if tt.headers != nil {
				if headers == nil {
					headers = make(map[string]string)
				}
				for k, v := range tt.headers {
					headers[k] = v
				}
			}

			// Make request and measure time
			start := time.Now()
			resp, err := cfg.makeRequest(tt.method, tt.path, tt.body, headers)
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Check performance
			if duration > tt.threshold {
				t.Errorf("Request too slow: %v (threshold: %v)", duration, tt.threshold)
			} else {
				t.Logf("✅ %s: %v (threshold: %v)", tt.name, duration, tt.threshold)
			}

			// Basic status check
			if resp.StatusCode >= 500 {
				t.Errorf("Server error: %d", resp.StatusCode)
			}
		})
	}
}

// Test complete user journey
func TestE2E_UserJourney(t *testing.T) {
	cfg := setupE2E(t)

	// Generate unique email for this test
	testEmail := fmt.Sprintf("e2e-test-%d@example.com", time.Now().Unix())

	t.Logf("🧪 Starting E2E user journey test with email: %s", testEmail)

	// Step 1: Try to login with non-existent user (should fail)
	t.Run("Login with non-existent user", func(t *testing.T) {
		loginReq := map[string]interface{}{
			"email":    testEmail,
			"password": "password123",
		}

		resp, err := cfg.makeRequest("POST", "/api/login", loginReq, nil)
		if err != nil {
			t.Fatalf("Login request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 401 or 404 for non-existent user, got %d", resp.StatusCode)
		}

		t.Logf("✅ Non-existent user login correctly rejected")
	})

	// Step 2: Register new user
	var authToken string
	var userID string

	t.Run("Register new user", func(t *testing.T) {
		registerReq := map[string]interface{}{
			"name":     "E2E Test User",
			"email":    testEmail,
			"password": "password123",
			"role":     "CLIENT",
		}

		resp, err := cfg.makeRequest("POST", "/api/register", registerReq, nil)
		if err != nil {
			t.Fatalf("Registration failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			// If registration fails (e.g., user exists), try login instead
			if resp.StatusCode == http.StatusConflict {
				t.Logf("User already exists, attempting login...")
				
				loginReq := map[string]interface{}{
					"email":    testEmail,
					"password": "password123",
				}

				loginResp, err := cfg.makeRequest("POST", "/api/login", loginReq, nil)
				if err != nil {
					t.Fatalf("Login after conflict failed: %v", err)
				}
				defer loginResp.Body.Close()

				if loginResp.StatusCode != http.StatusOK {
					t.Fatalf("Login after conflict failed with status %d", loginResp.StatusCode)
				}

				var loginResponse auth.AuthResponse
				if err := json.NewDecoder(loginResp.Body).Decode(&loginResponse); err != nil {
					t.Fatalf("Failed to decode login response: %v", err)
				}

				authToken = loginResponse.AccessToken
				userID = loginResponse.ID
			} else {
				t.Fatalf("Registration failed with status %d", resp.StatusCode)
			}
		} else {
			var registerResponse auth.AuthResponse
			if err := json.NewDecoder(resp.Body).Decode(&registerResponse); err != nil {
				t.Fatalf("Failed to decode registration response: %v", err)
			}

			authToken = registerResponse.AccessToken
			userID = registerResponse.ID

			if authToken == "" {
				t.Fatalf("No access token received")
			}

			t.Logf("✅ User registered successfully")
		}
	})

	// Step 3: Access protected endpoint
	t.Run("Access protected endpoint", func(t *testing.T) {
		headers := map[string]string{
			"Authorization": "Bearer " + authToken,
		}

		resp, err := cfg.makeRequest("GET", "/api/loginUser", nil, headers)
		if err != nil {
			t.Fatalf("Protected endpoint request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for protected endpoint, got %d", resp.StatusCode)
		}

		var userResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
			t.Fatalf("Failed to decode user response: %v", err)
		}

		if userResponse["email"] != testEmail {
			t.Errorf("Expected email %s, got %v", testEmail, userResponse["email"])
		}

		t.Logf("✅ Protected endpoint access successful")
	})

	// Step 4: Test performance of authenticated requests
	t.Run("Authenticated request performance", func(t *testing.T) {
		headers := map[string]string{
			"Authorization": "Bearer " + authToken,
		}

		// Test multiple requests to check caching
		var durations []time.Duration
		for i := 0; i < 3; i++ {
			start := time.Now()
			resp, err := cfg.makeRequest("GET", "/api/loginUser", nil, headers)
			duration := time.Since(start)
			durations = append(durations, duration)

			if err != nil {
				t.Fatalf("Performance test request %d failed: %v", i+1, err)
			}
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Performance test request %d failed with status %d", i+1, resp.StatusCode)
			}
		}

		// Calculate average
		var total time.Duration
		for _, d := range durations {
			total += d
		}
		average := total / time.Duration(len(durations))

		t.Logf("✅ Authenticated request performance:")
		t.Logf("   Request 1: %v", durations[0])
		t.Logf("   Request 2: %v", durations[1]) 
		t.Logf("   Request 3: %v", durations[2])
		t.Logf("   Average: %v", average)

		// Check if requests are getting faster (caching effect)
		if len(durations) >= 2 && durations[1] < durations[0] {
			t.Logf("✅ Caching appears to be working (request 2 faster than request 1)")
		}

		// Performance threshold for authenticated requests
		if average > 2*time.Second {
			t.Errorf("Authenticated requests too slow: %v average", average)
		}
	})

	t.Logf("🎉 E2E user journey completed successfully")
}

// Test API documentation endpoints
func TestE2E_Documentation(t *testing.T) {
	cfg := setupE2E(t)

	tests := []struct {
		name         string
		path         string
		expectedType string
	}{
		{
			name:         "Swagger Documentation",
			path:         "/docs",
			expectedType: "text/html",
		},
		{
			name:         "OpenAPI Specification",
			path:         "/api/openapi.yaml",
			expectedType: "application/x-yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := cfg.makeRequest("GET", tt.path, nil, nil)
			if err != nil {
				t.Fatalf("Documentation request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}

			contentType := resp.Header.Get("Content-Type")
			if contentType == "" {
				t.Logf("Warning: No Content-Type header for %s", tt.name)
			} else {
				t.Logf("✅ %s available with Content-Type: %s", tt.name, contentType)
			}
		})
	}
}

// Test rate limiting (if implemented)
func TestE2E_RateLimit(t *testing.T) {
	cfg := setupE2E(t)

	// Make many requests quickly to test rate limiting
	const numRequests = 20
	const requestWindow = 1 * time.Second

	t.Logf("🔄 Testing rate limiting with %d requests in %v", numRequests, requestWindow)

	start := time.Now()
	var rateLimitHit bool

	for i := 0; i < numRequests; i++ {
		resp, err := cfg.makeRequest("GET", "/", nil, nil)
		if err != nil {
			t.Logf("Request %d failed: %v", i+1, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			rateLimitHit = true
			t.Logf("✅ Rate limit triggered at request %d", i+1)
			
			// Check rate limit headers
			if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
				t.Logf("   Rate limit: %s", limit)
			}
			if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
				t.Logf("   Remaining: %s", remaining)
			}
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				t.Logf("   Retry after: %s seconds", retryAfter)
			}
			break
		}

		// Small delay to spread requests
		time.Sleep(requestWindow / numRequests)
	}

	duration := time.Since(start)
	t.Logf("Rate limit test completed in %v", duration)

	if !rateLimitHit {
		t.Logf("ℹ️ Rate limiting not triggered (may not be implemented or limits are high)")
	}
}

// Benchmark production performance
func BenchmarkE2E_HealthCheck(b *testing.B) {
	baseURL := os.Getenv("FINSOLVZ_E2E_URL")
	if baseURL == "" {
		b.Skip("Skipping benchmark: FINSOLVZ_E2E_URL not set")
	}

	cfg := &E2EConfig{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := cfg.makeRequest("GET", "/", nil, nil)
		if err != nil {
			b.Fatalf("Health check failed: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Health check failed with status %d", resp.StatusCode)
		}
	}
}