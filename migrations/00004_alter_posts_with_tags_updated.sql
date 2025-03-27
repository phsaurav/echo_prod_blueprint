-- +goose Up
-- +goose StatementBegin
ALTER TABLE posts
  ADD COLUMN IF NOT EXISTS tags VARCHAR(100)[];
  
ALTER TABLE posts
  ADD COLUMN IF NOT EXISTS updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN IF EXISTS tags;

ALTER TABLE posts DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
