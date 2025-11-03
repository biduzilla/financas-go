package model

import (
	"financas/utils/validator"
	"time"
)

type GoalStatus int

const (
	GoalStatusPending GoalStatus = iota + 1
	GoalStatusInProgress
	GoalStatusFinished
	GoalStatusFailed
)

type Goal struct {
	ID           int64
	Name         string
	Description  string
	Color        string
	User         *User
	Deadline     time.Time
	Amount       float64
	Current      float64
	Status       GoalStatus
	Version      int
	CreatedAt    time.Time
	Deleted      bool
	Installments *Installments
}
type GoalDTO struct {
	ID           *int64        `json:"goal_id"`
	Name         *string       `json:"name"`
	Description  *string       `json:"description"`
	Color        *string       `json:"color"`
	User         *UserDTO      `json:"user"`
	Deadline     *string       `json:"deadline"`
	Amount       *float64      `json:"amount"`
	Current      *float64      `json:"current"`
	Status       *string       `json:"status"`
	Version      *int          `json:"version"`
	CreatedAt    *time.Time    `json:"created_at"`
	Installments *Installments `json:"Installments"`
}

type Installments struct {
	Amount   float64
	Quantity int
}

type GoalProgress struct {
	ID        int64
	Goal      *Goal
	Amount    float64
	Date      time.Time
	Version   *int
	CreatedAt *time.Time
	Deleted   bool
}

type GoalProgressDTO struct {
	ID        *int64     `json:"goal_progress_id"`
	Goal      *GoalDTO   `json:"goal"`
	Amount    *float64   `json:"amount"`
	Date      *time.Time `json:"date"`
	Version   *int       `json:"version"`
	CreatedAt *time.Time `json:"created_at"`
	Deleted   *bool      `json:"deleted"`
}

func (s GoalStatus) String() string {
	switch s {
	case GoalStatusPending:
		return "PENDING"
	case GoalStatusInProgress:
		return "IN_PROGRESS"
	case GoalStatusFinished:
		return "FINISHED"
	case GoalStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

func GoalStatusFromString(str string) GoalStatus {
	switch str {
	case "PENDING":
		return GoalStatusPending
	case "IN_PROGRESS":
		return GoalStatusInProgress
	case "FINISHED":
		return GoalStatusFinished
	case "FAILED":
		return GoalStatusFailed
	default:
		return 0
	}
}

func (g *Goal) ToDTO() *GoalDTO {
	goal := &GoalDTO{}
	goal.ID = &g.ID
	goal.Name = &g.Name
	goal.Description = &g.Description
	goal.Color = &g.Color
	goal.User = g.User.ToDTO()
	deadlineStr := g.Deadline.Format("02/01/2006")
	goal.Deadline = &deadlineStr
	goal.Amount = &g.Amount
	goal.Current = &g.Current
	statusStr := g.Status.String()
	goal.Status = &statusStr
	goal.Version = &g.Version
	goal.CreatedAt = &g.CreatedAt
	goal.Installments = g.Installments
	return goal
}

func (m *GoalDTO) ToModel() *Goal {
	goal := &Goal{}

	if m.ID != nil {
		goal.ID = *m.ID
	}
	if m.Name != nil {
		goal.Name = *m.Name
	}
	if m.Description != nil {
		goal.Description = *m.Description
	}
	if m.Color != nil {
		goal.Color = *m.Color
	}
	if m.User != nil {
		goal.User = m.User.ToModel()
	}
	if m.Deadline != nil {
		parsedTime, err := time.Parse("02/01/2006", *m.Deadline)
		if err != nil {
			parsedTime = time.Now()
		}
		goal.Deadline = parsedTime
	}
	if m.Amount != nil {
		goal.Amount = *m.Amount
	}
	if m.Current != nil {
		goal.Current = *m.Current
	}
	if m.Status != nil {
		goal.Status = GoalStatusFromString(*m.Status)
	}
	return goal
}

func (g *GoalProgress) ToDTO() *GoalProgressDTO {
	goalProgress := &GoalProgressDTO{}
	goalProgress.ID = &g.ID
	goalProgress.Goal = g.Goal.ToDTO()
	goalProgress.Amount = &g.Amount
	goalProgress.Date = &g.Date
	goalProgress.Version = g.Version
	goalProgress.CreatedAt = g.CreatedAt
	goalProgress.Deleted = &g.Deleted
	return goalProgress
}

func (g *GoalProgressDTO) ToModel() *GoalProgress {
	goalProgress := &GoalProgress{}
	if g.ID != nil {
		goalProgress.ID = *g.ID
	}
	if g.Goal != nil {
		goalProgress.Goal = g.Goal.ToModel()
	}
	if g.Amount != nil {
		goalProgress.Amount = *g.Amount
	}
	if g.Date != nil {
		goalProgress.Date = *g.Date
	}
	if g.Deleted != nil {
		goalProgress.Deleted = *g.Deleted
	}
	if g.Version != nil {
		goalProgress.Version = g.Version
	}
	if g.CreatedAt != nil {
		goalProgress.CreatedAt = g.CreatedAt
	}
	return goalProgress
}

func (g *Goal) ValidateGoal(v *validator.Validator) {
	v.Check(g.Name != "", "name", "must be provided")
	v.Check(len(g.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(g.Description != "", "description", "must be provided")
	v.Check(len(g.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	v.Check(g.Status.String() != "", "status", "must be provided")
	v.Check(g.Status.String() != "Unknown", "status", "invalid status value")
	v.Check(g.Color != "", "color", "must be provided")
	v.Check(g.Amount != 0, "amount", "must be provided")
}

func (g *GoalProgress) ValidateGoalProgress(v *validator.Validator) {
	v.Check(g.Amount != 0, "amount", "must be provided")
	v.Check(g.Goal != nil, "goal", "must be provided")
}
