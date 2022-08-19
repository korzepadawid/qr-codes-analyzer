// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: encode-code.sql

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