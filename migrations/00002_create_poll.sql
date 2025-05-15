-- +goose Up
-- +goose StatementBegin

-- Polls table
CREATE TABLE IF NOT EXISTS polls (
  id SERIAL PRIMARY KEY,
  question VARCHAR(255) NOT NULL,
  user_id INTEGER NOT NULL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Poll options table
CREATE TABLE IF NOT EXISTS poll_options (
  id SERIAL PRIMARY KEY,
  poll_id INTEGER NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
  text VARCHAR(255) NOT NULL
);

-- Poll votes table
CREATE TABLE IF NOT EXISTS poll_votes (
  id SERIAL PRIMARY KEY,
  poll_id INTEGER NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
  option_id INTEGER NOT NULL REFERENCES poll_options(id) ON DELETE CASCADE,
  user_id INTEGER NOT NULL REFERENCES users(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(poll_id, user_id) -- Ensures one vote per user per poll
);

-- Create indexes for performance
CREATE INDEX idx_poll_votes_poll_id ON poll_votes(poll_id);
CREATE INDEX idx_poll_votes_user_id ON poll_votes(user_id);
CREATE INDEX idx_poll_options_poll_id ON poll_options(poll_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS polls;
DROP TABLE IF EXISTS poll_options;
DROP TABLE IF EXISTS poll_votes;

-- +goose StatementEnd