-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS weeks ADD COLUMN finals_date timestamptz;
ALTER TABLE IF EXISTS weeks ADD COLUMN is_finals boolean DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS weeks DROP COLUMN IF EXISTS finals_date;
ALTER TABLE IF EXISTS weeks DROP COLUMN IF EXISTS is_finals;
-- +goose StatementEnd
