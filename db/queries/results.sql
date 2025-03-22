-- name: GetLeaderboardSummary :many 
SELECT
  weeks.*,
  coalesce(w.points, 0) AS points,
  coalesce(r.is_published, FALSE) AS has_results
FROM weeks
JOIN results r ON r.week_id = weeks.id
JOIN winners w ON w.week_id = weeks.id AND w.player_id = ?
ORDER BY weeks.nr;

-- name: GetManageResultView :many
SELECT
  w.id,
  w.nr,
  w.is_finals,
  w.course_name,
  w.tee_name,
  r.id AS result_id,
  coalesce(r.is_published, FALSE) as is_published,
  coalesce(rd.participants, 0) AS participants,
  coalesce(win.winners, 0) as winners
FROM week_details w
LEFT JOIN results r ON r.week_id = w.id
LEFT JOIN  (
  SELECT result_id, COUNT(*) as participants FROM rounds GROUP BY result_id
) rd ON rd.result_id = r.id
LEFT JOIN (
    SELECT week_id, COUNT(*) AS winners FROM winners GROUP BY week_id
) win ON win.week_id = w.id
GROUP BY w.id, r.id, win.winners
ORDER BY w.nr ASC;

-- name: GetResultById :one
SELECT
  r.*,
  w.nr as week_nr,
  w.start_date as week_start_date,
  w.end_date as week_end_date,
  w.course_id,
  t.slope,
  t.cr
FROM results r
JOIN weeks w ON w.id = r.week_id
JOIN tees t ON t.id = w.tee_id
WHERE r.id = ?;

-- name: GetRoundsByResultId :many
SELECT
  r.*,
  p.first_name,
  p.last_name
FROM full_rounds r
JOIN players p ON p.id = r.player_id
WHERE r.result_id = ?
ORDER BY r.net_total ASC;

-- name: GetRoundById :one
SELECT * FROM full_rounds WHERE id = ?;

-- name: GetRemainingPlayersByResultId :many
SELECT 
  p.*,
  cast(coalesce((
      SELECT h.new_hcp FROM hcp_changes h
      WHERE h.player_id = p.id AND datetime(h.valid_from) < datetime(w.end_date)
      ORDER BY datetime(h.valid_from) DESC
      LIMIT 1
  ), (
      SELECT h.new_hcp FROM hcp_changes h 
      WHERE h.player_id = p.id 
      ORDER BY datetime(h.valid_from) ASC
      LIMIT 1
  ), 36.0) as real) AS hcp
FROM players p
LEFT JOIN rounds r ON r.player_id = p.id AND r.result_id = ?
LEFT JOIN results res ON res.id = r.result_id
LEFT JOIN weeks w ON w.id = res.week_id
WHERE r.player_id IS NULL
ORDER BY p.last_name ASC, p.first_name ASC;

-- name: CreateResult :one 
INSERT INTO results (id, week_id) VALUES (?, ?) RETURNING *;

-- name: CreateRound :one
INSERT INTO rounds (id, player_id, result_id) VALUES (?, ?, ?) RETURNING *;

-- name: CreateRoundDetail :one 
INSERT INTO round_details (net_in, net_out, gross_in, gross_out, round_id) VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: CreateScores :one
INSERT INTO scores (round_id, hole_id, strokes) VALUES (?, ?, ?) RETURNING *;

-- name: CreateWinner :one
INSERT INTO winners (id, points, player_id, week_id) VALUES (?, ?, ?, ?) RETURNING *;

-- name: PublishRound :exec
UPDATE results SET is_published = TRUE WHERE id = ?;

-- name: DeleteResult :exec
DELETE FROM results WHERE id = ?;

-- name: DeleteRound :exec
DELETE FROM rounds WHERE id = ?;

-- name: DeleteWinnersByWeek :exec
DELETE FROM winners WHERE week_id = ?;
