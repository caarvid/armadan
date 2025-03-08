-- name: GetCourse :one
SELECT * FROM course_details c WHERE c.id = ?;

-- name: GetCourses :many
SELECT * FROM course_details c;

-- name: CreateCourse :one
INSERT INTO courses (id, name, par) VALUES (?, ?, ?) RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = ?;

-- name: UpdateCourse :one
UPDATE courses SET name = ?, par = ? WHERE id = ? RETURNING *;

-- name: CreateHoles :one
INSERT INTO holes (id, nr, par, stroke_index, course_id) VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateHoles :exec
UPDATE holes SET nr = ?, par = ?, stroke_index = ? WHERE id = ?;

-- name: CreateTees :one
INSERT INTO tees (id, name, slope, cr, course_id) VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateTees :exec
UPDATE tees SET name = ?, slope = ?, cr = ? WHERE id = ?;

-- name: GetTeesByCourse :many
SELECT * FROM tees WHERE course_id = ?;

-- name: DeleteTee :exec
DELETE FROM tees WHERE id = ?;

