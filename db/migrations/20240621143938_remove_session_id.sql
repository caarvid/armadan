-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS user_sessions DROP COLUMN IF EXISTS session_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS user_sessions ADD COLUMN session_id text NOT NULL UNIQUE;
-- +goose StatementEnd
