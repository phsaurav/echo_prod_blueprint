package poll

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockRepository implements poll.Repository interface for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, p *Poll) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*Poll, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Poll), args.Error(1)
}

func (m *MockRepository) Vote(ctx context.Context, pollID, optionID, userID int64) error {
	args := m.Called(ctx, pollID, optionID, userID)
	return args.Error(0)
}

func (m *MockRepository) GetResults(ctx context.Context, pollID int64) ([]Option, error) {
	args := m.Called(ctx, pollID)
	return args.Get(0).([]Option), args.Error(1)
}

func (m *MockRepository) HasUserVoted(ctx context.Context, pollID int64, userID int64) (bool, error) {
	args := m.Called(ctx, pollID, userID)
	return args.Bool(0), args.Error(1)
}

// MockDBService implements database.Service interface for testing
type MockDBService struct {
	mock.Mock
	db *sql.DB
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

// MockPollService implements poll.PollService for testing
type MockPollService struct {
	mock.Mock
}

func (m *MockPollService) CreatePoll(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockPollService) GetPoll(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockPollService) VotePoll(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockPollService) GetResults(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
