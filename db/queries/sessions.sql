-- name: GetSessionByToken :one
SELECT 
  s.*,
  u.email,
  u.user_role
FROM sessions s 
JOIN users u ON u.id = s.user_id 
WHERE token = ?; 

-- name: CreateSession :one
INSERT INTO sessions (token, user_id, expires_at) VALUES (?, ?, ?) RETURNING *;   

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = ?;
