package user

import "github.com/nickeskov/db_forum/internal/pkg/models"

type Repository interface {
	Create(user models.User) error
	UpdateByNickname(user models.User) (models.User, error)
	GetByNickname(nickname string) (models.User, error)
	GetWithSameNicknameAndEmail(nickname, email string) (models.Users, error)
}
