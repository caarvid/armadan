-- name: GetCourse :one
SELECT
  c.id,
  c.name,
  c.par,
  COALESCE(
    (
      SELECT jsonb_agg(to_jsonb(t))
      FROM tees t 
      WHERE t.course_id = c.id
    ), '[]'
  )::jsonb AS tees,
  COALESCE(
    (
      SELECT jsonb_agg(to_jsonb(h) ORDER BY h.nr)
      FROM holes h 
      WHERE h.course_id = c.id
  ), '[]'
  )::jsonb AS holes
FROM courses c
WHERE c.id=$1;

-- name: GetCourses :many
SELECT
  c.id,
  c.name,
  c.par,
  COALESCE(
    (
      SELECT jsonb_agg(to_jsonb(t))
      FROM tees t 
      WHERE t.course_id = c.id
    ), '[]'
  )::jsonb AS tees,
  COALESCE(
    (
      SELECT jsonb_agg(to_jsonb(h) ORDER BY h.nr)
      FROM holes h 
      WHERE h.course_id = c.id
  ), '[]'
  )::jsonb AS holes
FROM courses c
GROUP BY c.id;

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
