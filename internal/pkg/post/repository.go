package post

import "github.com/nickeskov/db_forum/internal/pkg/models"

type Repository interface {
	CreatePostsInThread(thread models.Thread, posts models.Posts) (models.Posts, error)
	GetPostByID(id int64) (models.Post, error)
	UpdatePostByID(post models.Post) (models.Post, error)
}
