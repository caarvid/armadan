-- name: GetWeek :one
WITH week_course_data AS (
	SELECT
		w.id,
		w.nr,
		jsonb_build_object(
			'id', c.id,
			'name', c.name
		) AS course
	FROM weeks w
	LEFT JOIN courses c ON c.id = w.course_id
	GROUP BY w.id, w.nr, c.id, c.name
), week_tee_data AS (
	SELECT
		w.id,
		jsonb_build_object(
			'id', t.id,
			'name', t.name
		) AS tee
	FROM weeks w
	LEFT JOIN tees t ON t.id = w.tee_id
	GROUP BY w.id, t.id, t.name
)
SELECT cd.id, cd.nr, cd.course, td.tee FROM week_course_data cd JOIN week_tee_data td USING (id) WHERE id = $1::UUID;

-- name: GetWeeks :many
WITH week_course_data AS (
	SELECT
		w.id,
		w.nr,
		jsonb_build_object(
			'id', c.id,
			'name', c.name
		) AS course
	FROM weeks w
	LEFT JOIN courses c ON c.id = w.course_id
	GROUP BY w.id, w.nr, c.id, c.name
), week_tee_data AS (
	SELECT
		w.id,
		jsonb_build_object(
			'id', t.id,
			'name', t.name
		) AS tee
	FROM weeks w
	LEFT JOIN tees t ON t.id = w.tee_id
	GROUP BY w.id, t.id, t.name
)
SELECT cd.id, cd.nr, cd.course, td.tee FROM week_course_data cd JOIN week_tee_data td USING (id) ORDER BY nr ASC;

-- name: CreateWeek :one
INSERT INTO weeks (nr, course_id, tee_id) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateWeek :one
UPDATE weeks SET nr = $1, course_id = $2, tee_id = $3 WHERE id = $4 RETURNING *;

-- name: DeleteWeek :exec
DELETE FROM weeks WHERE id = $1;
