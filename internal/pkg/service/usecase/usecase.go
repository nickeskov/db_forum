package usecase

import (
	serviceModels "github.com/nickeskov/db_forum/internal/pkg/models/service"
	"github.com/nickeskov/db_forum/internal/pkg/service"
)

type UseCase struct {
	repo service.Repository
}

func NewUseCase(repo service.Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}

func (useCase UseCase) DropAllData() error {
	return useCase.repo.DropAllData()
}

func (useCase UseCase) GetStatus() (serviceModels.Status, error) {
	return useCase.repo.GetStatus()
}
