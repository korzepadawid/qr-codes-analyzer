// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: group.sql

package db

import (
	"context"
)

const createGroup = `-- name: CreateGroup :one
INSERT INTO groups (owner, title, description)
VALUES ($1, $2, $3) RETURNING id, owner, title, description, created_at
`

type CreateGroupParams struct {
	Owner       string `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, createGroup, arg.Owner, arg.Title, arg.Description)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const getGroupByOwnerAndID = `-- name: GetGroupByOwnerAndID :one
SELECT g.id, g.owner, g.title, g.description, g.created_at
FROM users u
         JOIN groups g ON g.owner = u.username
WHERE g.owner = $1
  AND g.id = $2
`

type GetGroupByOwnerAndIDParams struct {
	Owner   string `json:"owner"`
	GroupID int64  `json:"group_id"`
}

func (q *Queries) GetGroupByOwnerAndID(ctx context.Context, arg GetGroupByOwnerAndIDParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, getGroupByOwnerAndID, arg.Owner, arg.GroupID)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.Description,
		&i.CreatedAt,
	)
	return i, err
}

const getGroupsByOwner = `-- name: GetGroupsByOwner :many
SELECT g.id, g.owner, g.title, g.description, g.created_at
FROM users u
         JOIN groups g on u.username = g.owner AND u.username = $3
ORDER BY g.created_at DESC LIMIT $1
OFFSET $2
`

type GetGroupsByOwnerParams struct {
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Owner  string `json:"owner"`
}

func (q *Queries) GetGroupsByOwner(ctx context.Context, arg GetGroupsByOwnerParams) ([]Group, error) {
	rows, err := q.db.QueryContext(ctx, getGroupsByOwner, arg.Limit, arg.Offset, arg.Owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Group{}
	for rows.Next() {
		var i Group
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Title,
			&i.Description,
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

const getGroupsCountByOwner = `-- name: GetGroupsCountByOwner :one
SELECT COUNT(*)
FROM users u
         JOIN groups g on u.username = g.owner AND u.username = $1
`

func (q *Queries) GetGroupsCountByOwner(ctx context.Context, owner string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getGroupsCountByOwner, owner)
	var count int64
	err := row.Scan(&count)
	return count, err
}
