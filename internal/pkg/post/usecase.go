package post

import "github.com/nickeskov/db_forum/internal/pkg/models"

type UseCase interface {
	CreatePostsByThreadSlugOrID(threadSlugOrID string, posts models.Posts) (models.Posts, error)
	GetPostInfoByID(id int64, related []string) (models.PostFullInfo, error)
	UpdatePostByID(post models.Post) (models.Post, error)
	GetSortedPostsByThreadSlugOrID(threadSlugOrID, sincePostID,
		sort, desc, limit string) (models.Posts, error)
}
