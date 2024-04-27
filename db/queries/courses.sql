-- name: GetCourse :one
WITH hole_data AS (
  SELECT
    c.id,
    c.name,
    c.par,
    COALESCE(jsonb_agg(to_jsonb(h) ORDER BY h.nr ASC) FILTER (WHERE h.course_id IS NOT NULL), '[]') AS holes
  FROM courses c
  LEFT JOIN holes h ON h.course_id = c.id
  GROUP BY c.id, c.name, c.par
), tee_data AS (
  SELECT
    c.id,
    COALESCE(jsonb_agg(to_jsonb(t)) FILTER (WHERE t.course_id IS NOT NULL), '[]') AS tees
  FROM courses c
  LEFT JOIN tees t ON t.course_id = c.id
  GROUP BY c.id, c.name, c.par
)
SELECT hd.id, hd.name, hd.par, hd.holes, td.tees FROM hole_data hd JOIN tee_data td USING (id) WHERE id = $1::UUID;

-- name: GetCourses :many
WITH hole_data AS (
  SELECT
    c.id,
    c.name,
    c.par,
    COALESCE(jsonb_agg(to_jsonb(h) ORDER BY h.nr ASC) FILTER (WHERE h.course_id IS NOT NULL), '[]') AS holes
  FROM courses c
  LEFT JOIN holes h ON h.course_id = c.id
  GROUP BY c.id, c.name, c.par
), tee_data AS (
  SELECT
    c.id,
    COALESCE(jsonb_agg(to_jsonb(t)) FILTER (WHERE t.course_id IS NOT NULL), '[]') AS tees
  FROM courses c
  LEFT JOIN tees t ON t.course_id = c.id
  GROUP BY c.id, c.name, c.par
)
SELECT hd.id, hd.name, hd.par, hd.holes, td.tees FROM hole_data hd JOIN tee_data td USING (id);


-- name: CreateCourse :one
INSERT INTO courses (name, par) VALUES ($1, $2) RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1;

-- name: UpdateCourse :one
UPDATE courses SET name = $1, par = $2 WHERE id = $3 RETURNING *;

-- name: CreateHoles :copyfrom
INSERT INTO holes (nr, par, index, course_id) VALUES ($1, $2, $3, $4);

-- name: UpdateHoles :batchexec
UPDATE holes SET nr = $1, par = $2, index = $3 WHERE id = $4;

-- name: CreateTees :copyfrom
INSERT INTO tees (name, slope, cr, course_id) values ($1, $2, $3, $4);

-- name: UpdateTees :batchexec
UPDATE tees SET name = $1, slope = $2, cr = $3 WHERE id = $4;

-- name: DeleteTee :exec
DELETE FROM tees WHERE id = $1;
