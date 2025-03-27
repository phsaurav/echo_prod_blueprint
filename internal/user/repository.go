package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/phsaurav/go_echo_base/internal/database"
	errs "github.com/phsaurav/go_echo_base/pkg/error"
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

func (r *Repo) Create(ctx context.Context, u *User) error {
	query := `
		INSERT INTO users (username, email, created_at, is_active)
		VALUES ($1, $2, NOW(), false)
		RETURNING id
	`
	err := r.DB.QueryRowContext(ctx, query, u.Username, u.Email).Scan(&u.ID)
	if err != nil {
		return errs.InternalServerError(err)
	}
	return nil
}

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
		return nil, err
	}
	return u, nil
}

func (r *Repo) Follow(ctx context.Context, followedID, followerID int64) error {
	query := `INSERT INTO followers (followed_id, follower_id) VALUES ($1, $2)`
	_, err := r.DB.ExecContext(ctx, query, followedID, followerID)
	return err
}

func (r *Repo) Unfollow(ctx context.Context, followerID, unfollowedID int64) error {
	query := `DELETE FROM followers WHERE follower_id = $1 AND followed_id = $2`
	_, err := r.DB.ExecContext(ctx, query, followerID, unfollowedID)
	return err
}

func (r *Repo) Activate(ctx context.Context, token string) error {
	query := `UPDATE users SET is_active = true WHERE activation_token = $1`
	res, err := r.DB.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return errs.BaseErr("user not found or token invalid", err)
	}
	return nil
}
