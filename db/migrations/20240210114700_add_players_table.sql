-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS players (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  first_name text NOT NULL,
  last_name text NOT NULL,
  points integer NOT NULL DEFAULT 0,
  user_id uuid NOT NULL,
  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE players;
-- +goose StatementEnd
