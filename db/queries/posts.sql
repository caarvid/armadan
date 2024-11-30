-- name: GetPost :one
SELECT * FROM posts WHERE id = $1 LIMIT 1;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY created_at DESC;

-- name: CreatePost :one
INSERT INTO posts (title, body, author) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePost :one
UPDATE posts SET title = $1, body = $2, author = $3 WHERE id = $4 RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;
