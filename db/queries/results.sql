-- name: GetLeaderboardSummary :many 
WITH winner_data AS (
  SELECT w.id, w.points, w.player_id, w.result_id FROM winners w
), result_data AS (
  SELECT r.id, r.week_id, r.published FROM results r
)
SELECT wk.id, wk.nr, COALESCE(wd.points, 0) AS points, COALESCE(rd.published, false) as has_results FROM weeks wk
LEFT JOIN result_data rd ON rd.week_id = wk.id 
LEFT JOIN winner_data wd ON rd.id = wd.result_id AND wd.player_id = $1::UUID 
GROUP BY wk.id, wd.points, rd.published
ORDER BY wk.nr ASC;

-- name: GetManageResultView :many
WITH result_data AS (
  SELECT r.id, r.published, r.week_id 
  FROM results r
), participant_counts AS (
  SELECT r.result_id, COUNT(r.id) AS participants
  FROM rounds r
  GROUP BY r.result_id
), winner_counts AS (
  SELECT w.result_id, COUNT(w.id) AS winners
  FROM winners w
  GROUP BY w.result_id
)
SELECT 
  wd.id,
  wd.nr,
  wd.is_finals,
  wd.course_name,
  wd.tee_name,
  rd.id::UUID AS result_id,
  COALESCE(rd.published, false) AS published,
  COALESCE(pc.participants, 0) AS participants,
  COALESCE(wc.winners, 0) AS winners
FROM week_details wd
LEFT JOIN result_data rd ON rd.week_id = wd.id
LEFT JOIN participant_counts pc ON pc.result_id = rd.id
LEFT JOIN winner_counts wc ON wc.result_id = rd.id
ORDER BY wd.nr ASC;

-- name: GetResultById :one
WITH week_data AS (
  SELECT w.id, w.nr, w.course_id, w.tee_id FROM weeks w
), tee_data AS (
  SELECT t.id, t.slope, t.cr FROM tees t
)
SELECT
  r.id,
  td.slope,
  td.cr,
  wd.nr::integer as week_nr,
  wd.course_id::UUID as course_id
FROM results r
LEFT JOIN week_data wd ON wd.id = r.week_id
LEFT JOIN tee_data td ON td.id = wd.tee_id
WHERE r.id = $1;

-- name: GetRoundsByResultId :many
WITH player_data AS (
  SELECT p.id, p.first_name, p.last_name, p.hcp FROM players p
)
SELECT 
  r.*, 
  r.net_in + r.net_out AS net_total,
  r.gross_in + r.gross_out AS gross_total,
  p.first_name, 
  p.last_name, 
  p.hcp
FROM rounds r
LEFT JOIN player_data p ON p.id = r.player_id
WHERE r.result_id = $1
ORDER BY net_total ASC;

-- name: CreateResult :one 
INSERT INTO results (week_id) VALUES ($1) RETURNING *;

-- name: CreateRound :one
INSERT INTO rounds (player_id, result_id, net_in, net_out, gross_in, gross_out, new_hcp, old_hcp) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: CreateScores :copyfrom
INSERT INTO scores (round_id, hole_id, strokes) VALUES ($1, $2, $3);

-- name: DeleteResult :exec
DELETE FROM results WHERE id = $1;

-- name: DeleteRound :exec
DELETE FROM rounds WHERE id = $1;
