package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew tests the creation of a new database service
func TestNew(t *testing.T) {
	// Skip integration tests when running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test with valid configuration
	t.Run("ValidConfig", func(t *testing.T) {
		// Using an in-memory SQLite database for testing
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		// Replace the sql.Open call with our mocked DB
		origSqlOpen := sqlOpen
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return db, nil
		}
		defer func() { sqlOpen = origSqlOpen }()

		// Setup expectations
		mock.ExpectPing()

		// Call New with test parameters
		dbService, err := New("postgres://user:password@localhost/testdb?sslmode=disable", 10, 5, "15m")

		// Assert no error and service is created
		assert.NoError(t, err)
		assert.NotNil(t, dbService)

		// Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Test with connection error
	t.Run("ConnectionError", func(t *testing.T) {
		// First reset the singleton
		dbInstance = nil

		// Replace the sql.Open call to simulate connection error
		origSqlOpen := sqlOpen
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return nil, errors.New("connection error")
		}
		defer func() { sqlOpen = origSqlOpen }()

		// Call New with test parameters
		dbService, err := New("invalid-connection-string", 10, 5, "15m")

		// Assert error and nil service
		assert.Error(t, err, "Expected an error when connection fails")
		assert.Nil(t, dbService, "Expected nil service when connection fails")
		assert.Contains(t, err.Error(), "connection error", "Error should contain the connection error message")
	})

	// Test with invalid duration
	t.Run("InvalidDuration", func(t *testing.T) {
		// Using an in-memory SQLite database for testing
		db, _, err := sqlmock.New()
		require.NoError(t, err)

		// Replace the sql.Open call with our mocked DB
		origSqlOpen := sqlOpen
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return db, nil
		}
		defer func() { sqlOpen = origSqlOpen }()

		// Call New with invalid duration
		dbService, err := New("postgres://user:password@localhost/testdb?sslmode=disable", 10, 5, "invalid")

		// Assert error and nil service
		assert.Error(t, err)
		assert.Nil(t, dbService)
		assert.Contains(t, err.Error(), "time: invalid duration")
	})

	// Test reusing existing connection
	t.Run("ReuseConnection", func(t *testing.T) {
		// First reset the singleton
		dbInstance = nil

		// Using an in-memory SQLite database for testing
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		// Replace the sql.Open call with our mocked DB
		origSqlOpen := sqlOpen
		sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
			return db, nil
		}
		defer func() { sqlOpen = origSqlOpen }()

		// Setup expectations - ping should only be called once
		mock.ExpectPing()

		// First call to New
		service1, err := New("postgres://user:password@localhost/testdb?sslmode=disable", 10, 5, "15m")
		assert.NoError(t, err)
		assert.NotNil(t, service1)

		// Second call should reuse the existing connection
		service2, err := New("postgres://user:password@localhost/testdb?sslmode=disable", 10, 5, "15m")
		assert.NoError(t, err)
		assert.NotNil(t, service2)

		// Both services should be the same instance
		assert.Equal(t, service1, service2)

		// Verify all expectations were met (ping should only be called once)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestService_Health tests the Health method of the service
func TestService_Health(t *testing.T) {
	t.Run("HealthyDatabase", func(t *testing.T) {
		// Create a mock database
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)

		// Create a service with mocked DB
		svc := &service{db: db}

		// Setup expectations - ping succeeds
		mock.ExpectPing()

		// Call Health method
		stats := svc.Health()

		// Assert health stats
		assert.Equal(t, "up", stats["status"])
		assert.Equal(t, "It's healthy", stats["message"])
		assert.Contains(t, stats, "open_connections")
		assert.Contains(t, stats, "in_use")
		assert.Contains(t, stats, "idle")

		// Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UnhealthyDatabase", func(t *testing.T) {
		// Create a mock database WITH ping monitoring enabled
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true)) // <-- Add this option
		require.NoError(t, err)

		// Create a service with mocked DB
		svc := &service{db: db}

		// Setup expectations - ping fails
		mock.ExpectPing().WillReturnError(errors.New("connection lost"))

		stats := svc.Health()

		// Assert health stats for down database
		assert.Equal(t, "down", stats["status"])
		assert.Contains(t, stats["error"], "db down: connection lost")

		// Verify all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// TestService_Close tests the Close method of the service
func TestService_Close(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Create a service with mocked DB
	svc := &service{db: db}

	// Setup expectations - close succeeds
	mock.ExpectClose()

	// Call Close method
	err = svc.Close()

	// Assert no error
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestService_DB tests the DB method of the service
func TestService_DB(t *testing.T) {
	// Create a mock database
	db, _, err := sqlmock.New()
	require.NoError(t, err)

	// Create a service with mocked DB
	svc := &service{db: db}

	// Call DB method
	returnedDB := svc.DB()

	// Assert returned DB is same as input DB
	assert.Equal(t, db, returnedDB)
}
