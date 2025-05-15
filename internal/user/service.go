package user

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	errs "github.com/phsaurav/echo_prod_blueprint/pkg/error"
	"github.com/phsaurav/echo_prod_blueprint/pkg/response"

	"golang.org/x/crypto/bcrypt"
)

// Repository is declared on the consumer side.
// Any concrete repository must implement these methods.
type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
	ActivateUser(ctx context.Context, id int64) error
}

// Service contains business logic for user operations
type Service struct {
	Repo       Repository
	JWTSecret  string
	JWTExpires time.Duration
}

// NewService creates a new user service
func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		Repo:       repo,
		JWTSecret:  jwtSecret,
		JWTExpires: 24 * time.Hour, // Default expiration of 24 hours
	}
}

// RegisterUser handles user registration
// @Summary Register a new user
// @Description Create a new user account with username, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration details"
// @Success 200 {object} User "Successfully registered user"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid input"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/user/register [post]
func (s *Service) RegisterUser(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return response.ErrorBuilder(errs.BadRequest(errors.New("all fields required"))).Send(c)
	}

	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
		IsActive: true, // Auto-activate for simplicity
	}

	if err := s.Repo.Create(c.Request().Context(), user); err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	// Don't return the password hash
	user.Password = ""

	return response.SuccessBuilder(user).Send(c)
}

// LoginUser handles user login and returns a JWT token
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} TokenResponse "Successfully authenticated with JWT token"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid input"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/user/login [post]
func (s *Service) LoginUser(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return response.ErrorBuilder(errs.BadRequest(err)).Send(c)
	}

	user, err := s.Repo.GetByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return response.ErrorBuilder(errs.NotFound(errors.New("invalid credentials"))).Send(c)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return response.ErrorBuilder(errs.BaseErr("invalid credentials")).Send(c)
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return response.ErrorBuilder(errs.InternalServerError(err)).Send(c)
	}

	return response.SuccessBuilder(map[string]string{"token": token}).Send(c)
}

// GetUser retrieves user details by ID
// @Summary Get user profile
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User "User details"
// @Failure 400 {object} response.ErrorResponse "Bad request - invalid ID format"
// @Failure 404 {object} response.ErrorResponse "Not found - user doesn't exist"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/user/{id} [get]
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

// generateJWT creates a new JWT token for the user
func (s *Service) generateJWT(user *User) (string, error) {
	// Create token with claims
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(s.JWTExpires).Unix()

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
