package repository

import (
	"database/sql"
)

type Repository struct {
	User UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
