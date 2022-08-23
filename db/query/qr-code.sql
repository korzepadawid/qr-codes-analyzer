-- name: CreateQRCode :one
INSERT INTO qr_codes (owner, group_id, redirection_url, title, description, storage_url, uuid)
VALUES (sqlc.arg(owner), sqlc.arg(group_id), sqlc.arg(redirection_url), sqlc.arg(title), sqlc.arg(description),
        sqlc.arg(storage_url), sqlc.arg(UUID))
RETURNING *;

-- name: GetQRCode :one
SELECT *
FROM qr_codes
WHERE uuid = sqlc.arg(UUID)
LIMIT 1;

-- name: GetQRCodeForUpdate :one
SELECT *
FROM qr_codes
WHERE uuid = sqlc.arg(UUID)
LIMIT 1 FOR NO KEY UPDATE;

-- name: IncrementQRCodeEntries :exec
UPDATE qr_codes
SET usages_count = usages_count + 1
WHERE uuid = sqlc.arg(UUID);

-- name: DeleteQRCode :exec
DELETE
FROM qr_codes
WHERE uuid = sqlc.arg(UUID)
  AND owner = sqlc.arg(owner);

-- name: GetQRCodesPageByGroupAndOwner :many
SELECT qc.*
FROM groups g
         JOIN qr_codes qc on g.id = qc.group_id
WHERE g.id = sqlc.arg(group_id)
  AND g.owner = sqlc.arg(owner)
LIMIT $1 OFFSET $2;

-- name: GetQRCodesCountByGroupAndOwner :one
SELECT COUNT(*)
FROM groups g
         JOIN qr_codes qc on g.id = qc.group_id
WHERE g.id = sqlc.arg(group_id)
  AND g.owner = sqlc.arg(owner);

-- name: UpdateQRCode :exec
UPDATE qr_codes
SET title       = $1,
    description = $2
WHERE uuid = $3
  AND owner = $4;

-- name: GetQRCodeForUpdateTitleAndDesc :one
SELECT *
FROM qr_codes
WHERE uuid = sqlc.arg(UUID)
  AND owner = sqlc.arg(owner)
LIMIT 1 FOR NO KEY UPDATE;
