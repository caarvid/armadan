-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: GetUserById :one
SELECT * FROM users WHERE id = ?;

-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (id, email, password) VALUES (?, ?, ?) RETURNING *;

-- name: UpdateUserRole :one
UPDATE users SET user_role = ? WHERE id = ? RETURNING *;

-- name: UpdateUserEmail :one
UPDATE users SET email = ? WHERE id = ? RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET password = ? WHERE id = ? RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;
