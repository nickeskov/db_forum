package repository

import (
	"github.com/jackc/pgx"
	"github.com/nickeskov/db_forum/internal/pkg/forum"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
	"github.com/pkg/errors"
	"time"
)

type Repository struct {
	db        *pgx.ConnPool
	forumRepo forum.Repository
}

func NewRepository(db *pgx.ConnPool, forumRepo forum.Repository) Repository {
	return Repository{
		db:        db,
		forumRepo: forumRepo,
	}
}

func (repo Repository) GetByID(id int32) (models.Thread, error) {
	var thread models.Thread

	err := repo.db.QueryRow(`
			SELECT id,
				   slug,
				   forum_slug,
				   author_nickname,
				   title,
				   message,
				   votes,
				   created
			FROM threads
			WHERE id = $1`,
		id,
	).Scan(
		&thread.ID,
		&thread.Slug,
		&thread.Forum,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Votes,
		&thread.Created,
	)

	if err == pgx.ErrNoRows {
		return models.Thread{}, models.ErrDoesNotExist
	}

	return thread, errors.WithStack(err)
}

func (repo Repository) GetBySlug(slug string) (models.Thread, error) {
	var thread models.Thread

	err := repo.db.QueryRow(`
			SELECT id,
				   slug,
				   forum_slug,
				   author_nickname,
				   title,
				   message,
				   votes,
				   created
			FROM threads
			WHERE slug = $1`,
		slug,
	).Scan(
		&thread.ID,
		&thread.Slug,
		&thread.Forum,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Votes,
		&thread.Created,
	)

	if err == pgx.ErrNoRows {
		return models.Thread{}, models.ErrDoesNotExist
	}

	return thread, errors.WithStack(err)
}

func (repo Repository) Create(thread models.Thread) (models.Thread, error) {
	var threadSlug *string
	if thread.Slug != "" {
		threadSlug = &thread.Slug
	}

	err := repo.db.QueryRow(`
			INSERT INTO threads (slug, author_nickname, title, message, created, forum_slug)
			VALUES ($1, $2, $3, $4, $5, (SELECT slug FROM forums WHERE slug = $6)) 
			RETURNING id, author_nickname, forum_slug, created`,
		threadSlug,
		thread.Author,
		thread.Title,
		thread.Message,
		thread.Created,
		thread.Forum,
	).Scan(
		&thread.ID,
		&thread.Author,
		&thread.Forum,
		&thread.Created,
	)

	if pgxErr := codes.ExtractErrorCode(err); pgxErr != nil {
		switch pgxErr.Error() {
		case codes.ErrCodeUnique:
			return models.Thread{}, models.ErrConflict
		case codes.ErrCodeNotNull:
			fallthrough
		case codes.ErrCodeForeignKey:
			return models.Thread{}, models.ErrBadForeign
		default:
			return models.Thread{}, errors.Wrapf(err,
				"some error while creating thread=%+v", thread)
		}
	}

	return thread, errors.WithStack(err)
}

func (repo Repository) UpdateByID(thread models.Thread) (models.Thread, error) {
	err := repo.db.QueryRow(`
			UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE id = $1
			RETURNING id, slug, forum_slug, author_nickname, title, message, votes, created`,
		thread.ID,
		thread.Title,
		thread.Message,
	).Scan(
		&thread.ID,
		&thread.Slug,
		&thread.Forum,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Votes,
		&thread.Created,
	)

	switch {
	case err == pgx.ErrNoRows:
		return models.Thread{}, models.ErrDoesNotExist
	case err != nil:
		return models.Thread{}, errors.Wrapf(err,
			"some error in thread repo in UpdateByID with thread=%+v", thread)
	}

	return thread, nil
}

func (repo Repository) UpdateBySlug(thread models.Thread) (models.Thread, error) {
	err := repo.db.QueryRow(`
			UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE slug = $1
			RETURNING id, slug, forum_slug, author_nickname, title, message, votes, created`,
		thread.Slug,
		thread.Title,
		thread.Message,
	).Scan(
		&thread.ID,
		&thread.Slug,
		&thread.Forum,
		&thread.Author,
		&thread.Title,
		&thread.Message,
		&thread.Votes,
		&thread.Created,
	)

	switch {
	case err == pgx.ErrNoRows:
		return models.Thread{}, models.ErrDoesNotExist
	case err != nil:
		return models.Thread{}, errors.Wrapf(err,
			"some error in thread repo in UpdateBySlug with thread=%+v", thread)
	}

	return thread, nil
}

func (repo Repository) GetThreadsByForumSlug(forumSlug string, since *time.Time, desc bool, limit int32) (models.Threads, error) {
	var rows *pgx.Rows
	var err error

	if since != nil {
		sqlQuery := sqlGetThreadsByForumSlugSince[desc]
		rows, err = repo.db.Query(sqlQuery, forumSlug, *since, limit)
	} else {
		sqlQuery := sqlGetThreadsByForumSlug[desc]
		rows, err = repo.db.Query(sqlQuery, forumSlug, limit)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	threads := make(models.Threads, 0)
	for rows.Next() {
		var thread models.Thread

		err = rows.Scan(
			&thread.ID,
			&thread.Slug,
			&thread.Forum,
			&thread.Author,
			&thread.Title,
			&thread.Message,
			&thread.Votes,
			&thread.Created,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		threads = append(threads, thread)
	}

	if len(threads) == 0 {
		_, err := repo.forumRepo.GetBySlug(forumSlug)
		switch {
		case errors.Is(err, models.ErrDoesNotExist):
			return nil, models.ErrDoesNotExist

		case err != nil:
			return nil, errors.WithStack(err)
		}
	}

	return threads, nil
}
