package db

import (
	"context"
	"fmt"
	"github.com/korzepadawid/qr-codes-analyzer/ipapi"
)

type UpdateGroupTxParams struct {
	Title       string
	Description string
	Owner       string
	ID          int64
}

type IncrementRedirectEntriesTxParams struct {
	UUID      string
	IPv4      string
	IPDetails *ipapi.IPDetails
}

// Transactor contains method signatures of db transactions
type Transactor interface {
	UpdateGroupTx(context.Context, UpdateGroupTxParams) (Group, error)

	IncrementRedirectEntriesTx(context.Context, IncrementRedirectEntriesTxParams) error
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

func (s *SQLStore) IncrementRedirectEntriesTx(ctx context.Context, params IncrementRedirectEntriesTxParams) error {
	return s.execTx(ctx, func(queries *Queries) error {
		if err := queries.IncrementQRCodeEntries(ctx, params.UUID); err != nil {
			return err
		}

		_, err := queries.CreateRedirect(ctx, CreateRedirectParams{
			QrCodeUuid:    params.UUID,
			Ipv4:          params.IPv4,
			Isp:           params.IPDetails.ISP,
			AutonomousSys: params.IPDetails.AS,
			Lat:           fmt.Sprintf("%9.5f", params.IPDetails.Lat),
			Lon:           fmt.Sprintf("%9.5f", params.IPDetails.Lon),
			City:          params.IPDetails.City,
			Country:       params.IPDetails.Country,
		})

		if err != nil {
			return err
		}

		return nil
	})
}
