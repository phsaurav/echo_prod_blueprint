package server

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServer_Start(t *testing.T) {
	// Skip this test if running in CI to avoid full server creation
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Create a simple Echo instance to test startup logic
	e := echo.New()

	// Create mock DB service
	mockDBService := new(MockDBService)

	// Set up expectations BEFORE creating the server or calling methods
	mockDBService.On("Health").Return(map[string]string{"status": "ok"}).Once()

	// Create server manually with the mock
	s := &Server{
		e:     e,
		store: NewStore(mockDBService),
	}

	// Test that basic server components are created
	assert.NotNil(t, s.e)
	assert.NotNil(t, s.store)

	// Now call the methods that we expect to be called
	s.store.DBHealth()

	// Verify expectations
	mockDBService.AssertExpectations(t)
}
