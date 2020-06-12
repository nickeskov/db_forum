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

	insertedPosts, err = getInsertedPosts(batchResults, batch.Len())

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

func (repo Repository) GetPostByID(id int64) (models.Post, error) {
	ctx := context.Background()

	var post models.Post
	var postParent *int64

	err := repo.db.QueryRow(ctx, `
			SELECT id,
				   thread_id,
				   author_nickname,
				   forum_slug,
				   is_edited,
				   message,
				   parent,
				   created
			FROM posts
			WHERE id = $1`,
		id,
	).Scan(
		&post.ID,
		&post.Thread,
		&post.Author,
		&post.Forum,
		&post.IsEdited,
		&post.Message,
		&postParent,
		&post.Created,
	)

	if err == pgx.ErrNoRows {
		return models.Post{}, models.ErrDoesNotExist
	}

	if postParent != nil {
		post.Parent = *postParent
	}

	return post, errors.WithStack(err)
}

func (repo Repository) UpdatePostByID(post models.Post) (models.Post, error) {
	ctx := context.Background()

	var postParent *int64
	var postMessage *string

	if post.Message != "" {
		postMessage = &post.Message
	}

	err := repo.db.QueryRow(ctx, `
			UPDATE posts
			SET message   = COALESCE($2, message),
				is_edited = CASE
								WHEN (is_edited = TRUE
									OR (is_edited = FALSE AND $2 IS NOT NULL AND $2 <> message)) THEN TRUE
								ELSE FALSE
					END
			WHERE id = $1
			RETURNING id, thread_id, author_nickname, forum_slug, is_edited, message, parent, created`,
		post.ID,
		postMessage,
	).Scan(
		&post.ID,
		&post.Thread,
		&post.Author,
		&post.Forum,
		&post.IsEdited,
		&post.Message,
		&postParent,
		&post.Created,
	)

	switch err {
	case nil:
		if postParent != nil {
			post.Parent = *postParent
		}
		return post, nil

	case pgx.ErrNoRows:
		return models.Post{}, models.ErrDoesNotExist

	default:
		return models.Post{}, errors.WithStack(err)
	}
}

func createPostsBatch(thread models.Thread, posts models.Posts) *pgx.Batch {
	batch := new(pgx.Batch)

	created := time.Now()

	for _, post := range posts {
		// TODO(nickeskov): maybe not use parent as nullable column
		var parent *int64
		if post.Parent != 0 {
			parent = new(int64)
			*parent = post.Parent
		}

		batch.Queue(`
				INSERT INTO posts (thread_id, author_nickname, forum_slug, message, parent, created)
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING id, thread_id, author_nickname, forum_slug, is_edited, message, parent, created;`,
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
	insertedPosts := make(models.Posts, batchLen)

	for i := 0; i < batchLen; i++ {
		post := &insertedPosts[i]

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
	}

	return insertedPosts, nil
}
