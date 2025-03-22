-- name: GetPlayerHcp :one 
SELECT new_hcp FROM hcp_changes WHERE player_id = ? ORDER BY datetime(valid_from) DESC LIMIT 1;

-- name: CreateHcpChange :one 
INSERT INTO hcp_changes (old_hcp, new_hcp, round_id, player_id, valid_from) VALUES (?, ?, ?, ?, ?) RETURNING *;
