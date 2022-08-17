package db

import (
	"context"
)

type UpdateGroupTxParams struct {
	Title       string
	Description string
	Owner       string
	ID          int64
}

// Transactor contains method signatures of db transactions
type Transactor interface {
	UpdateGroupTx(context.Context, UpdateGroupTxParams) (Group, error)
}

func (s *SQLStore) UpdateGroupTx(ctx context.Context, params UpdateGroupTxParams) (Group, error) {
	var result Group
	err := s.execTx(ctx, func(queries *Queries) error {
		group, err := s.GetGroupByOwnerAndIDForUpdate(ctx, GetGroupByOwnerAndIDForUpdateParams{
			Owner:   params.Owner,
			GroupID: params.ID,
		})

		if err != nil {
			return err
		}

		updateArg := createUpdateArgs(group, params)

		result, err = s.UpdateGroupByOwnerAndID(ctx, updateArg)

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func createUpdateArgs(group Group, params UpdateGroupTxParams) UpdateGroupByOwnerAndIDParams {
	updateArg := UpdateGroupByOwnerAndIDParams{
		Title:       group.Title,
		Description: group.Description,
		ID:          params.ID,
		Owner:       params.Owner,
	}

	if len(params.Title) > 0 {
		updateArg.Title = params.Title
	}

	if len(params.Description) > 0 {
		updateArg.Description = params.Description
	}

	return updateArg
}
