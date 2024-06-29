-- name: GetToken :one
SELECT * FROM password_reset_tokens WHERE hash = $1;

-- name: CreateToken :one
INSERT INTO password_reset_tokens (hash, user_id, expires_at) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteToken :exec
DELETE FROM password_reset_tokens WHERE user_id = $1; 
