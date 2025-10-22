package service

import (
	"database/sql"
	"financas/internal/config"
	"financas/internal/repository"
)

type Service struct {
	User        UserServiceInterface
	Auth        AuthServiceInterface
	Category    CategoryServiceInterface
	Transaction TransactionServiceInterface
}

func NewService(db *sql.DB, config config.Config) *Service {
	repository := repository.NewRepository(db)
	userService := NewUserService(repository.User)

	return &Service{
		User:        userService,
		Auth:        NewAuthService(userService, config),
		Category:    NewCategoryService(repository.Category),
		Transaction: NewTransactionService(repository.Transaction),
	}
}
