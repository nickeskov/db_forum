package codes

import (
	"github.com/jackc/pgconn"
	oldPgx "github.com/jackc/pgx"
	"github.com/nickeskov/db_forum/internal/pkg/models"
)

const (
	ErrCodeNotNull        = "23502"
	ErrCodeForeignKey     = "23503"
	ErrCodeUnique         = "23505"
	ErrCodeRaiseException = "P0001"
)

func ExtractPgx4ErrorCode(err error) error {
	if err == nil {
		return nil
	}

	pgxErr, ok := err.(*pgconn.PgError)
	if !ok {
		return models.ErrInvalid
	}
	return models.NewError(pgxErr.Code)
}

func ConvertToPgx4Error(err error) (*pgconn.PgError, bool) {
	pgError, ok := err.(*pgconn.PgError)
	return pgError, ok
}

func ExtractPgxErrorCode(err error) error {
	if err == nil {
		return nil
	}

	pgxErr, ok := err.(oldPgx.PgError)
	if !ok {
		return models.ErrInvalid
	}
	return models.NewError(pgxErr.Code)
}

func ConvertToPgxError(err error) (oldPgx.PgError, bool) {
	pgError, ok := err.(oldPgx.PgError)
	return pgError, ok
}
