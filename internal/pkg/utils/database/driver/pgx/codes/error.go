package codes

import (
	"github.com/jackc/pgx"
	"github.com/nickeskov/db_forum/internal/pkg/models"
)

var (
	ErrCodeNotNull        = "23502"
	ErrCodeForeignKey     = "23503"
	ErrCodeUnique         = "23505"
	ErrCodeRaiseException = "P0001"
)

func ExtractErrorCode(err error) error {
	if err != nil {
		pgxErr, ok := err.(pgx.PgError)
		if !ok {
			return models.ErrInvalid
		}
		return models.NewError(pgxErr.Code)
	}
	return nil
}

func ConvertToPgError(err error) (pgx.PgError, bool) {
	pgError, ok := err.(pgx.PgError)
	return pgError, ok
}
