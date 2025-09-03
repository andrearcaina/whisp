-- name: ListMessages :many
SELECT * FROM messages
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateMessage :one
INSERT INTO messages (message)
VALUES ($1)
RETURNING *;
