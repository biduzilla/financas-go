package model

import "time"

type FinancialSummary struct {
	TotalIncome     float64           `json:"total_income"`
	TotalExpenses   float64           `json:"total_expenses"`
	Balance         float64           `json:"balance"`
	CategorySummary []CategorySummary `json:"category_summary"`
	MonthlyTrends   []MonthlyTrend    `json:"monthly_trends"`
	Period          PeriodSummary     `json:"period"`
}

type CategorySummary struct {
	Category   *CategoryDTO `json:"category"`
	Total      float64      `json:"total"`
	Count      int          `json:"count"`
	Percentage float64      `json:"percentage"`
}

type MonthlyTrend struct {
	Month    string  `json:"month"`
	Income   float64 `json:"income"`
	Expenses float64 `json:"expenses"`
	Balance  float64 `json:"balance"`
}

type PeriodSummary struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Days      int       `json:"days"`
}
