package service

import (
	"database/sql"
	"financas/internal/config"
	"financas/internal/repository"
)

type Service struct {
	User         UserServiceInterface
	Auth         AuthServiceInterface
	Category     CategoryServiceInterface
	Transaction  TransactionServiceInterface
	Report       ReportServiceInterface
	Goal         GoalServiceInterface
	GoalProgress GoalProgressServiceInterface
}

func NewService(db *sql.DB, config config.Config) *Service {
	repository := repository.NewRepository(db)
	userService := NewUserService(repository.User)
	categoryService := NewCategoryService(repository.Category, db)
	transactionService := NewTransactionService(repository.Transaction)
	goalService := NewGoalService(repository.Goal)

	return &Service{
		User:         userService,
		Auth:         NewAuthService(userService, config),
		Category:     categoryService,
		Transaction:  transactionService,
		Report:       NewReportService(transactionService, categoryService),
		Goal:         goalService,
		GoalProgress: NewGoalProgressService(repository.GoalProgress, goalService),
	}
}
