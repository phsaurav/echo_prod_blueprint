package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	// Create mock DB
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create mock service
	mockDBService := new(MockDBService)
	mockDBService.On("DB").Return(db)

	// Create repository
	repo := NewRepo(mockDBService)

	// Assert repository is created with the expected DB
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.DB)

	// Verify mock expectations
	mockDBService.AssertExpectations(t)
}

func TestRepo_Create(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Test data
	now := time.Now().Truncate(time.Second)
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	// Setup expectations
	userRows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, now)
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.Email, user.Password).
		WillReturnRows(userRows)

	// Call function under test
	err = repo.Create(context.Background(), user)

	// Assert no error
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Verify user was updated with ID and created time
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, now, user.CreatedAt)
}

func TestRepo_Create_Failure(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Test data
	user := &User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	// Setup expectations - insert fails
	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Username, user.Email, user.Password).
		WillReturnError(sql.ErrConnDone)

	// Call function under test
	err = repo.Create(context.Background(), user)

	// Assert error occurred
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sql: connection is already closed")

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_GetByID(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Current time for consistent testing
	now := time.Now().Truncate(time.Second)

	// Setup expectations
	userRows := sqlmock.NewRows([]string{"id", "username", "email", "created_at", "is_active"}).
		AddRow(1, "testuser", "test@example.com", now, true)
	mock.ExpectQuery("SELECT id, username, email, created_at, is_active FROM users").
		WithArgs(1).
		WillReturnRows(userRows)

	// Call function under test
	user, err := repo.GetByID(context.Background(), 1)

	// Assert no error
	assert.NoError(t, err)

	// Verify user data
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, now, user.CreatedAt)
	assert.True(t, user.IsActive)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_GetByID_NotFound(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Setup expectations - user not found
	mock.ExpectQuery("SELECT id, username, email, created_at, is_active FROM users").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	// Call function under test
	user, err := repo.GetByID(context.Background(), 999)

	// Assert not found error
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "sql: no rows in result set")

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_GetByEmail(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Current time for consistent testing
	now := time.Now().Truncate(time.Second)

	// Setup expectations
	userRows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "is_active"}).
		AddRow(1, "testuser", "test@example.com", "hashedpassword", now, true)
	mock.ExpectQuery("SELECT id, username, email, password, created_at, is_active FROM users").
		WithArgs("test@example.com").
		WillReturnRows(userRows)

	// Call function under test
	user, err := repo.GetByEmail(context.Background(), "test@example.com")

	// Assert no error
	assert.NoError(t, err)

	// Verify user data
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.Equal(t, now, user.CreatedAt)
	assert.True(t, user.IsActive)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_GetByEmail_NotFound(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Setup expectations - user not found
	mock.ExpectQuery("SELECT id, username, email, password, created_at, is_active FROM users").
		WithArgs("nonexistent@example.com").
		WillReturnError(sql.ErrNoRows)

	// Call function under test
	user, err := repo.GetByEmail(context.Background(), "nonexistent@example.com")

	// Assert not found error
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "sql: no rows in result set")

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_UpdatePassword(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Setup expectations
	mock.ExpectExec("UPDATE users SET password").
		WithArgs("newhashpassword", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call function under test
	err = repo.UpdatePassword(context.Background(), 1, "newhashpassword")

	// Assert no error
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_ActivateUser(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Setup expectations
	mock.ExpectExec("UPDATE users SET is_active").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call function under test
	err = repo.ActivateUser(context.Background(), 1)

	// Assert no error
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
