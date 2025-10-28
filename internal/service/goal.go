package service

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/repository"
	e "financas/utils/errors"
	"financas/utils/validator"
)

type GoalService struct {
	Goal repository.GoalRepositoryInterface
}

type GoalServiceInterface interface {
	GetAllByUserId(name string, userID int64, f filters.Filters, v *validator.Validator) ([]*model.Goal, filters.Metadata, error)
	GetById(id, userID int64) (*model.Goal, error)
	Create(v *validator.Validator, goal *model.Goal) error
	Update(v *validator.Validator, goal *model.Goal, userID int64) error
	Delete(id, userID int64) error
}

func NewGoalService(g repository.GoalRepositoryInterface) *GoalService {
	return &GoalService{
		Goal: g,
	}
}

func (s *GoalService) GetAllByUserId(
	name string,
	userID int64,
	f filters.Filters,
	v *validator.Validator,
) ([]*model.Goal, filters.Metadata, error) {
	if filters.ValidateFilters(v, f); !v.Valid() {
		return nil, filters.Metadata{}, e.ErrInvalidData
	}

	goals, meta, err := s.Goal.GetAllByUserId(name, userID, f)
	if err != nil {
		return nil, filters.Metadata{}, err
	}

	return goals, meta, nil
}

func (s *GoalService) GetById(id, userID int64) (*model.Goal, error) {
	goal, err := s.Goal.GetById(id, userID)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) Create(v *validator.Validator, goal *model.Goal) error {
	if goal.ValidateGoal(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.Goal.Create(goal)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalService) Update(v *validator.Validator, goal *model.Goal, userID int64) error {
	if goal.ValidateGoal(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.Goal.Update(goal, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalService) Delete(id, userID int64) error {
	err := s.Goal.Delete(id, userID)
	if err != nil {
		return err
	}
	return nil
}
