package user

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockRepository implements user.Repository interface for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, u *User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) UpdatePassword(ctx context.Context, id int64, password string) error {
	args := m.Called(ctx, id, password)
	return args.Error(0)
}

func (m *MockRepository) ActivateUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

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

// MockUserService implements user.UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetProfile(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockUserService) UpdatePassword(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockUserService) GetUser(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockUserService) LoginUser(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockUserService) RegisterUser(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
