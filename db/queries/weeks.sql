-- name: GetWeek :one
SELECT * FROM week_details WHERE id = ?;

-- name: GetWeeks :many
SELECT * FROM week_details ORDER BY nr ASC;

-- name: CreateWeek :one
INSERT INTO weeks (id, nr, course_id, tee_id, is_finals, finals_date) VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateWeek :one
UPDATE weeks SET nr = ?, course_id = ?, tee_id = ? WHERE id = ? RETURNING *;

-- name: DeleteWeek :exec
DELETE FROM weeks WHERE id = ?;
