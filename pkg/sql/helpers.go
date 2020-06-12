package sql

import (
	"context"
	pgxInterfaces "github.com/nickeskov/db_forum/pkg/sql/pgx"
	pgx4Interfaces "github.com/nickeskov/db_forum/pkg/sql/pgx/v4"
	"github.com/pkg/errors"
)

func FinishPgxTransaction(tx pgxInterfaces.Transactioner, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Wrap(err, rollbackErr.Error())
		}
	} else if commitErr := tx.Commit(); commitErr != nil {
		err = errors.Wrap(commitErr, commitErr.Error())
	}

	return err
}

func FinishPgx4Transaction(ctx context.Context, tx pgx4Interfaces.Transactioner, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			err = errors.Wrap(err, rollbackErr.Error())
		}
	} else if commitErr := tx.Commit(ctx); commitErr != nil {
		err = errors.Wrap(commitErr, commitErr.Error())
	}

	return err
}
