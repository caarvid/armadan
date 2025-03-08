-- name: GetToken :one
SELECT * FROM password_reset_tokens WHERE token = ?;

-- name: CreateToken :one
INSERT INTO password_reset_tokens (token, user_id, expires_at) VALUES (?, ?, ?) RETURNING *;

-- name: DeleteToken :exec
DELETE FROM password_reset_tokens WHERE user_id = ?; 
