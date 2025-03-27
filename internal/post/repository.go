package post

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/phsaurav/go_echo_base/internal/database"
	errs "github.com/phsaurav/go_echo_base/pkg/error"
)

// Repo is a concrete implementation of the post repository.
// It must satisfy the consumer-side Repository interface.
type Repo struct {
	DB *sql.DB
}

// NewRepo creates a new repository instance.
func NewRepo(db database.Service) *Repo {
	return &Repo{DB: db.DB()}
}

// Ensure Repo implements the consumer-side Repository interface.
var _ Repository = (*Repo)(nil)

// Create inserts a new post into the database.
func (r *Repo) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRowContext(ctx, query, p.Content, p.Title, p.UserID, pq.Array(p.Tags)).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return errs.InternalServerError(err)
	}
	return nil
}

// GetByID retrieves a post by its ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`
	p := new(Post)
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Content, &p.UserID, pq.Array(&p.Tags),
		&p.CreatedAt, &p.UpdatedAt, &p.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound(err)
		}
		return nil, err
	}
	return p, nil
}
