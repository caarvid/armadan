-- name: GetResetPasswordToken :one
SELECT * FROM password_reset_tokens WHERE token = ?;

-- name: CreateResetPasswordToken :one
INSERT INTO password_reset_tokens (token, user_id, expires_at) VALUES (?, ?, ?) RETURNING *;

-- name: DeleteResetPasswordToken :exec
DELETE FROM password_reset_tokens WHERE token = ?; 
