package post

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/labstack/echo/v4"
	errs "github.com/phsaurav/go_echo_base/pkg/error"
	"github.com/phsaurav/go_echo_base/pkg/response"
)

// Repository defines expected behavior for the post repository.
type Repository interface {
	Create(ctx context.Context, p *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
}

// Service implements business logic for posts.
type Service struct {
	Repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

// CreatePostPayload is the request structure for creating a post.
type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
	UserID  int64    `json:"user_id"`
}

// CreatePost handles POST /posts.
func (s *Service) CreatePost(c echo.Context) error {
	payload := new(CreatePostPayload)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	post := &Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  payload.UserID,
	}

	if err := s.Repo.Create(c.Request().Context(), post); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return response.SuccessBuilder(post).Send(c)
}

// GetPost handles GET /posts/:id.
func (s *Service) GetPost(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	post, err := s.Repo.GetByID(c.Request().Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return response.ErrorBuilder(errs.NotFound(err)).Send(c)
		}
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return response.SuccessBuilder(post).Send(c)
}
