package user

import (
	"context"
	"database/sql"
	"github.com/phsaurav/go_echo_base/config"
	errs "github.com/phsaurav/go_echo_base/pkg/error"

	"github.com/phsaurav/go_echo_base/internal/user/models"
)

type UserRepo interface {
	Create(context.Context, *sql.Tx, *models.User) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) UserRepo {
	return &repository{db: db}
}

func (repo *repository) Create(ctx context.Context, tx *sql.Tx, user *models.User) error {
	query := `
		INSERT INTO users (username, models.Password, email, role_id) VALUES 
    ($1, $2, $3, (SELECT id FROM roles WHERE name = $4))
    RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	//role := user.Role.Name
	//if role == "" {
	//	role = "user"
	//}

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return errs.BasicErr("Email already exists")
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return errs.BasicErr("Username already exists")
		default:
			return err
		}
	}

	return nil
}
