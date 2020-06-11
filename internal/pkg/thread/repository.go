package thread

import (
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"time"
)

type Repository interface {
	GetByID(id int32) (models.Thread, error)
	GetBySlug(slug string) (models.Thread, error)

	Create(thread models.Thread) (models.Thread, error)
	GetThreadsByForumSlug(forumSlug string, since *time.Time, desc bool, limit int32) (models.Threads, error)
}
