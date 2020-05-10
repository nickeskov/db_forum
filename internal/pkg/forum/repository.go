package forum

import "github.com/nickeskov/db_forum/internal/pkg/models"

type Repository interface {
	Create(forum models.Forum) (models.Forum, error)
	GetBySlug(slug string) (models.Forum, error)
	GetForumUsersBySlug(slug, sinceNickname string, desc bool, limit int32) (models.Users, error)
}
