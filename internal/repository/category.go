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

	"github.com/lib/pq"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) GetByID(id int64, userID int64) (*model.Category, error) {
	query := `
	SELECT id, created_at, name, type, color, user_id,version
	FROM categories
	WHERE id = $1 AND user_id = $2 AND deleted = false
	`

	category := model.Category{
		User: &model.User{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Name,
		&category.Type,
		&category.Color,
		&category.User.ID,
		&category.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (r *CategoryRepository) GetAll(name string, userID int64, f filters.Filters) ([]*model.Category, filters.Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, name, type, color, user_id,version
	FROM categories
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND user_id = $2 AND deleted = false
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4
	`, f.SortColumn(), f.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, userID, f.Limit(), f.Offset()}

	rows, err := r.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	categories := []*model.Category{}

	for rows.Next() {
		category := model.Category{
			User: &model.User{},
		}

		err := rows.Scan(
			&totalRecords,
			&category.ID,
			&category.CreatedAt,
			&category.Name,
			&category.Type,
			&category.Color,
			&category.User.ID,
			&category.Version,
		)

		if err != nil {
			return nil, filters.Metadata{}, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metaData := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)
	return categories, metaData, nil
}

func (r *CategoryRepository) Insert(category *model.Category) error {
	query := `
	INSERT INTO categories (name, type, color, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{
		category.Name,
		category.Type,
		category.Color,
		category.User.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "categories_nome_key":
				return e.ErrDuplicateName
			}
		}

		return err
	}

	return nil
}

func (r *CategoryRepository) Update(category *model.Category, userID int64) error {
	query := `
	UPDATE categories
	SET 
		name = $1, 
		type = $2, 
		color = $3, 
		version = version + 1
	WHERE 
		id = $4 
		AND user_id = $5 
		AND deleted = false 
		AND version = $6
	RETURNING version
	`

	args := []any{
		category.Name,
		category.Type,
		category.Color,
		category.ID,
		userID,
		category.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&category.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "categories_nome_key":
				return e.ErrDuplicateName
			}
		}

		return err
	}

	return nil
}

func (r *CategoryRepository) Delete(id int64, userID int64) error {
	query := `
	UPDATE from categories
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
