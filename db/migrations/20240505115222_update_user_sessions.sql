-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS user_sessions DROP COLUMN IF EXISTS data;
ALTER TABLE IF EXISTS user_sessions ADD COLUMN expires_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE IF EXISTS user_sessions ADD COLUMN is_active boolean NOT NULL;
ALTER TABLE IF EXISTS user_sessions ADD COLUMN token text NOT NULL UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS user_sessions DROP COLUMN IF EXISTS is_active;
ALTER TABLE IF EXISTS user_sessions DROP COLUMN IF EXISTS expires_at;
ALTER TABLE IF EXISTS user_sessions DROP COLUMN IF EXISTS token;
ALTER TABLE IF EXISTS user_sessions ADD COLUMN data jsonb NOT NULL;
-- +goose StatementEnd
