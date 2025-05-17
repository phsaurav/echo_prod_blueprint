package poll

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	errs "github.com/phsaurav/echo_prod_blueprint/pkg/error"
	"github.com/phsaurav/echo_prod_blueprint/pkg/response"
)

type Repository interface {
	Create(ctx context.Context, p *Poll) error
	GetByID(ctx context.Context, id int64) (*Poll, error)
	Vote(ctx context.Context, pollID, optionID, userID int64) error
	GetResults(ctx context.Context, pollID int64) ([]Option, error)
	HasUserVoted(ctx context.Context, pollID int64, userID int64) (bool, error)
}

// Service implements the consumer-side PollService interface.
type Service struct {
	Repo Repository
}

// NewService creates a new poll service instance.
func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) CreatePoll(c echo.Context) error {
	// Parse request payload to get question/options
	var req struct {
		Question string   `json:"question"`
		Options  []string `json:"options"`
	}
	if err := c.Bind(&req); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}
	// Get authenticated user ID from context/JWT
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	poll := &Poll{
		Question:  req.Question,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := s.Repo.Create(c.Request().Context(), poll); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}
	return response.SuccessBuilder(poll).Send(c)
}

func (s *Service) GetPoll(c echo.Context) error {
	idStr := c.Param("id")
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	p, err := s.Repo.GetByID(c.Request().Context(), pollID)
	if err != nil {
		return response.ErrorBuilder(errs.NotFound(err)).Send(c)
	}

	return response.SuccessBuilder(p).Send(c)
}

func (s *Service) VotePoll(c echo.Context) error {
	idStr := c.Param("id")
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	var body struct {
		OptionID int64 `json:"option_id"`
	}
	if err := c.Bind(&body); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}
	if body.OptionID == 0 {
		return response.ErrorBuilder(errs.BaseErr("option_id is required")).Send(c)
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	alreadyVoted, err := s.Repo.HasUserVoted(c.Request().Context(), pollID, userID)
	if err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}
	if alreadyVoted {
		return response.ErrorBuilder(errs.BaseErr("already voted")).Send(c)
	}

	if err := s.Repo.Vote(c.Request().Context(), pollID, body.OptionID, userID); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Service) GetResults(c echo.Context) error {
	idStr := c.Param("id")
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	options, err := s.Repo.GetResults(c.Request().Context(), pollID)
	if err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return response.SuccessBuilder(options).Send(c)
}
