package v4

import (
	"context"
	"github.com/pkg/errors"
)

func FinishPgx4Transaction(ctx context.Context, tx Transactioner, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			err = errors.Wrap(err, rollbackErr.Error())
		}
	} else if commitErr := tx.Commit(ctx); commitErr != nil {
		err = errors.Wrap(commitErr, commitErr.Error())
	}

	return err
}
