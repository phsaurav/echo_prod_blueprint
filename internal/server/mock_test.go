package server

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockDBService implements database.Service for testing
type MockDBService struct {
	mock.Mock
}

func (m *MockDBService) Health() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

func (m *MockDBService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDBService) DB() *sql.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sql.DB)
}

// MockStore implements Store for testing
type MockStore struct {
	mock.Mock
}

func (m *MockStore) DBHealth() map[string]string {
	args := m.Called()
	return args.Get(0).(map[string]string)
}

// MockMiddleware implements middleware functions for testing
type MockMiddleware struct {
	mock.Mock
}

func (m *MockMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	args := m.Called(next)
	return args.Get(0).(echo.HandlerFunc)
}
