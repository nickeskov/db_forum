package sql

import "github.com/pkg/errors"

func FinishTransaction(tx Transactioner, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Wrap(err, rollbackErr.Error())
		}
	} else if commitErr := tx.Commit(); commitErr != nil {
		err = errors.Wrap(commitErr, commitErr.Error())
	}

	return err
}
