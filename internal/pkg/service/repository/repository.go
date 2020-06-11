package repository

import (
	"github.com/jackc/pgx"
	"github.com/nickeskov/db_forum/internal/pkg/models/service"
	"github.com/nickeskov/db_forum/pkg/sql"
	"github.com/pkg/errors"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewRepository(db *pgx.ConnPool) Repository {
	return Repository{
		db: db,
	}
}

func (repo Repository) DropAllData() error {
	_, err := repo.db.Exec(`TRUNCATE users, forums, threads, votes, posts, forums_users_nicknames`)
	return errors.WithStack(err)
}

func (repo Repository) GetStatus() (status service.Status, err error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	defer func() {
		err = sql.FinishTransaction(tx, err)
	}()

	if err = tx.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&status.User); err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	if err = tx.QueryRow(`SELECT COUNT(*) FROM forums`).Scan(&status.Forum); err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	if err = tx.QueryRow(`SELECT COUNT(*) FROM threads`).Scan(&status.Thread); err != nil {
		return service.Status{}, errors.WithStack(err)
	}
	if err = tx.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&status.Post); err != nil {
		return service.Status{}, errors.WithStack(err)
	}

	return status, nil
}
