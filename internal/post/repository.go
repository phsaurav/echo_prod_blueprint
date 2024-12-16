package post

import (
	"context"
	"database/sql"

	"github.com/phsaurav/go_echo_base/config"
	"github.com/phsaurav/go_echo_base/internal/post/models"
)

type PostRepo interface {
	Create(ctx context.Context, post *models.Post) error
	GetByID(ctx context.Context, id int64) (*models.Post, error)
	// Add other methods as needed
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) PostRepo {
	return &repository{db: db}
}

func (repo *repository) Create(ctx context.Context, post *models.Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	err := repo.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	// Implement the GetByID method
	return nil, nil
}
