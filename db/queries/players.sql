-- name: GetPlayer :one
SELECT * FROM players_extended p WHERE p.id = ?;

-- name: GetPlayers :many
SELECT * FROM players_extended p ORDER BY p.last_name ASC, p.first_name ASC;

-- name: GetLeaderboard :many
SELECT 
  p.*,
  count(r.id) as nr_of_rounds
FROM players_extended p
LEFT JOIN rounds r ON r.player_id = p.id
GROUP BY p.id
ORDER BY p.points DESC, nr_of_rounds DESC;

-- name: CreatePlayer :one
INSERT INTO players (id, first_name, last_name, user_id) VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdatePlayer :one
UPDATE players SET first_name = ?, last_name = ?  WHERE id = ? RETURNING *; 

