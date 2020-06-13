package post

import "github.com/nickeskov/db_forum/internal/pkg/models"

type Repository interface {
	CreatePostsInThread(thread models.Thread, posts models.Posts) (models.Posts, error)
	GetPostByID(id int64) (models.Post, error)
	UpdatePostByID(post models.Post) (models.Post, error)
	GetSortedPostsByThreadSlugOrID(threadID int32, sincePostID *int64,
		sort PostsSortType, desc bool, limit int64) (models.Posts, error)
}

type PostsSortType string

const (
	FlatSort       PostsSortType = "flat"
	TreeSort       PostsSortType = "tree"
	ParentTreeSort PostsSortType = "parent_tree"
)
