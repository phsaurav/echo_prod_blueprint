package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	errs "github.com/phsaurav/go_echo_base/pkg/error"
	"github.com/phsaurav/go_echo_base/pkg/response"
)

// Repository is declared on the consumer side.
// Any concrete repository must implement these methods.
type Repository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, u *User) error
	Follow(ctx context.Context, followedID, followerID int64) error
	Unfollow(ctx context.Context, followerID, unfollowedID int64) error
	Activate(ctx context.Context, token string) error
}

// Service implements the consumer-side UserService interface.
type Service struct {
	Repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

// CreateUser handles POST /users.
func (s *Service) CreateUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	if err := s.Repo.Create(c.Request().Context(), u); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return response.SuccessBuilder(u).Send(c)
}

// GetUser handles GET /users/:id.
func (s *Service) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	u, err := s.Repo.GetByID(c.Request().Context(), userID)
	if err != nil {
		return response.ErrorBuilder(errs.NotFound(err)).Send(c)
	}

	return response.SuccessBuilder(u).Send(c)
}

// FollowUser handles PUT /users/:id/follow.
func (s *Service) FollowUser(c echo.Context) error {
	// In a real app, obtain followerID (e.g. from auth context)
	followerID := int64(1)
	idStr := c.Param("id")
	followedID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	if err := s.Repo.Follow(c.Request().Context(), followedID, followerID); err != nil {
		return response.ErrorBuilder(errs.Conflict(err)).Send(c)
	}

	return c.NoContent(http.StatusNoContent)
}

// UnfollowUser handles PUT /users/:id/unfollow.
func (s *Service) UnfollowUser(c echo.Context) error {
	// In a real app, obtain followerID (e.g. from auth context)
	followerID := int64(1)
	idStr := c.Param("id")
	unfollowedID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	if err := s.Repo.Unfollow(c.Request().Context(), followerID, unfollowedID); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return c.NoContent(http.StatusNoContent)
}

// ActivateUser handles PUT /users/activate/:token.
func (s *Service) ActivateUser(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return response.ErrorBuilder(errs.BadRequest(errs.BaseErr("missing token"))).Send(c)
	}

	if err := s.Repo.Activate(c.Request().Context(), token); err != nil {
		return response.ErrorBuilder(errs.NotFound(err)).Send(c)
	}

	return c.NoContent(http.StatusNoContent)
}
