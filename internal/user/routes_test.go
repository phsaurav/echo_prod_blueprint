package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/phsaurav/echo_prod_blueprint/config"
	"github.com/phsaurav/echo_prod_blueprint/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegisterRoutes tests that routes are registered correctly
// TestRegisterRoutes tests that routes are registered correctly
func TestRegisterRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api/v1/user")

	mockService := new(MockUserService)
	authMiddleware := testutils.CreateAuthMiddleware()

	// Register routes
	RegisterRoutes(g, mockService, authMiddleware)

	// Test POST /api/v1/user/register
	mockService.On("RegisterUser", mock.Anything).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Test POST /api/v1/user/login
	mockService.On("LoginUser", mock.Anything).Return(nil)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/login", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Important: Set up GetUser expectation BEFORE GetProfile
	// because of how Echo matches routes
	mockService.On("GetUser", mock.Anything).Return(nil).Maybe()

	// Test GET /api/v1/user/profile - Add auth token
	mockService.On("GetProfile", mock.Anything).Return(nil).Maybe()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/user/profile", nil)
	// Add auth header for protected routes
	req.Header.Set("Authorization", "Bearer test-token")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Test GET /api/v1/user/:id - Add auth token
	req = httptest.NewRequest(http.MethodGet, "/api/v1/user/123", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Test PUT /api/v1/user/password - Add auth token
	mockService.On("UpdatePassword", mock.Anything).Return(nil).Maybe()
	req = httptest.NewRequest(http.MethodPut, "/api/v1/user/password", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify all expectations were met
	mockService.AssertExpectations(t)
}

// TestRegister tests the Register function
func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api/v1/user")

	// Create mock database service
	mockDB := new(MockDBService)
	mockDB.On("DB").Return(nil)

	// Mock auth middleware
	authMiddleware := testutils.CreateAuthMiddleware()

	cfg := config.Config{
		TokenConfig: config.TokenConfig{
			Secret: "test-secret-key",
			Exp:    24 * time.Hour,
			Iss:    "TestIssuer",
		},
		Env: "test",
		// Add any other required fields
	}

	assert.NotPanics(t, func() {
		Register(g, mockDB, cfg, authMiddleware)
	})

	// Verify mock was called
	mockDB.AssertExpectations(t)
}
