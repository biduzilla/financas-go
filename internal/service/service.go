package service

import (
	"database/sql"
	"financas/internal/repository"
	"financas/utils/validator"
)

type Service struct {
	User UserServiceInterface
}

func NewService(db *sql.DB, v *validator.Validator) *Service {
	repository := repository.NewRepository(db)
	return &Service{
		User: NewUserService(repository.User, v),
	}
}
