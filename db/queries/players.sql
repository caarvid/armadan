-- name: GetPlayer :one
SELECT 
  p.id, 
  p.first_name, 
  p.last_name, 
  p.points, 
  u.email,
  u.id::UUID AS user_id
FROM players p LEFT JOIN users u ON u.id = p.user_id WHERE p.id = $1::UUID;

-- name: GetPlayers :many
SELECT 
  p.id, 
  p.first_name, 
  p.last_name, 
  p.points, 
  u.email,
  u.id::UUID AS user_id
FROM players p LEFT JOIN users u ON u.id = p.user_id ORDER BY p.last_name ASC;

-- name: CreatePlayer :one
INSERT INTO players (first_name, last_name, user_id) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePlayer :one
UPDATE players SET first_name = $1, last_name = $2 WHERE id = $3 RETURNING *; 

-- name: DeletePlayer :exec
DELETE FROM players WHERE id = $1;
