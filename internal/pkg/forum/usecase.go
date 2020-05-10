package forum

import "github.com/nickeskov/db_forum/internal/pkg/models"

type UseCase interface {
	Create(user models.Forum) (models.Forum, error)
	GetBySlug(slug string) (models.Forum, error)
	GetForumUsersBySlug(slug, sinceNickname, desc, limit string) (models.Users, error)
}
