package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
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

func (repo Repository) Create(forum models.Forum) (models.Forum, error) {
	err := repo.db.QueryRow(context.Background(),
		`	INSERT INTO forums (slug, title, threads, posts, owner_nickname)
				VALUES ($1, $2, $3, $4, (
							SELECT nickname FROM users WHERE nickname = $5
				)) RETURNING owner_nickname`,
		forum.Slug,
		forum.Title,
		forum.Threads,
		forum.Posts,
		forum.User,
	).Scan(
		&forum.User,
	)

	if pgxErr := codes.ExtractPgx4ErrorCode(err); pgxErr != nil {
		switch pgxErr.Error() {
		case codes.ErrCodeUnique:
			return models.Forum{}, models.ErrConflict
		case codes.ErrCodeNotNull:
			return models.Forum{}, models.ErrBadForeign
		}
	}

	return forum, errors.WithStack(err)
}

func (repo Repository) GetBySlug(slug string) (models.Forum, error) {
	var forum models.Forum
	err := repo.db.QueryRow(context.Background(),
		`	SELECT slug, title, threads, posts, owner_nickname
				FROM forums
				WHERE slug = $1`,
		slug,
	).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.Threads,
		&forum.Posts,
		&forum.User,
	)

	switch {
	case err == pgx.ErrNoRows:
		return models.Forum{}, models.ErrDoesNotExist

	case err != nil:
		return models.Forum{}, errors.WithStack(err)
	}

	return forum, nil
}

func (repo Repository) GetForumUsersBySlug(slug, sinceNickname string, desc bool, limit int32) (models.Users, error) {
	var err error
	var rows pgx.Rows

	if sinceNickname != "" {
		rows, err = repo.db.Query(context.Background(),
			sqlGetForumUserWithSince[desc],
			slug,
			sinceNickname,
			limit,
		)
	} else {
		rows, err = repo.db.Query(context.Background(),
			sqlGetForumUser[desc],
			slug,
			limit,
		)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "error in forum repository"+
			"GetForumUsersBySlug with slug=%s: %v", slug, err)
	}

	defer rows.Close()

	users := make(models.Users, 0)
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Nickname,
			&user.Email,
			&user.Fullname,
			&user.About,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error in forum repository"+
				"GetForumUsersBySlug with slug=%s while scanning rows: %v", slug, err)
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		var isExists bool
		err := repo.db.QueryRow(context.Background(),
			`	SELECT EXISTS(
               			SELECT 1
               			FROM forums
               			WHERE slug = $1
           			)`,
			slug,
		).Scan(
			&isExists,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error in forum repository"+
				"GetForumUsersBySlug while checking exists slug=%s: %v", slug, err)
		}

		if !isExists {
			return nil, models.ErrDoesNotExist
		}
	}

	return users, nil
}
