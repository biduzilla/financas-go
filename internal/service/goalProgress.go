package service

import (
	"financas/internal/model"
	"financas/internal/repository"
	e "financas/utils/errors"
	"financas/utils/validator"
)

type GoalProgressService struct {
	gP repository.GoalProgressRepositoryInterface
}

type GoalProgressServiceInterface interface {
	GetGoalProgressIDGoal(userID, goalID int) ([]*model.GoalProgress, error)
	Insert(v *validator.Validator, gP *model.GoalProgress) error
	Update(v *validator.Validator, gP *model.GoalProgress, userID int64) error
	Delete(goalProgressID, userID int64) error
}

func NewGoalProgressService(gP repository.GoalProgressRepositoryInterface) *GoalProgressService {
	return &GoalProgressService{gP: gP}
}

func (s *GoalProgressService) GetGoalProgressIDGoal(userID, goalID int) ([]*model.GoalProgress, error) {
	gPs, err := s.gP.GetGoalProgressIDGoal(userID, goalID)
	if err != nil {
		return nil, err
	}
	return gPs, nil
}

func (s *GoalProgressService) Insert(v *validator.Validator, gP *model.GoalProgress) error {
	if gP.ValidateGoalProgress(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.gP.Insert(gP)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalProgressService) Update(v *validator.Validator, gP *model.GoalProgress, userID int64) error {
	if gP.ValidateGoalProgress(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.gP.Update(gP, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalProgressService) Delete(goalProgressID, userID int64) error {
	err := s.gP.Delete(goalProgressID, userID)
	if err != nil {
		return err
	}
	return nil
}
