package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	Transactor
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

type txFunc func(*Queries) error

func (s *SQLStore) execTx(ctx context.Context, fn txFunc) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error %s %s", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
