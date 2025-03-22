-- name: GetPost :one
SELECT * FROM posts WHERE id = ? LIMIT 1;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY created_at DESC;

-- name: CreatePost :one
INSERT INTO posts (id, title, body, author) VALUES (?, ?, ?, ?) RETURNING *;

-- name: UpdatePost :one
UPDATE posts SET title = ?, body = ?, author = ? WHERE id = ? RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = ?;
