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

type GoalRepository struct {
	db *sql.DB
}

type GoalRepositoryInterface interface {
	GetAllByUserId(name string, userID int64, f filters.Filters) ([]*model.Goal, filters.Metadata, error)
	GetById(id, idUser int64) (*model.Goal, error)
	Create(goal *model.Goal) error
	Update(goal *model.Goal, idUser int64) error
	Delete(id, idUser int64) error
}

func NewGoalRepository(db *sql.DB) *GoalRepository {
	return &GoalRepository{db: db}
}

func (r *GoalRepository) GetAllByUserId(name string, userID int64, f filters.Filters) ([]*model.Goal, filters.Metadata, error) {
	query := fmt.Sprintf(`
	SELECT 
		count(*) OVER(),
		id, 
		name, 
		description,  
		user_id,
		deadline, 
		amount, 
		current,
		status,
		version,
		created_at,
		deleted,
		u.created_at as u_created_at, 
		u.name as u_name,
		u.phone as u_phone,
		u.email as u_email,
		u.cod as u_cod,
		u.activated as u_activated,
		u.version as u_version
	FROM goals
	inner join users u on (goals.user_id = u.id)
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

	var goals []*model.Goal
	var totalRecords int

	for rows.Next() {
		goal := &model.Goal{
			User: &model.User{},
		}
		err := rows.Scan(
			&totalRecords,
			&goal.ID,
			&goal.Name,
			&goal.Description,
			&goal.User.ID,
			&goal.Deadline,
			&goal.Amount,
			&goal.Current,
			&goal.Status,
			&goal.Version,
			&goal.CreatedAt,
			&goal.Deleted,
			&goal.User.CreatedAt,
			&goal.User.Name,
			&goal.User.Phone,
			&goal.User.Email,
			&goal.User.Cod,
			&goal.User.Activated,
			&goal.User.Version,
		)

		if err != nil {
			return nil, filters.Metadata{}, err
		}

		goals = append(goals, goal)
	}

	if err = rows.Err(); err != nil {
		return nil, filters.Metadata{}, err
	}

	metaData := filters.CalculateMetadata(totalRecords, f.Page, f.PageSize)
	return goals, metaData, nil
}

func (r *GoalRepository) GetById(id, idUser int64) (*model.Goal, error) {
	query := `
	SELECT
		id,
		name,
		description,
		color,
		user_id,
		deadline,
		amount,
		current,
		status,
		version,
		created_at,
		deleted,
		u.created_at as u_created_at, 
		u.name as u_name,
		u.phone as u_phone,
		u.email as u_email,
		u.cod as u_cod,
		u.activated as u_activated,
		u.version as u_version
	from goals
	inner join users u on (goals.user_id = u.id)
	where 
		goals.id = $1 
		and goals.user_id = $2 
		and goals.deleted = false
	`

	goal := &model.Goal{}
	goal.User = &model.User{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, id, idUser).Scan(
		&goal.ID,
		&goal.Name,
		&goal.Description,
		&goal.User.ID,
		&goal.Deadline,
		&goal.Amount,
		&goal.Current,
		&goal.Status,
		&goal.Version,
		&goal.CreatedAt,
		&goal.Deleted,
		&goal.User.CreatedAt,
		&goal.User.Name,
		&goal.User.Phone,
		&goal.User.Email,
		&goal.User.Cod,
		&goal.User.Activated,
		&goal.User.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return goal, nil
}

func (r *GoalRepository) Create(goal *model.Goal) error {
	query := `
	INSERT INTO goals 
		(name, 
		description, 
		color, 
		deadline, 
		amount, 
		current, 
		status, 
		created_at, 
		deleted)
	VALUES ($1,
		$2, 
		$3, 
		$4, 
		$5, 
		$6, 
		$7, 
		1, 
		NOW(), 
		false)
	RETURNING 
		id, 
		created_at, 
		version
	`
	args := []any{
		goal.Name,
		goal.Description,
		goal.Color,
		goal.User.ID,
		goal.Deadline,
		goal.Amount,
		goal.Current,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&goal.ID,
		&goal.CreatedAt,
		&goal.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "unique_user_goal_name":
				return e.ErrDuplicateName
			}
		}
		return err
	}

	return nil
}

func (r *GoalRepository) Update(goal *model.Goal, idUser int64) error {
	query := `
	UPDATE goals
	SET	
		name = $1, 
		description = $2, 
		color = $3,
		deadline = $4, 
		amount = $5, 
		current = $6,
		status = $7,
		version = version + 1
	WHERE 
		id = $8 
		AND version = $9
		AND user_id = $10
		AND deleted = false
	RETURNING version
	`

	args := []any{
		goal.Name,
		goal.Description,
		goal.Color,
		goal.Deadline,
		goal.Amount,
		goal.Current,
		goal.Status,
		goal.ID,
		goal.Version,
		idUser,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&goal.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "unique_user_goal_name":
				return e.ErrDuplicateName
			}
		}

		return err
	}

	return nil
}

func (r *GoalRepository) Delete(id, idUser int64) error {
	query := `
	UPDATE goals
	SET deleted = true
	WHERE 
		id = $1 
		AND user_id = $2 
		AND deleted = false
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id, idUser)

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
