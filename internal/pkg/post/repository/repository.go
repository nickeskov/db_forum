package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
	"github.com/nickeskov/db_forum/pkg/sql"
	"github.com/pkg/errors"
	"time"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		db: db,
	}
}

func (repo Repository) CreatePostsInThread(thread models.Thread,
	posts models.Posts) (insertedPosts models.Posts, err error) {

	if len(posts) == 0 {
		return make(models.Posts, 0), nil
	}

	ctx := context.Background()

	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		err = sql.FinishPgx4Transaction(ctx, tx, err)
	}()

	batch := createPostsBatch(thread, posts)

	batchResults := tx.SendBatch(ctx, batch)
	defer func() {
		if closeErr := batchResults.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	insertedPosts, err = getInsertedPosts(batchResults, len(posts))

	if pgxErr := codes.ExtractPgx4ErrorCode(err); pgxErr != nil {
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

	return insertedPosts, errors.WithStack(err)
}

func createPostsBatch(thread models.Thread, posts models.Posts) *pgx.Batch {
	batch := &pgx.Batch{}

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
			thread.ID,
			post.Author,
			thread.Forum,
			post.Message,
			parent,
			created,
		)
	}

	return batch
}

func getInsertedPosts(batchResults pgx.BatchResults, batchLen int) (models.Posts, error) {
	var insertedPosts models.Posts

	for i := 0; i < batchLen; i++ {
		var post models.Post
		var postParent *int64

		err := batchResults.QueryRow().Scan(
			&post.ID,
			&post.Thread,
			&post.Author,
			&post.Forum,
			&post.IsEdited,
			&post.Message,
			&postParent,
			&post.Created,
		)
		if err != nil {
			return nil, err
		}

		if postParent != nil {
			post.Parent = *postParent
		}

		insertedPosts = append(insertedPosts, post)
	}

	return insertedPosts, nil
}
