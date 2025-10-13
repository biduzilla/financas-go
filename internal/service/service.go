package service

import (
	"database/sql"
	"financas/configuration"
	"financas/internal/repository"
)

type Service struct {
	User UserServiceInterface
	Auth AuthServiceInterface
}

func NewService(db *sql.DB, config *configuration.Conf) *Service {
	repository := repository.NewRepository(db)
	userService := NewUserService(repository.User)
	return &Service{
		User: userService,
		Auth: NewAuthService(userService, config),
	}
}
