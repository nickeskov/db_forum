package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickeskov/db_forum/internal/pkg/forum"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
	sqlHelpers "github.com/nickeskov/db_forum/pkg/sql"
	pgx4Helpers "github.com/nickeskov/db_forum/pkg/sql/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

type Repository struct {
	db        *pgxpool.Pool
	forumRepo forum.Repository
}

func NewRepository(db *pgxpool.Pool, forumRepo forum.Repository) Repository {
	return Repository{
		db:        db,
		forumRepo: forumRepo,
	}
}

func (repo Repository) GetByID(id int32) (models.Thread, error) {
	return getByID(repo.db, id)
}

func (repo Repository) GetBySlug(slug string) (models.Thread, error) {
	return getBySlug(repo.db, slug)
}

func (repo Repository) VoteByID(id int32, vote models.Vote) (thread models.Thread, err error) {
	ctx := context.Background()

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return models.Thread{}, errors.WithStack(err)
	}
	defer func() {
		err = pgx4Helpers.FinishPgx4Transaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, `
			INSERT INTO votes (thread_id, author_nickname, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, author_nickname) DO UPDATE SET voice = $3`,
		id,
		vote.Nickname,
		vote.Voice,
	)

	if pgxErr := codes.ExtractPgx4ErrorCode(err); pgxErr != nil {
		switch pgxErr.Error() {
		case codes.ErrCodeForeignKey:
			return models.Thread{}, models.ErrDoesNotExist // author or thread does not exist
		default:
			return models.Thread{}, errors.Wrapf(err,
				"some error while voting threadID=%d, vote=%+v", id, vote)
		}
	}

	if thread, err = getByID(tx, id); err != nil {
		return models.Thread{}, errors.Wrapf(err,
			"some error while voting threadID=%d, vote=%+v", id, vote)
	}

	return thread, errors.WithStack(err)
}

func (repo Repository) VoteBySlug(slug string, vote models.Vote) (thread models.Thread, err error) {
	ctx := context.Background()

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return models.Thread{}, errors.WithStack(err)
	}
	defer func() {
		err = pgx4Helpers.FinishPgx4Transaction(ctx, tx, err)
	}()

	thread, err = getBySlug(tx, slug) // get thread model for id
	switch {
	case err == models.ErrDoesNotExist:
		return models.Thread{}, models.ErrDoesNotExist // thread does not exist
	case err != nil:
		return models.Thread{}, errors.Wrapf(err,
			"some error while voting threadSlug=%s, vote=%+v", slug, vote)
	}

	_, err = tx.Exec(ctx, `
			INSERT INTO votes (thread_id, author_nickname, voice)
			VALUES ($1, $2, $3)
			ON CONFLICT (thread_id, author_nickname) DO UPDATE SET voice = $3`,
		thread.ID,
		vote.Nickname,
		vote.Voice,
	)

	if pgxErr := codes.ExtractPgx4ErrorCode(err); pgxErr != nil {
		switch pgxErr.Error() {
		case codes.ErrCodeForeignKey:
			return models.Thread{}, models.ErrDoesNotExist // author does not exist
		default:
			return models.Thread{}, errors.Wrapf(err,
				"some error while voting threadSlug=%s, vote=%+v", slug, vote)
		}
	}

	err = tx.QueryRow(ctx,
		`SELECT votes FROM threads WHERE id = $1`, // update votes in thread
		thread.ID,
	).Scan(
		&thread.Votes,
	)
	if err != nil {
		return models.Thread{}, errors.Wrapf(err,
			"some error while voting threadSlug=%s, vote=%+v", slug, vote)
	}

	return thread, errors.WithStack(err)
}

