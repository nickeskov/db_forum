package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/post"
	"github.com/nickeskov/db_forum/internal/pkg/utils/database/driver/pgx/codes"
	sqlHelpers "github.com/nickeskov/db_forum/pkg/sql"
	pgx4Helpers "github.com/nickeskov/db_forum/pkg/sql/pgx/v4"
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
		err = pgx4Helpers.FinishPgx4Transaction(ctx, tx, err)
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

	var postModel models.Post

	row := repo.db.QueryRow(ctx, `
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
	)

	switch err := scanPosts(row, &postModel); err {
	case nil:
		return postModel, nil

	case pgx.ErrNoRows:
		return models.Post{}, models.ErrDoesNotExist

	default:
		return models.Post{}, errors.WithStack(err)
	}
}

func (repo Repository) UpdatePostByID(post models.Post) (models.Post, error) {
	ctx := context.Background()

	var postMessage *string

	if post.Message != "" {
		postMessage = &post.Message
	}

	row := repo.db.QueryRow(ctx, `
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
	)

	switch err := scanPosts(row, &post); err {
	case nil:
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

	for _, postModel := range posts {
		// TODO(nickeskov): maybe not use parent as nullable column
		var parent *int64
		if postModel.Parent != 0 {
			parent = new(int64)
			*parent = postModel.Parent
		}

		batch.Queue(`
				INSERT INTO posts (thread_id, author_nickname, forum_slug, message, parent, created)
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING id, thread_id, author_nickname, forum_slug, is_edited, message, parent, created`,
			thread.ID,
			postModel.Author,
			thread.Forum,
			postModel.Message,
			parent,
			created,
		)
	}

	return batch
}

func getInsertedPosts(batchResults pgx.BatchResults, batchLen int) (models.Posts, error) {
	insertedPosts := make(models.Posts, batchLen)

	for i := 0; i < batchLen; i++ {
		if err := scanPosts(batchResults.QueryRow(), &insertedPosts[i]); err != nil {
			return nil, err
		}
	}

	return insertedPosts, nil
}

func (repo Repository) GetSortedPostsByThreadSlugOrID(threadID int32, sincePostID *int64,
	sort post.PostsSortType, desc bool, limit int64) (models.Posts, error) { // err = {nil, modesl.ErrInvalid unknown}

	ctx := context.Background()

	var err error
	var rows pgx.Rows

	if sincePostID != nil {
		query, ok := sqlGetSortedPostsSince[desc][sort]
		if !ok {
			return nil, models.ErrInvalid
		}
		rows, err = repo.db.Query(ctx, query, threadID, *sincePostID, limit)
	} else {
		query, ok := sqlGetSortedPosts[desc][sort]
		if !ok {
			return nil, models.ErrInvalid
		}
		rows, err = repo.db.Query(ctx, query, threadID, limit)
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer rows.Close()

	posts := make(models.Posts, 0)
	for rows.Next() {
		var postModel models.Post

		if err := scanPosts(rows, &postModel); err != nil {
			return nil, errors.WithStack(err)
		}

		posts = append(posts, postModel)
	}

	return posts, nil
}

func scanPosts(scanner sqlHelpers.Scanner, postDst *models.Post) error {
	var postParent *int64

	err := scanner.Scan(
		&postDst.ID,
		&postDst.Thread,
		&postDst.Author,
		&postDst.Forum,
		&postDst.IsEdited,
		&postDst.Message,
		&postParent,
		&postDst.Created,
	)
	if err != nil {
		return err
	}

	if postParent != nil {
		postDst.Parent = *postParent
	}

	return nil
}
