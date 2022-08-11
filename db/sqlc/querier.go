// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"
)

type Querier interface {
	CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetGroupsByOwner(ctx context.Context, arg GetGroupsByOwnerParams) ([]Group, error)
	GetUser(ctx context.Context, username string) (User, error)
}

var _ Querier = (*Queries)(nil)