func (repo Repository) Create(thread models.Thread) (models.Thread, error) {
	ctx := context.Background()

	var threadSlug *string
	if thread.Slug != "" {
		threadSlug = &thread.Slug
	}

	// TODO(nickeskov): maybe remove nickname select, because tests passing without it
	err := repo.db.QueryRow(ctx, `
			INSERT INTO threads (slug, author_nickname, title, message, created, forum_slug)
			VALUES ($1, (SELECT nickname FROM users WHERE nickname = $2), $3, $4, $5,
					(SELECT slug FROM forums WHERE slug = $6))
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

	if pgxErr := codes.ExtractPgx4ErrorCode(err); pgxErr != nil {
		switch pgxErr.Error() {
		case codes.ErrCodeUnique:
			return models.Thread{}, models.ErrConflict
		case codes.ErrCodeForeignKey, codes.ErrCodeNotNull:
			return models.Thread{}, models.ErrBadForeign
		default:
			return models.Thread{}, errors.Wrapf(err,
				"some error while creating thread=%+v", thread)
		}
	}

	return thread, errors.WithStack(err)
}

func (repo Repository) UpdateByID(thread models.Thread) (models.Thread, error) {
	ctx := context.Background()

	row := repo.db.QueryRow(ctx, `
			UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE id = $1
			RETURNING id, slug, forum_slug, author_nickname, title, message, votes, created`,
		thread.ID,
		thread.Title,
		thread.Message,
	)

	err := scanThread(row, &thread)

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
	ctx := context.Background()

	row := repo.db.QueryRow(ctx, `
			UPDATE threads
			SET title   = COALESCE(NULLIF($2, ''), title),
				message = COALESCE(NULLIF($3, ''), message)
			WHERE slug = $1
			RETURNING id, slug, forum_slug, author_nickname, title, message, votes, created`,
		thread.Slug,
		thread.Title,
		thread.Message,
	)

	err := scanThread(row, &thread)

	switch {
	case err == pgx.ErrNoRows:
		return models.Thread{}, models.ErrDoesNotExist
	case err != nil:
		return models.Thread{}, errors.Wrapf(err,
			"some error in thread repo in UpdateBySlug with thread=%+v", thread)
	}

	return thread, nil
}

func (repo Repository) GetThreadsByForumSlug(forumSlug string, since *time.Time, desc bool,
	limit int32) (models.Threads, error) {

	ctx := context.Background()

	var rows pgx.Rows
	var err error

	if since != nil {
		sqlQuery := sqlGetThreadsByForumSlugSince[desc]
		rows, err = repo.db.Query(ctx, sqlQuery, forumSlug, *since, limit)
	} else {
		sqlQuery := sqlGetThreadsByForumSlug[desc]
		rows, err = repo.db.Query(ctx, sqlQuery, forumSlug, limit)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	threads := make(models.Threads, 0)
	for rows.Next() {
		var thread models.Thread

		if err := scanThread(rows, &thread); err != nil {
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

func getByID(querier pgx4Helpers.Querier, id int32) (models.Thread, error) {
	ctx := context.Background()

	var thread models.Thread

	row := querier.QueryRow(ctx, `
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
	)

	err := scanThread(row, &thread)

	if err == pgx.ErrNoRows {
		return models.Thread{}, models.ErrDoesNotExist
	}

	return thread, errors.WithStack(err)
}

func getBySlug(queryer pgx4Helpers.Querier, slug string) (models.Thread, error) {
	ctx := context.Background()

	var thread models.Thread

	row := queryer.QueryRow(ctx, `
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
	)

	err := scanThread(row, &thread)

	if err == pgx.ErrNoRows {
		return models.Thread{}, models.ErrDoesNotExist
	}

	return thread, errors.WithStack(err)
}

func scanThread(scanner sqlHelpers.Scanner, threadDst *models.Thread) error {
	return scanner.Scan(
		&threadDst.ID,
		&threadDst.Slug,
		&threadDst.Forum,
		&threadDst.Author,
		&threadDst.Title,
		&threadDst.Message,
		&threadDst.Votes,
		&threadDst.Created,
	)
}
