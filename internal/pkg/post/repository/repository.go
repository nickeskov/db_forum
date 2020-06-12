package repository

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
	"github.com/nickeskov/db_forum/pkg/sql"
	"github.com/pkg/errors"
	"time"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewRepository(db *pgx.ConnPool) Repository {
	return Repository{
		db: db,
	}
}

func (repo Repository) CreatePostsInThread(thread models.Thread,
	posts models.Posts) (insertedPosts models.Posts, err error) {

	if len(posts) == 0 {
		return make(models.Posts, 0), nil
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		err = sql.FinishTransaction(tx, err)
	}()

	batch := tx.BeginBatch()
	defer func() {
		if closeErr := batch.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	err = sendPostsBatch(batch, thread, posts)

	if pgxErr := codes.ExtractErrorCode(err); pgxErr != nil {
		switch pgxErr.Error() {
		case codes.ErrCodeForeignKey:
			return nil, models.ErrDoesNotExist // author or thread does not exist
		case codes.ErrCodeRaiseException:
			return nil, models.ErrConflict // parent post does not exist in this thread
		default:
			return nil, errors.Wrapf(err,
				"some error post creating with thread=%+v, posts=%+v", thread, posts)
		}
	}

	insertedPosts, err = getInsertedPosts(batch, len(posts))

	return insertedPosts, errors.WithStack(err)
}

func sendPostsBatch(batch *pgx.Batch, thread models.Thread, posts models.Posts) error {
	created := time.Now()

	for _, post := range posts {
		var parent *int64
		if post.Parent != 0 {
			parent = &post.Parent
		}

		batch.Queue(`
			INSERT INTO posts (thread_id, author_nickname, forum_slug, message, parent, created)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, thread_id, author_nickname, forum_slug, is_edited, message, parent, created`,
			[]interface{}{
				thread.ID,
				post.Author,
				thread.Forum,
				post.Message,
				parent,
				created,
			}, nil, nil,
		)
	}

	err := batch.Send(context.Background(), &pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})

	return err
}

func getInsertedPosts(batch *pgx.Batch, batchLen int) (models.Posts, error) {
	var insertedPosts models.Posts

	for i := 0; i < batchLen; i++ {
		// TODO(nickeskov): this not work
		row := batch.QueryRowResults()

		var post models.Post

		err := row.Scan(
			&post.ID,
			&post.Thread,
			&post.Author,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&post.Parent,
			&post.Created,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		insertedPosts = append(insertedPosts, post)
	}

	return insertedPosts, nil
}
