package service

import "github.com/nickeskov/db_forum/internal/pkg/models/service"

type UseCase interface {
	DropAllData() error
	GetStatus() (service.Status, error)
}
