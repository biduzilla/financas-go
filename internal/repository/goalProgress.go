package repository

import (
	"context"
	"database/sql"
	"errors"
	"financas/internal/model"
	e "financas/utils/errors"
	"time"
)

type GoalProgressRepository struct {
	db *sql.DB
}

type GoalProgressRepositoryInterface interface {
	GetGoalProgressIDGoal(userID, goalID int) ([]*model.GoalProgress, error)
	Insert(gP *model.GoalProgress) error
	Update(gP *model.GoalProgress, userID int64) error
	Delete(goalProgressID, userID int64) error
}

func NewGoalProgressRepository(db *sql.DB) *GoalProgressRepository {
	return &GoalProgressRepository{db: db}
}

func (r *GoalProgressRepository) GetGoalProgressIDGoal(
	userID, goalID int,
) ([]*model.GoalProgress, error) {
	query := `
	SELECT
		gp.id,
		gp.current,
		gp.date,
		gp.version,
		gp.created_at,
		gp.deleted,
		
		g.id AS g_id, 
		g.name AS g_name, 
		g.description AS g_description,  
		g.user_id AS g_user_id,
		g.deadline AS g_deadline, 
		g.amount AS g_amount, 
		g.current AS g_current,
		g.status AS g_status,
		g.version AS g_version,
		g.created_at AS g_created_at,
		g.deleted AS g_deleted,

		u.id AS u_id,
		u.name AS u_name,
		u.phone AS u_phone,
		u.email AS u_email,
		u.cod AS u_cod,
		u.activated AS u_activated,
		u.version AS u_version,
		u.created_at AS u_created_at
		
	FROM
		goal_progress gp
	INNER JOIN goals g ON gp.goal_id = g.id
	INNER JOIN users u ON g.user_id = u.id
	WHERE
		g.user_id = $1
		AND gp.deleted = FALSE
		AND gp.id = $2
	ORDER BY gp.created_at DESC;
`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userID, goalID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var gPs []*model.GoalProgress
	for rows.Next() {
		gP := &model.GoalProgress{
			Goal: &model.Goal{
				User: &model.User{},
			},
		}

		err := rows.Scan(
			&gP.ID,
			&gP.Current,
			&gP.Date,
			&gP.Version,
			&gP.Deleted,
			&gP.Goal.ID,
			&gP.Goal.Name,
			&gP.Goal.Description,
			&gP.Goal.User.ID,
			&gP.Goal.Deadline,
			&gP.Goal.Amount,
			&gP.Goal.Current,
			&gP.Goal.Status,
			&gP.Goal.Version,
			&gP.Goal.CreatedAt,
			&gP.Goal.Deleted,
			&gP.Goal.User.CreatedAt,
			&gP.Goal.User.Name,
			&gP.Goal.User.Phone,
			&gP.Goal.User.Email,
			&gP.Goal.User.Cod,
			&gP.Goal.User.Activated,
			&gP.Goal.User.Version,
			&gP.Goal.User.ID,
		)

		if err != nil {
			return nil, err
		}

		gPs = append(gPs, gP)
	}

	return gPs, nil
}

func (r *GoalProgressRepository) Insert(gP *model.GoalProgress) error {
	query := `
	INSERT INTO goal_progress (
		current,
		date,
		created_at,
		deleted,
		goal_id
	) VALUES (
		$1, $2, $3, false, $5
	) RETURNING 
	 	id,
		created_at,
		version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		gP.Current,
		gP.Date,
		gP.Version,
		gP.CreatedAt,
		gP.Deleted,
		gP.Goal.ID,
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&gP.ID,
		&gP.CreatedAt,
		&gP.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}
		return err
	}

	return nil
}

func (r *GoalProgressRepository) Update(gP *model.GoalProgress, userID int64) error {
	query := `
	UPDATE goal_progress
	SET
		current = $1,
		date = $2,
		version = version + 1
	inner join goals g on goal_progress.goal_id = g.id
	WHERE
		id = $3
		and goal_id = $4
		AND deleted = false
		AND g.user_id = $5
	RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		gP.Current,
		gP.Date,
		gP.ID,
		gP.Goal.ID,
		userID,
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&gP.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrEditConflict
		}
		return err
	}

	return nil
}

func (r *GoalProgressRepository) Delete(goalProgressID, userID int64) error {
	query := `
	UPDATE goal_progress
	SET
		deleted = true,
	WHERE
		id = $1
		AND goal_id IN (
			SELECT id FROM goals WHERE user_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, goalProgressID, userID)

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
