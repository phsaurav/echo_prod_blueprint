package poll

import (
	"context"
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

// CreatePoll creates a new poll with options
// @Summary Create a new poll
// @Description Create a new poll with a question and multiple options
// @Tags polls
// @Accept json
// @Produce json
// @Param request body CreatePollRequest true "Poll creation request"
// @Success 200 {object} CreatePollResponse "Successfully created poll"
// @Failure 400 {object} response.FailedResponse "Bad request - invalid input"
// @Failure 401 {object} response.FailedResponse "Unauthorized - authentication required"
// @Failure 500 {object} response.FailedResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/poll [post]
func (s *Service) CreatePoll(c echo.Context) error {
	var req CreatePollRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	if req.Question == "" {
		return response.ErrorBuilder(errs.BaseErr("question is required")).Send(c)
	}
	if len(req.Options) < 2 {
		return response.ErrorBuilder(errs.BaseErr("at least two options are required")).Send(c)
	}

	// Get authenticated user ID
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	// Create poll and options
	poll := &Poll{
		Question:  req.Question,
		UserID:    userID,
		CreatedAt: time.Now(),
		Options:   make([]Option, len(req.Options)),
	}

	// Populate options
	for i, text := range req.Options {
		poll.Options[i] = Option{
			Text: text,
		}
	}

	if err := s.Repo.Create(c.Request().Context(), poll); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return response.SuccessBuilder(CreatePollResponse{Poll: *poll}).Send(c)
}

// GetPoll retrieves poll details by ID
// @Summary Get poll information
// @Description Get poll details including available options
// @Tags polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Success 200 {object} Poll "Poll details with options"
// @Failure 400 {object} response.FailedResponse "Bad request - invalid ID format"
// @Failure 404 {object} response.FailedResponse "Not found - poll doesn't exist"
// @Failure 500 {object} response.FailedResponse "Internal server error"
// @Router /api/v1/poll/{id} [get]
func (s *Service) GetPoll(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	poll, err := s.Repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	return response.SuccessBuilder(poll).Send(c)
}

// VotePoll records a user's vote for a specific poll option
// @Summary Vote on a poll
// @Description Submit a vote for a specific option in a poll
// @Tags polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Param request body VotePollRequest true "Vote details with option_id"
// @Success 200 {object} VotePollResponse "Vote successfully recorded with details"
// @Failure 400 {object} response.FailedResponse "Bad request - invalid input or poll ID"
// @Failure 401 {object} response.FailedResponse "Unauthorized - authentication required"
// @Failure 403 {object} response.FailedResponse "Forbidden - user has already voted"
// @Failure 500 {object} response.FailedResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/poll/{id}/vote [post]
func (s *Service) VotePoll(c echo.Context) error {
	idStr := c.Param("id")
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	var req VotePollRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}
	if req.OptionID == 0 {
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

	if err := s.Repo.Vote(c.Request().Context(), pollID, req.OptionID, userID); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	// Return a more informative response instead of just a status code
	resp := VotePollResponse{
		Message:   "Vote recorded successfully",
		PollID:    pollID,
		OptionID:  req.OptionID,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return response.SuccessBuilder(resp).Send(c)
}

// GetResults retrieves the current results of a poll
// @Summary Get poll results
// @Description Get the current vote counts for each option in a poll
// @Tags polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Success 200 {object} PollResultsResponse "Poll results with options and vote counts"
// @Failure 400 {object} response.FailedResponse "Bad request - invalid poll ID format"
// @Failure 404 {object} response.FailedResponse "Not found - poll doesn't exist"
// @Failure 500 {object} response.FailedResponse "Internal server error"
// @Router /api/v1/poll/{id}/results [get]
func (s *Service) GetResults(c echo.Context) error {
	idStr := c.Param("id")
	pollID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	// First get the poll details
	poll, err := s.Repo.GetByID(c.Request().Context(), pollID)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	// Then get the results
	options, err := s.Repo.GetResults(c.Request().Context(), pollID)
	if err != nil {
		return response.ErrorBuilder(err).Send(c)
	}

	// Calculate total votes
	var totalVotes int64
	for _, opt := range options {
		totalVotes += opt.Votes
	}

	// Create the response
	results := PollResultsResponse{
		PollID:     pollID,
		Question:   poll.Question,
		TotalVotes: totalVotes,
		CreatedAt:  poll.CreatedAt,
		Options:    options,
	}

	return response.SuccessBuilder(results).Send(c)
}
