-- name: GetWeek :one
SELECT * FROM week_details WHERE id = $1::UUID;

-- name: GetWeeks :many
SELECT * FROM week_details ORDER BY nr ASC;

-- name: CreateWeek :one
INSERT INTO weeks (nr, course_id, tee_id, is_finals, finals_date) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateWeek :one
UPDATE weeks SET nr = $1, course_id = $2, tee_id = $3 WHERE id = $4 RETURNING *;

-- name: DeleteWeek :exec
DELETE FROM weeks WHERE id = $1;
