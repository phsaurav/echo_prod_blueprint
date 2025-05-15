package poll

import (
	"context"
	"database/sql"
	"errors"

	"github.com/phsaurav/echo_prod_blueprint/internal/database"
	errs "github.com/phsaurav/echo_prod_blueprint/pkg/error"
)

// Repo is the concrete implementation of the poll repository.
type Repo struct {
	DB *sql.DB
}

// NewRepo creates a new poll repository instance.
func NewRepo(db database.Service) *Repo {
	return &Repo{DB: db.DB()}
}

var _ Repository = (*Repo)(nil)

// Create inserts a new poll and its options into the DB.
func (r *Repo) Create(ctx context.Context, p *Poll) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return errs.InternalServerError(err)
	}
	defer tx.Rollback()

	// Insert poll
	pollQuery := `
		INSERT INTO polls (question, created_at)
		VALUES ($1, NOW())
		RETURNING id, created_at
	`
	err = tx.QueryRowContext(ctx, pollQuery, p.Question).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return errs.InternalServerError(err)
	}

	// Insert options
	optionQuery := `INSERT INTO poll_options (poll_id, text) VALUES ($1, $2) RETURNING id`
	for i := range p.Options {
		err := tx.QueryRowContext(ctx, optionQuery, p.ID, p.Options[i].Text).Scan(&p.Options[i].ID)
		if err != nil {
			return errs.InternalServerError(err)
		}
		p.Options[i].PollID = p.ID
	}

	if err := tx.Commit(); err != nil {
		return errs.InternalServerError(err)
	}
	return nil
}

// GetByID fetches a poll and its options by poll ID.
func (r *Repo) GetByID(ctx context.Context, id int64) (*Poll, error) {
	query := `SELECT id, question, created_at FROM polls WHERE id = $1`
	p := new(Poll)
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Question, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound(err)
		}
		return nil, errs.InternalServerError(err)
	}

	// Fetch options
	optQuery := `SELECT id, poll_id, text FROM poll_options WHERE poll_id = $1`
	rows, err := r.DB.QueryContext(ctx, optQuery, p.ID)
	if err != nil {
		return nil, errs.InternalServerError(err)
	}
	defer rows.Close()

	var opts []Option
	for rows.Next() {
		var opt Option
		if err := rows.Scan(&opt.ID, &opt.PollID, &opt.Text); err != nil {
			return nil, errs.InternalServerError(err)
		}
		opts = append(opts, opt)
	}
	p.Options = opts
	return p, nil
}

// Vote records a user's vote for a specific poll option.
func (r *Repo) Vote(ctx context.Context, pollID, optionID, userID int64) error {
	voteQuery := `INSERT INTO poll_votes (poll_id, option_id, user_id, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.DB.ExecContext(ctx, voteQuery, pollID, optionID, userID)
	if err != nil {
		return errs.InternalServerError(err)
	}
	return nil
}

// GetResults fetches poll options and their vote counts for a poll.
func (r *Repo) GetResults(ctx context.Context, pollID int64) ([]Option, error) {
	query := `
		SELECT o.id, o.poll_id, o.text, COUNT(v.id) as votes
		FROM poll_options o
		LEFT JOIN poll_votes v ON o.id = v.option_id
		WHERE o.poll_id = $1
		GROUP BY o.id
		ORDER BY o.id
	`
	rows, err := r.DB.QueryContext(ctx, query, pollID)
	if err != nil {
		return nil, errs.InternalServerError(err)
	}
	defer rows.Close()

	var opts []Option
	for rows.Next() {
		var opt Option
		if err := rows.Scan(&opt.ID, &opt.PollID, &opt.Text, &opt.Votes); err != nil {
			return nil, errs.InternalServerError(err)
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

func (r *Repo) HasUserVoted(ctx context.Context, pollID int64, userID int64) (bool, error) {
	query := "SELECT 1 FROM poll_votes WHERE poll_id=$1 AND user_id=$2"
	row := r.DB.QueryRowContext(ctx, query, pollID, userID)
	var dummy int
	err := row.Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
