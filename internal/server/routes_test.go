package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelloWorldHandler tests the hello world handler
func TestHelloWorldHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	s := &Server{}

	// Assertions
	if err := s.HelloWorldHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}

	// For debugging
	t.Logf("Response body: %s", resp.Body.String())

	// Use map[string]interface{} to handle mixed type values
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Errorf("handler() error decoding response body: %v", err)
		return
	}

	// Check for standard response structure
	if _, ok := response["code"]; !ok {
		t.Errorf("response missing 'code' field: %v", response)
		return
	}
	if _, ok := response["message"]; !ok {
		t.Errorf("response missing 'message' field: %v", response)
		return
	}

	// Check for data field
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Errorf("response has wrong 'data' type or missing 'data': %v", response)
		return
	}

	// Check expected message in data
	message, ok := data["message"].(string)
	if !ok {
		t.Errorf("data missing string 'message' field: %v", data)
		return
	}
	if message != "Hello World" {
		t.Errorf("wrong message. expected='Hello World', got=%v", message)
		return
	}
}

// TestHealthHandler tests the health check handler
func TestHealthHandler(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create mock DB service
	mockDBService := new(MockDBService)
	healthData := map[string]string{"database": "healthy", "status": "ok"}
	mockDBService.On("Health").Return(healthData).Once()

	// Create server with mock store
	s := &Server{
		e:     e,
		store: NewStore(mockDBService),
	}

	// Test the handler
	err := s.healthHandler(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check response body
	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	// For debugging, print the actual response
	t.Logf("Response body: %s", rec.Body.String())

	// Assert on the response structure
	assert.Contains(t, response, "data")

	// Check the data field contains our health information
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Contains(t, data, "database")
		assert.Contains(t, data, "status")
	} else {
		t.Fatalf("'data' field is not a map: %v", response["data"])
	}

	// Verify expectations
	mockDBService.AssertExpectations(t)
}

// TestRegisterRoutes tests that routes are properly registered
func TestRegisterRoutes(t *testing.T) {
	// Setup
	e := echo.New()

	// Create mock DB service
	mockDBService := new(MockDBService)
	mockDBService.On("Health").Return(map[string]string{"status": "ok"}).Maybe()
	mockDBService.On("DB").Return(nil).Maybe()

	// Create server with the mock
	s := &Server{
		e:     e,
		store: NewStore(mockDBService),
	}

	// Register routes
	handler := s.RegisterRoutes()

	// Verify handler is created
	assert.NotNil(t, handler)

	// Test the root route
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Hello World")

	// Test the health route
	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	// Verify response
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify mock expectations
	mockDBService.AssertExpectations(t)
}

// TestRegisterAPIRoutes tests API route registration
func TestRegisterAPIRoutes(t *testing.T) {
	// Setup
	e := echo.New()

	// Create mock DB service
	mockDBService := new(MockDBService)
	mockDBService.On("DB").Return(nil).Maybe()

	// Create server
	s := &Server{
		e:     e,
		store: NewStore(mockDBService),
	}

	// Test that registering routes doesn't panic
	assert.NotPanics(t, func() {
		s.routes(e.Group("/api"), "v1")
	})

	// Note: Full integration testing of API routes would require a more
	// complex setup with mocked user and poll services
}
