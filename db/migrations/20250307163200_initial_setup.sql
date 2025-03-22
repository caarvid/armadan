-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
	id text PRIMARY KEY NOT NULL,
	title text NOT NULL,
	body text NOT NULL,
	author text NOT NULL,
	created_at text NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
	id text PRIMARY KEY NOT NULL,
	email text UNIQUE NOT NULL,
	password text NOT NULL,
	user_role text NOT NULL DEFAULT 'user' CHECK (user_role in ('user', 'moderator', 'admin'))
);

CREATE TABLE IF NOT EXISTS players (
	id text PRIMARY KEY NOT NULL,
	first_name text NOT NULL,
	last_name text NOT NULL,
	user_id text UNIQUE NOT NULL,
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions (
	id integer PRIMARY KEY,
	token text UNIQUE NOT NULL,
	expires_at text NOT NULL,
	user_id text NOT NULL,
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
	id integer PRIMARY KEY,
	token text UNIQUE NOT NULL,
	expires_at text NOT NULL,
	user_id text NOT NULL,
	CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens(token);

CREATE TABLE IF NOT EXISTS courses (
	id text PRIMARY KEY NOT NULL,
	name text UNIQUE NOT NULL,
	par integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS holes (
	id text PRIMARY KEY NOT NULL,
	nr integer NOT NULL,
	par integer NOT NULL,
	stroke_index integer NOT NULL,
	course_id text NOT NULL,
	CONSTRAINT fk_course_id FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_holes_course_id ON holes(course_id);

CREATE TABLE IF NOT EXISTS tees (
	id text PRIMARY KEY NOT NULL,
	name text NOT NULL,
	slope integer NOT NULL,
	cr REAL NOT NULL,
	course_id text NOT NULL,
	CONSTRAINT fk_course_id FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tees_course_id ON tees(course_id);

CREATE TABLE IF NOT EXISTS weeks (
	id text PRIMARY KEY NOT NULL,
	nr integer UNIQUE NOT NULL,
	is_finals integer NOT NULL DEFAULT FALSE CHECK (
		is_finals IN (
			FALSE,
			TRUE
		)
	),
	finals_date text,
	start_date text NOT NULL,
	end_date text NOT NULL,
	course_id text NOT NULL,
	tee_id text NOT NULL,
	CONSTRAINT fk_course_id FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
	CONSTRAINT fk_tee_id FOREIGN KEY (tee_id) REFERENCES tees(id)
);

CREATE TABLE IF NOT EXISTS results (
	id text PRIMARY KEY NOT NULL,
	is_published integer NOT NULL DEFAULT FALSE CHECK (
		is_published IN (
			FALSE,
			TRUE
		)
	),
	week_id text UNIQUE NOT NULL,
	CONSTRAINT fk_week_id FOREIGN KEY (week_id) REFERENCES weeks(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rounds (
	id text PRIMARY KEY NOT NULL,
	player_id text NOT NULL,
	result_id text NOT NULL,
	CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
	CONSTRAINT fk_result_id FOREIGN KEY (result_id) REFERENCES results(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS round_details (
	id integer PRIMARY KEY,
	net_in integer NOT NULL,
	net_out integer NOT NULL,
	net_total integer AS (net_in+net_out) STORED,
	gross_in integer NOT NULL,
	gross_out integer NOT NULL,
	gross_total integer AS (gross_in+gross_out) STORED,
	round_id text UNIQUE NOT NULL,
	CONSTRAINT fk_round_id FOREIGN KEY (round_id) REFERENCES rounds(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS hcp_changes (
	id integer PRIMARY KEY,
	new_hcp REAL NOT NULL,
	old_hcp REAL NOT NULL,
	valid_from text NOT NULL,
	round_id text UNIQUE,
	player_id text,
	CONSTRAINT fk_round_id FOREIGN KEY (round_id) REFERENCES rounds(id) ON DELETE CASCADE,
	CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS winners (
	id text PRIMARY KEY NOT NULL,
	points integer NOT NULL,
	week_id text NOT NULL,
	player_id text NOT NULL,
	CONSTRAINT fk_week_id FOREIGN KEY (week_id) REFERENCES weeks(id) ON DELETE CASCADE,
	CONSTRAINT fk_player_id FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS scores (
	id integer PRIMARY KEY,
	strokes integer NOT NULL,
	hole_id text NOT NULL,
	round_id text NOT NULL,
	CONSTRAINT fk_hole_id FOREIGN KEY (hole_id) REFERENCES holes(id) ON DELETE CASCADE,
	CONSTRAINT fk_round_id FOREIGN KEY (round_id) REFERENCES rounds(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;

DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS players;

DROP TABLE IF EXISTS sessions;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_token;

DROP TABLE IF EXISTS password_reset_tokens;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;

DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS holes;
DROP INDEX IF EXISTS idx_holes_course_id;
DROP TABLE IF EXISTS tees;
DROP INDEX IF EXISTS idx_tees_course_id;

DROP TABLE IF EXISTS weeks;
DROP TABLE IF EXISTS results;
DROP TABLE IF EXISTS rounds;
DROP TABLE IF EXISTS round_details;
DROP TABLE IF EXISTS hcp_changes;
DROP TABLE IF EXISTS winners;
DROP TABLE IF EXISTS scores;
-- +goose StatementEnd
