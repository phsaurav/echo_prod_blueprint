package poll

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/phsaurav/echo_prod_blueprint/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegisterRoutes tests that routes are registered correctly
func TestRegisterRoutes(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api/v1/poll")

	mockService := new(MockPollService)

	// Mock auth middleware
	authMiddleware := testutils.CreateAuthMiddleware()

	// Register routes
	RegisterRoutes(g, mockService, authMiddleware)

	// Test POST /api/v1/poll
	mockService.On("CreatePoll", mock.Anything).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/poll", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Test GET /api/v1/poll/:id
	mockService.On("GetPoll", mock.Anything).Return(nil)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/poll/1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Test POST /api/v1/poll/:id/vote
	mockService.On("VotePoll", mock.Anything).Return(nil)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/poll/1/vote", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Test GET /api/v1/poll/:id/results
	mockService.On("GetResults", mock.Anything).Return(nil)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/poll/1/results", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Verify all expected methods were called
	mockService.AssertExpectations(t)
}

// TestRegister tests the Register function
func TestRegister(t *testing.T) {
	// Setup
	e := echo.New()
	g := e.Group("/api/v1/poll")

	// Create mock database service
	mockDB := new(MockDBService)
	mockDB.On("DB").Return(nil)

	// Mock auth middleware
	authMiddleware := testutils.CreateAuthMiddleware()

	assert.NotPanics(t, func() {
		Register(g, mockDB, authMiddleware)
	})

	// Verify mock was called
	mockDB.AssertExpectations(t)
}
