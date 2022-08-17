-- name: CreateGroup :one
INSERT INTO groups (owner, title, description)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetGroupsByOwner :many
SELECT g.*
FROM users u
         JOIN groups g on u.username = g.owner AND u.username = sqlc.arg(owner)
ORDER BY g.created_at DESC LIMIT $1
OFFSET $2;

-- name: GetGroupsCountByOwner :one
SELECT COUNT(*)
FROM users u
         JOIN groups g on u.username = g.owner AND u.username = sqlc.arg(owner);

-- name: GetGroupByOwnerAndID :one
SELECT *
FROM groups
WHERE owner = sqlc.arg(owner)
  AND id = sqlc.arg(group_id) LIMIT 1;

-- name: DeleteGroupByOwnerAndID :exec
DELETE
FROM groups
WHERE id = sqlc.arg(group_id)
  and owner = sqlc.arg(owner);