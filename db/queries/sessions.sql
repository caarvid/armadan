-- name: GetSessionByToken :one
SELECT * FROM user_sessions WHERE token = $1; 

-- name: CreateSession :one
INSERT INTO user_sessions (token, user_id, expires_at, is_active) VALUES ($1, $2, $3, $4) RETURNING *;   

-- name: DeleteSession :exec
DELETE FROM user_sessions WHERE token = $1;
