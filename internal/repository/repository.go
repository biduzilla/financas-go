package repository

import (
	"database/sql"
)

type Repository struct {
	User         UserRepository
	Category     CategoryRepositoryIntercafe
	Transaction  TransactionRepositoryInterface
	Goal         GoalRepositoryInterface
	GoalProgress GoalProgressRepositoryInterface
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Category:     NewCategoryRepository(db),
		Transaction:  NewTransactionRepository(db),
		Goal:         NewGoalRepository(db),
		GoalProgress: NewGoalProgressRepository(db),
	}
}
