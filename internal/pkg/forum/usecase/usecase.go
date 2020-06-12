package usecase

import (
	"github.com/nickeskov/db_forum/internal/pkg/forum"
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"strconv"
)

type UseCase struct {
	repository forum.Repository
}

func NewUseCase(repository forum.Repository) UseCase {
	return UseCase{
		repository: repository,
	}
}

func (useCase UseCase) Create(user models.Forum) (models.Forum, error) {
	return useCase.repository.Create(user)
}

func (useCase UseCase) GetBySlug(slug string) (models.Forum, error) {
	return useCase.repository.GetBySlug(slug)
}

func (useCase UseCase) GetForumUsersBySlug(slug, sinceNickname, desc, limit string) (models.Users, error) {
	convertedDesc, boolErr := strconv.ParseBool(desc)
	convertedLimit, intErr := strconv.Atoi(limit)
	if boolErr != nil || intErr != nil {
		return nil, models.ErrInvalid
	}

	return useCase.repository.GetForumUsersBySlug(slug, sinceNickname, convertedDesc, int32(convertedLimit))
}
