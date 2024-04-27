-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE users_role_enum AS ENUM ('user', 'admin', 'moderator');

CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text UNIQUE NOT NULL,
  password TEXT NOT NULL,
  role users_role_enum NOT NULL DEFAULT 'user'::users_role_enum
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "pgcrypto";

DROP TYPE IF EXISTS users_role_enum;

DROP TABLE users;
-- +goose StatementEnd
