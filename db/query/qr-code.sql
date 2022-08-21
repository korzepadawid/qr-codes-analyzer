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