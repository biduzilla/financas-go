package service

import (
	"financas/internal/model"
	"financas/internal/repository"
	e "financas/utils/errors"
	"financas/utils/validator"
	"time"
)

type GoalProgressService struct {
	gP   repository.GoalProgressRepositoryInterface
	goal GoalServiceInterface
}

type GoalProgressServiceInterface interface {
	GetGoalProgressIDGoal(userID, goalID int64) ([]*model.GoalProgress, error)
	Insert(v *validator.Validator, gP *model.GoalProgress, userID int64) error
	Update(v *validator.Validator, gP *model.GoalProgress, userID int64) error
	Delete(v *validator.Validator, goalProgressID, userID int64) error
}

func NewGoalProgressService(gP repository.GoalProgressRepositoryInterface, goal GoalServiceInterface) *GoalProgressService {
	return &GoalProgressService{gP: gP, goal: goal}
}

func (s *GoalProgressService) GetGoalProgressIDGoal(userID, goalID int64) ([]*model.GoalProgress, error) {
	gPs, err := s.gP.GetGoalProgressIDGoal(userID, goalID)
	if err != nil {
		return nil, err
	}
	return gPs, nil
}

func (s *GoalProgressService) GetGoalProgressByID(userID, gPID int64) (*model.GoalProgress, error) {
	gP, err := s.gP.GetGoalProgressByID(userID, gPID)
	if err != nil {
		return nil, err
	}
	return gP, nil
}

func (s *GoalProgressService) Insert(v *validator.Validator, gP *model.GoalProgress, userID int64) error {
	if gP.ValidateGoalProgress(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.gP.Insert(gP)
	if err != nil {
		return err
	}

	goal, err := s.goal.GetById(gP.Goal.ID, userID)

	if err != nil {
		return err
	}

	goal.Current += gP.Amount

	if goal.Current >= goal.Amount {
		goal.Status = model.GoalStatusFinished
	}

	err = s.goal.Update(v, goal, userID)
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

	return s.updateGoalStatus(v, gP.ID, userID)
}

func (s *GoalProgressService) Delete(v *validator.Validator, goalProgressID, userID int64) error {
	err := s.gP.Delete(goalProgressID, userID)
	if err != nil {
		return err
	}
	return s.updateGoalStatus(v, goalProgressID, userID)
}

func (s *GoalProgressService) updateGoalStatus(v *validator.Validator, gPID, userID int64) error {
	gP, err := s.GetGoalProgressByID(userID, gPID)
	if err != nil {
		return err
	}

	gPs, err := s.gP.GetGoalProgressIDGoal(userID, gP.Goal.ID)
	if err != nil {
		return err
	}

	var totalAmount float64 = 0
	for _, gp := range gPs {
		totalAmount += gp.Amount
	}

	goal, err := s.goal.GetById(gPID, userID)

	if err != nil {
		return err
	}

	goal.Current = totalAmount

	if time.Now().After(goal.Deadline) && goal.Current < goal.Amount {
		goal.Status = model.GoalStatusFailed
	} else if goal.Current >= goal.Amount && goal.Status != model.GoalStatusFailed {
		goal.Status = model.GoalStatusFinished
	} else if goal.Current < 0 {
		goal.Status = model.GoalStatusInProgress
	}

	err = s.goal.Update(v, goal, userID)
	if err != nil {
		return err
	}

	return nil
}
