package usecase

import (
	"github.com/nickeskov/db_forum/internal/pkg/models/service"
	"github.com/nickeskov/db_forum/internal/pkg/service/repository"
)

type UseCase struct {
	repo repository.Repository
}

func NewUseCase(repo repository.Repository) UseCase {
	return UseCase{
		repo: repo,
	}
}

func (useCase UseCase) DropAllData() error {
	return useCase.repo.DropAllData()
}

func (useCase UseCase) GetStatus() (service.Status, error) {
	return useCase.repo.GetStatus()
}
