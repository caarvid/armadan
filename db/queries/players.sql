-- name: GetPlayer :one
SELECT 
  p.id, 
  p.first_name, 
  p.last_name, 
  p.points, 
  p.hcp,
  u.email,
  u.id::UUID AS user_id
FROM players p LEFT JOIN users u ON u.id = p.user_id WHERE p.id = $1::UUID;

-- name: GetPlayers :many
SELECT 
  p.id, 
  p.first_name, 
  p.last_name, 
  p.points, 
  p.hcp,
  u.email,
  u.id::UUID AS user_id
FROM players p LEFT JOIN users u ON u.id = p.user_id ORDER BY p.last_name ASC, p.first_name ASC;

-- name: GetLeaderboard :many
SELECT
  p.id,
  p.first_name,
  p.last_name,
  p.points,
  COUNT(rounds.id) as nr_of_rounds
FROM players p
LEFT JOIN rounds ON rounds.player_id = p.id
GROUP BY p.id
ORDER BY p.points DESC, nr_of_rounds DESC;

-- name: CreatePlayer :one
INSERT INTO players (first_name, last_name, hcp, user_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePlayer :one
UPDATE players SET first_name = $1, last_name = $2, hcp = $3 WHERE id = $4 RETURNING *; 

-- name: DeletePlayer :exec
DELETE FROM players WHERE id = $1;
