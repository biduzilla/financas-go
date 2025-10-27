package repository

import (
	"context"
	"database/sql"
	"errors"
	"financas/internal/model"
	"financas/internal/model/filters"
	e "financas/utils/errors"
	"fmt"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

type TransactionRepositoryInterface interface {
	GetAllByUserAndCategory(description string, userID int64, categoryID int64, startDate, endDate *time.Time, f filters.Filters) ([]*model.Transaction, filters.Metadata, error)
	GetByID(id int64, userID int64) (*model.Transaction, error)
	Insert(transaction *model.Transaction) error
	Update(transaction *model.Transaction) error
	Delete(id int64, userID int64) error
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) GetAllByUserAndCategory(description string, userID int64, categoryID int64, startDate, endDate *time.Time, f filters.Filters) ([]*model.Transaction, filters.Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), 
	t.id, 
	t.created_at, 
	t.deleted, 
	t.version, 
	t.user_id, 
	t.category_id, 
	t.description, 
	t.amount
	FROM transactions t
	WHERE (to_tsvector('simple', t.description) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND t.user_id = $2 
	AND t.deleted = false 
	AND ($3 = 0 OR t.category_id = $3)
	AND ($4::timestamptz IS NULL OR t.created_at >= $4::timestamptz)
	AND ($5::timestamptz IS NULL OR t.created_at <= $5::timestamptz)
	ORDER BY %s %s, id ASC
	LIMIT $6 OFFSET $7
	`, f.SortColumn(), f.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := sql.NullTime{}
	if startDate != nil {
		start.Valid = true
		start.Time = *startDate
	}

	end := sql.NullTime{}
	if endDate != nil {
		end.Valid = true
		end.Time = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}
	args := []any{
		description,
		userID,
		categoryID,
		start,
		end,
		f.Limit(),
		f.Offset(),
	}
	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	transactions := []*model.Transaction{}

	for rows.Next() {
		transaction := model.Transaction{
			User:     &model.User{},
			Category: &model.Category{User: &model.User{}},
		}
		err := rows.Scan(
			&totalRecords,
			&transaction.ID,
			&transaction.CreatedAt,
			&transaction.Deleted,
			&transaction.Version,
			&transaction.User.ID,
			&transaction.Category.ID,
			&transaction.Description,
			&transaction.Amount,
		)
		if err != nil {
			return nil, filters.Metadata{}, err
		}
		transactions = append(transactions, &transaction)
	}
	metaData := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)
	return transactions, metaData, nil
}

func (r *TransactionRepository) GetByID(id int64, userID int64) (*model.Transaction, error) {
	query := `
	SELECT id, created_at, deleted, version, user_id, category_id, description, amount
	FROM transactions
	WHERE id = $1 AND user_id = $2 AND deleted = false
	`

	var tx model.Transaction
	tx.User = &model.User{}
	tx.Category = &model.Category{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&tx.ID,
		&tx.CreatedAt,
		&tx.Deleted,
		&tx.Version,
		&tx.User.ID,
		&tx.Category.ID,
		&tx.Description,
		&tx.Amount,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tx, nil
}

func (r *TransactionRepository) Insert(transaction *model.Transaction) error {
	query := `
	INSERT INTO transactions ( 
			user_id, 
			category_id, 
			description, 
			amount
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id,created_at, version
	`

	args := []any{
		transaction.User.ID,
		transaction.Category.ID,
		transaction.Description,
		transaction.Amount,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.Version,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) Update(transaction *model.Transaction) error {
	query := `
	UPDATE transactions
	SET user_id = $1, 
		category_id = $2, 
		description = $3, 
		amount = $4, 
		version = version + 1
	WHERE 
		id = $5
		AND user_id = $6
		AND deleted = false 
		AND version = $7
	RETURNING version
	`

	args := []any{
		transaction.User.ID,
		transaction.Category.ID,
		transaction.Description,
		transaction.Amount,
		transaction.ID,
		transaction.User.ID,
		transaction.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&transaction.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return e.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (r *TransactionRepository) Delete(id int64, userID int64) error {
	query := `
	UPDATE transactions
	SET 
		deleted = true
	WHERE 
		id = $1 
		AND user_id = $2 
		AND deleted = false
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return e.ErrRecordNotFound
	}

	return nil
}
