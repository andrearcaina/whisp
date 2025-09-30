-- name: ListMessages :many
SELECT * FROM messages
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateMessage :one
INSERT INTO messages (message, image_url, gif_url)
VALUES ($1, $2, $3)
RETURNING *;
