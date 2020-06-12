package post

import "github.com/nickeskov/db_forum/internal/pkg/models"

type UseCase interface {
	CreatePostsByThreadSlugOrID(threadSlugOrID string, posts models.Posts) (models.Posts, error)
}
