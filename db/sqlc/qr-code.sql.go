// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: qr-code.sql

package db

import (
	"context"
)

const createQRCode = `-- name: CreateQRCode :one
INSERT INTO qr_codes (owner, group_id, redirection_url, title, description, storage_url, uuid)
VALUES ($1, $2, $3, $4, $5,
        $6, $7)
RETURNING uuid, owner, group_id, usages_count, redirection_url, title, description, storage_url, created_at
`

type CreateQRCodeParams struct {
	Owner          string `json:"owner"`
	GroupID        int64  `json:"group_id"`
	RedirectionUrl string `json:"redirection_url"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	StorageUrl     string `json:"storage_url"`
	Uuid           string `json:"uuid"`
}

func (q *Queries) CreateQRCode(ctx context.Context, arg CreateQRCodeParams) (QrCode, error) {
	row := q.db.QueryRowContext(ctx, createQRCode,
		arg.Owner,
		arg.GroupID,
		arg.RedirectionUrl,
		arg.Title,
		arg.Description,
		arg.StorageUrl,
		arg.Uuid,
	)
	var i QrCode
	err := row.Scan(
		&i.Uuid,
		&i.Owner,
		&i.GroupID,
		&i.UsagesCount,
		&i.RedirectionUrl,
		&i.Title,
		&i.Description,
		&i.StorageUrl,
		&i.CreatedAt,
	)
	return i, err
}

const deleteQRCode = `-- name: DeleteQRCode :exec
DELETE
FROM qr_codes
WHERE uuid = $1
  AND owner = $2
`

type DeleteQRCodeParams struct {
	Uuid  string `json:"uuid"`
	Owner string `json:"owner"`
}

func (q *Queries) DeleteQRCode(ctx context.Context, arg DeleteQRCodeParams) error {
	_, err := q.db.ExecContext(ctx, deleteQRCode, arg.Uuid, arg.Owner)
	return err
}

const getQRCode = `-- name: GetQRCode :one
SELECT uuid, owner, group_id, usages_count, redirection_url, title, description, storage_url, created_at
FROM qr_codes
WHERE uuid = $1
LIMIT 1
`

func (q *Queries) GetQRCode(ctx context.Context, uuid string) (QrCode, error) {
	row := q.db.QueryRowContext(ctx, getQRCode, uuid)
	var i QrCode
	err := row.Scan(
		&i.Uuid,
		&i.Owner,
		&i.GroupID,
		&i.UsagesCount,
		&i.RedirectionUrl,
		&i.Title,
		&i.Description,
		&i.StorageUrl,
		&i.CreatedAt,
	)
	return i, err
}

const getQRCodeForUpdate = `-- name: GetQRCodeForUpdate :one
SELECT uuid, owner, group_id, usages_count, redirection_url, title, description, storage_url, created_at
FROM qr_codes
WHERE uuid = $1
LIMIT 1 FOR NO KEY UPDATE
`

func (q *Queries) GetQRCodeForUpdate(ctx context.Context, uuid string) (QrCode, error) {
	row := q.db.QueryRowContext(ctx, getQRCodeForUpdate, uuid)
	var i QrCode
	err := row.Scan(
		&i.Uuid,
		&i.Owner,
		&i.GroupID,
		&i.UsagesCount,
		&i.RedirectionUrl,
		&i.Title,
		&i.Description,
		&i.StorageUrl,
		&i.CreatedAt,
	)
	return i, err
}

const getQRCodeForUpdateTitleAndDesc = `-- name: GetQRCodeForUpdateTitleAndDesc :one
SELECT uuid, owner, group_id, usages_count, redirection_url, title, description, storage_url, created_at
FROM qr_codes
WHERE uuid = $1
  AND owner = $2
LIMIT 1 FOR NO KEY UPDATE
`

type GetQRCodeForUpdateTitleAndDescParams struct {
	Uuid  string `json:"uuid"`
	Owner string `json:"owner"`
}

func (q *Queries) GetQRCodeForUpdateTitleAndDesc(ctx context.Context, arg GetQRCodeForUpdateTitleAndDescParams) (QrCode, error) {
	row := q.db.QueryRowContext(ctx, getQRCodeForUpdateTitleAndDesc, arg.Uuid, arg.Owner)
	var i QrCode
	err := row.Scan(
		&i.Uuid,
		&i.Owner,
		&i.GroupID,
		&i.UsagesCount,
		&i.RedirectionUrl,
		&i.Title,
		&i.Description,
		&i.StorageUrl,
		&i.CreatedAt,
	)
	return i, err
}

const getQRCodesCountByGroupAndOwner = `-- name: GetQRCodesCountByGroupAndOwner :one
SELECT COUNT(*)
FROM groups g
         JOIN qr_codes qc on g.id = qc.group_id
WHERE g.id = $1
  AND g.owner = $2
`

type GetQRCodesCountByGroupAndOwnerParams struct {
	GroupID int64  `json:"group_id"`
	Owner   string `json:"owner"`
}

func (q *Queries) GetQRCodesCountByGroupAndOwner(ctx context.Context, arg GetQRCodesCountByGroupAndOwnerParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getQRCodesCountByGroupAndOwner, arg.GroupID, arg.Owner)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getQRCodesPageByGroupAndOwner = `-- name: GetQRCodesPageByGroupAndOwner :many
SELECT qc.uuid, qc.owner, qc.group_id, qc.usages_count, qc.redirection_url, qc.title, qc.description, qc.storage_url, qc.created_at
FROM groups g
         JOIN qr_codes qc on g.id = qc.group_id
WHERE g.id = $3
  AND g.owner = $4
LIMIT $1 OFFSET $2
`

type GetQRCodesPageByGroupAndOwnerParams struct {
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	GroupID int64  `json:"group_id"`
	Owner   string `json:"owner"`
}

func (q *Queries) GetQRCodesPageByGroupAndOwner(ctx context.Context, arg GetQRCodesPageByGroupAndOwnerParams) ([]QrCode, error) {
	rows, err := q.db.QueryContext(ctx, getQRCodesPageByGroupAndOwner,
		arg.Limit,
		arg.Offset,
		arg.GroupID,
		arg.Owner,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []QrCode{}
	for rows.Next() {
		var i QrCode
		if err := rows.Scan(
			&i.Uuid,
			&i.Owner,
			&i.GroupID,
			&i.UsagesCount,
			&i.RedirectionUrl,
			&i.Title,
			&i.Description,
			&i.StorageUrl,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const incrementQRCodeEntries = `-- name: IncrementQRCodeEntries :exec
UPDATE qr_codes
SET usages_count = usages_count + 1
WHERE uuid = $1
`

func (q *Queries) IncrementQRCodeEntries(ctx context.Context, uuid string) error {
	_, err := q.db.ExecContext(ctx, incrementQRCodeEntries, uuid)
	return err
}

const updateQRCode = `-- name: UpdateQRCode :exec
UPDATE qr_codes
SET title       = $1,
    description = $2
WHERE uuid = $3
  AND owner = $4
`

type UpdateQRCodeParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Uuid        string `json:"uuid"`
	Owner       string `json:"owner"`
}

func (q *Queries) UpdateQRCode(ctx context.Context, arg UpdateQRCodeParams) error {
	_, err := q.db.ExecContext(ctx, updateQRCode,
		arg.Title,
		arg.Description,
		arg.Uuid,
		arg.Owner,
	)
	return err
}
