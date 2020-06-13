package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickeskov/db_forum/internal/pkg/models/service"
	pgx4Helpers "github.com/nickeskov/db_forum/pkg/sql/pgx/v4"
	"github.com/pkg/errors"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		db: db,
	}
}

func (repo Repository) DropAllData() error {
	ctx := context.Background()

	_, err := repo.db.Exec(ctx,
		`TRUNCATE users, forums, threads, votes, posts, forums_users_nicknames`)
	return errors.WithStack(err)
}

func (repo Repository) GetStatus() (status service.Status, err error) {
	ctx := context.Background()

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	defer func() {
		err = pgx4Helpers.FinishPgx4Transaction(ctx, tx, err)
	}()

	if err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&status.User); err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	if err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM forums`).Scan(&status.Forum); err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	if err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM threads`).Scan(&status.Thread); err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	if err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM posts`).Scan(&status.Post); err != nil {
		return service.Status{}, errors.WithStack(err)
	}

	return status, nil
}
