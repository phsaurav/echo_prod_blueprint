package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/phsaurav/echo_prod_blueprint/internal/database"
	errs "github.com/phsaurav/echo_prod_blueprint/pkg/error"
)

// Repo is a concrete implementation of the user repository.
// It must satisfy the Repository interface declared on the consumer side.
type Repo struct {
	DB *sql.DB
}

// NewRepo creates a new repository instance.
func NewRepo(db database.Service) *Repo {
	return &Repo{DB: db.DB()}
}

// Ensure Repo implements the consumer-side Repository interface.
var _ Repository = (*Repo)(nil)

// Create adds a new user to the database
func (r *Repo) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, is_active)
		VALUES ($1, $2, $3, NOW(), false)
		RETURNING id, created_at
	`
	err := r.DB.QueryRowContext(ctx, query, u.Username, u.Email, u.Password).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return errs.InternalServerError(err)
	}
	return nil
}

// GetByID retrieves a user by their ID
func (r *Repo) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at, is_active 
		FROM users 
		WHERE id = $1
	`
	u := new(User)
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound(err)
		}
		return nil, errs.InternalServerError(err)
	}
	return u, nil
}

// GetByEmail retrieves a user by their email address
func (r *Repo) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at, is_active 
		FROM users 
		WHERE email = $1
	`
	u := new(User)
	err := r.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt, &u.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound(err)
		}
		return nil, errs.InternalServerError(err)
	}
	return u, nil
}

// UpdatePassword updates a user's password
func (r *Repo) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, password, id)
	if err != nil {
		return errs.InternalServerError(err)
	}
	return nil
}

// ActivateUser sets a user's is_active flag to true
func (r *Repo) ActivateUser(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_active = true WHERE id = $2`
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return errs.InternalServerError(err)
	}
	return nil
}

