-- name: GetTeesByCourse :many
SELECT * FROM tees WHERE course_id = $1;