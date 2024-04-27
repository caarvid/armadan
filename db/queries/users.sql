-- name: CreateUser :one
INSERT INTO users (email, password) VALUES ($1, crypt(@password, gen_salt('bf', 11))) RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users SET email = $1 WHERE id = $2 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
