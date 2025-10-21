package model

import (
	"financas/utils/validator"
	"time"
)

type Transaction struct {
	ID          int64
	CreatedAt   time.Time
	Deleted     bool
	Version     int
	User        *User
	Category    *Category
	Description string
	Amount      float64
}

type TransactionDTO struct {
	ID          *int64       `json:"transaction_id"`
	Version     *int         `json:"version"`
	User        *UserDTO     `json:"user"`
	Category    *CategoryDTO `json:"category"`
	Description *string      `json:"description"`
	Amount      *float64     `json:"amount"`
	CreatedAt   *time.Time   `json:"created_at"`
}

func (t *Transaction) ToDTO() *TransactionDTO {
	dto := TransactionDTO{}

	dto.ID = &t.ID
	dto.CreatedAt = &t.CreatedAt
	dto.Version = &t.Version
	dto.Description = &t.Description
	dto.Amount = &t.Amount

	if t.User != nil {
		dto.User = t.User.ToDTO()
	}

	if t.Category != nil {
		dto.Category = t.Category.ToDTO()
	}

	return &dto
}

func (t *TransactionDTO) ToModel() *Transaction {
	transaction := &Transaction{}

	if t.ID != nil {
		transaction.ID = *t.ID
	}
	if t.Version != nil {
		transaction.Version = *t.Version
	}
	if t.User != nil {
		transaction.User = t.User.ToModel()
	}
	if t.Category != nil {
		transaction.Category = t.Category.ToModel()
	}
	if t.Description != nil {
		transaction.Description = *t.Description
	}
	if t.Amount != nil {
		transaction.Amount = *t.Amount
	}

	return transaction
}

func (t *Transaction) ValidateTransaction(v *validator.Validator) {
	v.Check(t.User != nil, "user", "must be provided")
	v.Check(t.Category != nil, "category", "must be provided")
	v.Check(t.Description != "", "description", "must be provided")
	v.Check(len(t.Description) <= 500, "description", "must not be more than 500 bytes long")
	v.Check(t.Amount > 0, "amount", "must be positive")
	v.Check(t.Amount != 0, "amount", "must be provided")
}
