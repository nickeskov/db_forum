package usecase

import (
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/thread/repository"
	"github.com/nickeskov/db_forum/internal/pkg/utils"
	"strconv"
	"time"
)

type UseCase struct {
	repo repository.Repository
}

func NewUseCase(repo repository.Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}

func (useCase UseCase) GetBySlugOrID(slugOrID string) (models.Thread, error) {
	if id, err := strconv.Atoi(slugOrID); err != nil {
		return useCase.repo.GetBySlug(slugOrID)
	} else {
		return useCase.repo.GetByID(int32(id))
	}
}

func (useCase UseCase) Create(thread models.Thread) (models.Thread, error) {
	return useCase.repo.Create(thread)
}

func (useCase UseCase) GetThreadsByForumSlug(forumSlug, since, desc, limit string) (models.Threads, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, models.ErrInvalid
	}

	descBool, err := strconv.ParseBool(desc)
	if err != nil {
		return nil, models.ErrInvalid
	}

	if since != "" {
		if sinceTime, err := time.Parse(utils.TimestampFormat, since); err != nil {
			return nil, models.ErrInvalid
		} else {
			return useCase.repo.GetThreadsByForumSlug(forumSlug, &sinceTime, descBool, int32(limitInt))
		}
	}

	return useCase.repo.GetThreadsByForumSlug(forumSlug, nil, descBool, int32(limitInt))
}