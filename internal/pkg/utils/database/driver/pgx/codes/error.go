package codes

import (
	"github.com/jackc/pgx"
	"github.com/nickeskov/db_forum/internal/pkg/models"
)

var (
	ErrCodeForeignKey = "23503"
	ErrCodeUnique     = "23505"
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
