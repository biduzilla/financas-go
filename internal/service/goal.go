package service

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/repository"
	e "financas/utils/errors"
	"financas/utils/validator"
	"time"
)

type GoalService struct {
	Goal repository.GoalRepositoryInterface
}

type GoalServiceInterface interface {
	GetAllByUserId(name string, userID int64, f filters.Filters, v *validator.Validator) ([]*model.Goal, filters.Metadata, error)
	GetById(v *validator.Validator, id, userID int64) (*model.Goal, error)
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

	for _, g := range goals {
		if s.verifyFailedStatus(g) {
			g.Status = model.GoalStatusFailed
		}
		s.calculateInstallments(g)
	}

	return goals, meta, nil
}

func (s *GoalService) GetById(v *validator.Validator, id, userID int64) (*model.Goal, error) {
	goal, err := s.Goal.GetById(id, userID)
	if err != nil {
		return nil, err
	}

	s.updateStatus(v, userID, goal)
	s.calculateInstallments(goal)
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

func (s *GoalService) updateStatus(v *validator.Validator, userID int64, goal *model.Goal) error {
	if s.verifyFailedStatus(goal) {
		goal.Status = model.GoalStatusFailed
		return s.Update(v, goal, userID)
	}
	return nil
}

func (s *GoalService) verifyFailedStatus(goal *model.Goal) bool {
	return time.Now().After(goal.Deadline) && goal.Current < goal.Amount && goal.Status != model.GoalStatusFailed
}

func (s *GoalService) calculateInstallments(goal *model.Goal) {
	now := time.Now()
	if goal.Deadline.Before(now) {
		return
	}

	if goal.Installments == nil {
		goal.Installments = &model.Installments{}
	}

	yearDiff := goal.Deadline.Year() - now.Year()
	monthDiff := int(goal.Deadline.Month()) - int(now.Month())
	quantity := yearDiff*12 + monthDiff

	if quantity <= 0 {
		return
	}

	goal.Installments.Quantity = quantity
	remaining := goal.Amount - goal.Current
	if remaining <= 0 {
		return
	}

	goal.Installments.Amount = remaining / float64(quantity)
}
