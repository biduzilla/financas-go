package service

import (
	"database/sql"
	"financas/internal/repository"
)

type Service struct {
	User UserServiceInterface
}

func NewService(db *sql.DB) *Service {
	repository := repository.NewRepository(db)
	return &Service{
		User: NewUserService(repository.User),
	}
}
