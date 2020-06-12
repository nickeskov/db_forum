package post

import "github.com/nickeskov/db_forum/internal/pkg/models"

type Repository interface {
	CreatePostsInThread(thread models.Thread, posts models.Posts) (models.Posts, error)
}
