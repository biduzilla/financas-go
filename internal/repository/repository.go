package repository

import (
	"database/sql"
)

type Repository struct {
	User     UserRepository
	Category CategoryRepositoryIntercafe
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:     NewUserRepository(db),
		Category: NewCategoryRepository(db),
	}
}
