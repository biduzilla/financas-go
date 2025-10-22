package service

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/repository"
	e "financas/utils/errors"
	"financas/utils/validator"
	"time"
)

type TransactionService struct {
	Transaction repository.TransactionRepositoryInterface
}

func NewTransactionService(r repository.TransactionRepositoryInterface) *TransactionService {
	return &TransactionService{
		Transaction: r,
	}
}

type TransactionServiceInterface interface {
	GetByID(id, userID int64) (*model.Transaction, error)
	GetAllByUserAndCategory(v *validator.Validator, description string, userID int64, categoryID int64, startDate, endDate *time.Time, f filters.Filters) ([]*model.Transaction, filters.Metadata, error)
	Save(v *validator.Validator, t *model.Transaction) error
	Update(v *validator.Validator, t *model.Transaction, userID int64) error
	Delete(id, userID int64) error
}

func (s *TransactionService) GetByID(id, userID int64) (*model.Transaction, error) {
	return s.Transaction.GetByID(id, userID)
}

func (s *TransactionService) GetAllByUserAndCategory(v *validator.Validator, description string, userID int64, categoryID int64, startDate, endDate *time.Time, f filters.Filters) ([]*model.Transaction, filters.Metadata, error) {
	if filters.ValidateFilters(v, f); !v.Valid() {
		return nil, filters.Metadata{}, e.ErrInvalidData
	}

	t, m, err := s.Transaction.GetAllByUserAndCategory(description, userID, categoryID, startDate, endDate, f)
	if err != nil {
		return nil, filters.Metadata{}, err
	}

	return t, m, nil
}

func (s *TransactionService) Save(v *validator.Validator, t *model.Transaction) error {
	if t.ValidateTransaction(v); !v.Valid() {
		return e.ErrInvalidData
	}

	return s.Transaction.Insert(t)
}

func (s *TransactionService) Update(v *validator.Validator, t *model.Transaction, userID int64) error {
	if userID != t.User.ID {
		return e.ErrRecordNotFound
	}

	if t.ValidateTransaction(v); !v.Valid() {
		return e.ErrInvalidData
	}

	return s.Transaction.Update(t)
}

func (s *TransactionService) Delete(id, userID int64) error {
	return s.Transaction.Delete(id, userID)
}
