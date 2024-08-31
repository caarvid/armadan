-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS results DROP COLUMN IF EXISTS has_results;
ALTER TABLE IF EXISTS results ADD COLUMN published boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS results DROP COLUMN IF EXISTS published;
ALTER TABLE IF EXISTS results ADD COLUMN has_results boolean NOT NULL DEFAULT false;
-- +goose StatementEnd
