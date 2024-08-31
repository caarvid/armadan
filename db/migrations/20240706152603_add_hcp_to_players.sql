-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS players ADD COLUMN hcp decimal(3, 1) NOT NULL DEFAULT 0; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS players DROP COLUMN IF EXISTS hcp;
-- +goose StatementEnd
