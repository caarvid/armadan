-- name: GetCourse :one
WITH tee_data AS (
  SELECT 
    t.course_id, 
    COALESCE(jsonb_agg(to_jsonb(t)) FILTER (WHERE t.course_id IS NOT NULL), '[]') AS tee_agg 
  FROM tees t 
  GROUP BY t.course_id
), hole_data AS (
  SELECT 
    h.course_id, 
    COALESCE(jsonb_agg(to_jsonb(h) ORDER BY h.nr) FILTER (WHERE h.course_id IS NOT NULL), '[]') AS hole_agg 
  FROM holes h 
  GROUP BY h.course_id
)
SELECT
  c.id,
  c.name,
  c.par,
  t.tee_agg AS tees,
  h.hole_agg AS holes
FROM courses c
LEFT JOIN tee_data t ON t.course_id = c.id
LEFT JOIN hole_data h ON h.course_id = c.id
WHERE c.id = $1
GROUP BY c.id, t.tee_agg, h.hole_agg;

-- name: GetCourses :many
WITH tee_data AS (
  SELECT 
    t.course_id, 
    COALESCE(jsonb_agg(to_jsonb(t)) FILTER (WHERE t.course_id IS NOT NULL), '[]') AS tee_agg 
  FROM tees t 
  GROUP BY t.course_id
), hole_data AS (
  SELECT 
    h.course_id, 
    COALESCE(jsonb_agg(to_jsonb(h) ORDER BY h.nr) FILTER (WHERE h.course_id IS NOT NULL), '[]') AS hole_agg 
  FROM holes h 
  GROUP BY h.course_id
)
SELECT
  c.id,
  c.name,
  c.par,
  t.tee_agg AS tees,
  h.hole_agg AS holes
FROM courses c
LEFT JOIN tee_data t ON t.course_id = c.id
LEFT JOIN hole_data h ON h.course_id = c.id
GROUP BY c.id, t.tee_agg, h.hole_agg;


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
