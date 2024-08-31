-- name: GetSessionByToken :one
SELECT 
  us.id,
  us.user_id,
  us.is_active,
  us.expires_at,
  us.token,
  u.role,
  u.email
FROM user_sessions us LEFT JOIN users u ON u.id = us.user_id WHERE token = $1; 

-- name: CreateSession :one
INSERT INTO user_sessions (token, user_id, expires_at, is_active) VALUES ($1, $2, $3, $4) RETURNING *;   

-- name: DeleteSession :exec
DELETE FROM user_sessions WHERE token = $1;
