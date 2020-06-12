package thread

import "github.com/nickeskov/db_forum/internal/pkg/models"

type UseCase interface {
	GetBySlugOrID(slugOrID string) (models.Thread, error)
	VoteBySlugOrID(slugOrID string, vote models.Vote) (models.Thread, error)
	Create(thread models.Thread) (models.Thread, error)
	UpdateBySlugOrID(slugOrID string, thread models.Thread) (models.Thread, error)
	GetThreadsByForumSlug(forumSlug, since, desc, limit string) (models.Threads, error)
}
