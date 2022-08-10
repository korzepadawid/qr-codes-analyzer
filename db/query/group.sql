-- name: CreateGroup :one
INSERT INTO groups (owner, title, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGroupsByOwner :many
SELECT g.*
FROM users u
         JOIN groups g on u.username = g.owner AND u.username = sqlc.arg(owner)
LIMIT $1 OFFSET $2;