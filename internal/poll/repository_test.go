package poll

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
	mockDBService := &MockDBService{}
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
	poll := &Poll{
		Question: "What is your favorite color?",
		UserID:   1,
		Options: []Option{
			{Text: "Red"},
			{Text: "Blue"},
		},
	}

	// Setup expectations
	// 1. Transaction begins
	mock.ExpectBegin()

	// 2. Poll insertion
	pollRows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, time.Now())
	mock.ExpectQuery("INSERT INTO polls").
		WithArgs(poll.Question, poll.UserID).
		WillReturnRows(pollRows)

	// 3. Options insertion
	optionRows1 := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO poll_options").
		WithArgs(1, "Red").
		WillReturnRows(optionRows1)

	optionRows2 := sqlmock.NewRows([]string{"id"}).AddRow(2)
	mock.ExpectQuery("INSERT INTO poll_options").
		WithArgs(1, "Blue").
		WillReturnRows(optionRows2)

	// 4. Transaction commits
	mock.ExpectCommit()

	// Call function under test
	err = repo.Create(context.Background(), poll)

	// Assert no error
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())

	// Verify poll was updated with IDs
	assert.Equal(t, int64(1), poll.ID)
	assert.Equal(t, int64(1), poll.Options[0].ID)
	assert.Equal(t, int64(2), poll.Options[1].ID)
	assert.Equal(t, int64(1), poll.Options[0].PollID)
	assert.Equal(t, int64(1), poll.Options[1].PollID)
}

func TestRepo_Create_Failure(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Test data
	poll := &Poll{
		Question: "What is your favorite color?",
		UserID:   1,
		Options: []Option{
			{Text: "Red"},
			{Text: "Blue"},
		},
	}

	// Setup expectations - transaction begins but fails
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO polls").
		WithArgs(poll.Question, poll.UserID).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	// Call function under test
	err = repo.Create(context.Background(), poll)

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
	// 1. Poll query
	pollRows := sqlmock.NewRows([]string{"id", "question", "created_at"}).
		AddRow(1, "What is your favorite color?", now)
	mock.ExpectQuery("SELECT id, question, created_at FROM polls").
		WithArgs(1).
		WillReturnRows(pollRows)

	// 2. Options query
	optionRows := sqlmock.NewRows([]string{"id", "poll_id", "text"}).
		AddRow(1, 1, "Red").
		AddRow(2, 1, "Blue")
	mock.ExpectQuery("SELECT id, poll_id, text FROM poll_options").
		WithArgs(1).
		WillReturnRows(optionRows)

	// Call function under test
	poll, err := repo.GetByID(context.Background(), 1)

	// Assert no error
	assert.NoError(t, err)

	// Verify poll data
	assert.Equal(t, int64(1), poll.ID)
	assert.Equal(t, "What is your favorite color?", poll.Question)
	assert.Equal(t, now, poll.CreatedAt)
	assert.Len(t, poll.Options, 2)
	assert.Equal(t, "Red", poll.Options[0].Text)
	assert.Equal(t, "Blue", poll.Options[1].Text)

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

	// Setup expectations - poll not found
	mock.ExpectQuery("SELECT id, question, created_at FROM polls").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	// Call function under test
	poll, err := repo.GetByID(context.Background(), 999)

	// Assert not found error
	assert.Error(t, err)
	assert.Nil(t, poll)
	assert.Contains(t, err.Error(), "sql: no rows in result set")

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_Vote(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Setup expectations
	mock.ExpectExec("INSERT INTO poll_votes").
		WithArgs(1, 2, 3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call function under test
	err = repo.Vote(context.Background(), 1, 2, 3)

	// Assert no error
	assert.NoError(t, err)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_GetResults(t *testing.T) {
	// Create mock DB
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Create repository
	repo := &Repo{DB: db}

	// Setup expectations
	rows := sqlmock.NewRows([]string{"id", "poll_id", "text", "votes"}).
		AddRow(1, 1, "Red", 3).
		AddRow(2, 1, "Blue", 5)
	mock.ExpectQuery("SELECT o.id, o.poll_id, o.text, COUNT").
		WithArgs(1).
		WillReturnRows(rows)

	// Call function under test
	options, err := repo.GetResults(context.Background(), 1)

	// Assert no error
	assert.NoError(t, err)
	assert.Len(t, options, 2)
	assert.Equal(t, int64(3), options[0].Votes)
	assert.Equal(t, int64(5), options[1].Votes)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepo_HasUserVoted(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		setup    func(mock sqlmock.Sqlmock)
		pollID   int64
		userID   int64
		expected bool
		hasError bool
	}{
		{
			name: "User has voted",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
				mock.ExpectQuery("SELECT 1 FROM poll_votes").
					WithArgs(1, 2).
					WillReturnRows(rows)
			},
			pollID:   1,
			userID:   2,
			expected: true,
			hasError: false,
		},
		{
			name: "User has not voted",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1 FROM poll_votes").
					WithArgs(1, 3).
					WillReturnError(sql.ErrNoRows)
			},
			pollID:   1,
			userID:   3,
			expected: false,
			hasError: false,
		},
		{
			name: "Database error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT 1 FROM poll_votes").
					WithArgs(1, 4).
					WillReturnError(sql.ErrConnDone)
			},
			pollID:   1,
			userID:   4,
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock DB
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Create repository
			repo := &Repo{DB: db}

			// Setup expectations
			tt.setup(mock)

			// Call function under test
			hasVoted, err := repo.HasUserVoted(context.Background(), tt.pollID, tt.userID)

			// Assert results
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, hasVoted)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
