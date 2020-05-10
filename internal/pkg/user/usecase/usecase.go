package usecase

import (
	"github.com/nickeskov/db_forum/internal/pkg/models"
	"github.com/nickeskov/db_forum/internal/pkg/user"
)

type UseCase struct {
	repository user.Repository
}

func NewUseCase(repository user.Repository) UseCase {
	return UseCase{
		repository: repository,
	}
}

func (useCase UseCase) Create(user models.User) error {
	return useCase.repository.Create(user)
}

func (useCase UseCase) UpdateByNickname(user models.User) (models.User, error) {
	return useCase.repository.UpdateByNickname(user)
}

func (useCase UseCase) GetByNickname(nickname string) (user models.User, err error) {
	return useCase.repository.GetByNickname(nickname)
}

func (useCase UseCase) GetWithSameNicknameAndEmail(nickname, email string) (users models.Users, err error) {
	return useCase.repository.GetWithSameNicknameAndEmail(nickname, email)
}
