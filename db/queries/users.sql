-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserById :one
SELECT id, email, role FROM users WHERE id = $1;

-- name: GetUsers :many
SELECT id, email, role FROM users;

-- name: CreateUser :one
INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *;

-- name: UpdateUserRole :one
UPDATE users SET role = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users SET email = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET password = $1 WHERE id = $2 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
