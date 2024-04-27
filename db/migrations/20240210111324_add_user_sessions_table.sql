-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id text NOT NULL,
  data jsonb NOT NULL,
  user_id uuid NOT NULL,
  CONSTRAINT fk_user_id
    FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_sessions;
-- +goose StatementEnd
