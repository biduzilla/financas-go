package service

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/utils/validator"
	"time"
)

type ReportService struct {
	transaction TransactionServiceInterface
	category    CategoryServiceInterface
}

type ReportServiceInterface interface {
	GetFinancialSummary(
		v *validator.Validator,
		userID int64,
		startDate, endDate *time.Time,
		f filters.Filters,
	) (*model.FinancialSummary, error)

	GetCategoryReport(
		v *validator.Validator,
		userID int64,
		startDate, endDate *time.Time,
		f filters.Filters,
	) ([]model.CategorySummary, error)

	GetIncomeVsExpenses(
		v *validator.Validator,
		userID int64,
		startDate, endDate *time.Time,
		f filters.Filters,
	) (map[string]float64, error)

	GetTopCategories(
		v *validator.Validator,
		userID int64,
		startDate, endDate *time.Time,
		limit int,
		categoryType model.TypeCategoria,
		f filters.Filters,
	) ([]model.CategorySummary, error)
}

func NewReportService(transactionSvc TransactionServiceInterface, categorySvc CategoryServiceInterface) *ReportService {
	return &ReportService{
		transaction: transactionSvc,
		category:    categorySvc,
	}
}

func (s *ReportService) GetFinancialSummary(
	v *validator.Validator,
	userID int64,
	startDate,
	endDate *time.Time,
	f filters.Filters,
) (*model.FinancialSummary, error) {
	transactions, _, err := s.transaction.GetAllByUserAndCategory(
		v,
		"",
		userID,
		0,
		startDate,
		endDate,
		f,
	)

	if err != nil {
		return nil, err
	}

	if startDate == nil {
		temp := time.Now().AddDate(0, -1, 0)
		startDate = &temp
	}

	if endDate == nil {
		temp := time.Now().AddDate(0, 0, 0).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endDate = &temp
	}

	var totalIncome, totalExpenses float64
	categoryTotals := make(map[int64]*model.CategorySummary)

	for _, transaction := range transactions {
		amount := transaction.Amount
		category := transaction.Category

		switch category.Type {
		case model.RECEITA:
			totalIncome += amount
		case model.DESPESA:
			totalExpenses += amount
		}

		if _, exist := categoryTotals[category.ID]; !exist {
			categoryTotals[category.ID] = &model.CategorySummary{
				Category: category.ToDTO(),
				Total:    0,
				Count:    0,
			}
		}

		categoryTotals[category.ID].Total += amount
		categoryTotals[category.ID].Count++
	}

	categorySummary := make([]model.CategorySummary, 0, len(categoryTotals))
	for _, summary := range categoryTotals {
		var totalForPercentage float64
		if model.TypeCategoriaFromString(*summary.Category.Type) == model.RECEITA {
			totalForPercentage = totalIncome
		} else {
			totalForPercentage = totalExpenses
		}

		if totalForPercentage > 0 {
			summary.Percentage = (summary.Total / totalForPercentage) * 100
		}

		categorySummary = append(categorySummary, *summary)
	}

	monthlyTrends, err := s.generateMonthlyTrends(v, userID, 6)
	if err != nil {
		return nil, err
	}

	summary := &model.FinancialSummary{
		TotalIncome:     totalIncome,
		TotalExpenses:   totalExpenses,
		Balance:         totalIncome - totalExpenses,
		CategorySummary: categorySummary,
		MonthlyTrends:   monthlyTrends,
		Period: model.PeriodSummary{
			StartDate: *startDate,
			EndDate:   *endDate,
			Days:      int(endDate.Sub(*startDate).Hours() / 24),
		},
	}

	return summary, nil
}

func (s *ReportService) GetCategoryReport(
	v *validator.Validator,
	userID int64,
	startDate, endDate *time.Time,
	f filters.Filters,
) ([]model.CategorySummary, error) {
	summary, err := s.GetFinancialSummary(v, userID, startDate, endDate, f)
	if err != nil {
		return nil, err
	}
	return summary.CategorySummary, nil
}

func (s *ReportService) generateMonthlyTrends(v *validator.Validator, userID int64, months int) ([]model.MonthlyTrend, error) {
	trends := make([]model.MonthlyTrend, 0, months)
	now := time.Now()

	for i := months - 1; i >= 0; i-- {
		monthStart := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, -1)

		transactions, _, err := s.transaction.GetAllByUserAndCategory(
			v,
			"",
			userID,
			0,
			&monthStart,
			&monthEnd,
			filters.Filters{
				Page:         1,
				PageSize:     100,
				Sort:         "created_at",
				SortSafelist: []string{"created_at", "amount", "description"},
			},
		)

		if err != nil {
			return nil, err
		}

		var monthIncome, monthExpenses float64
		for _, transaction := range transactions {
			if transaction.Category.Type == model.RECEITA {
				monthIncome += transaction.Amount
			} else {
				monthExpenses += transaction.Amount
			}
		}

		trend := model.MonthlyTrend{
			Month:    monthStart.Format("Jan/2006"),
			Income:   monthIncome,
			Expenses: monthExpenses,
			Balance:  monthIncome - monthExpenses,
		}

		trends = append(trends, trend)
	}

	return trends, nil
}

func (s *ReportService) GetIncomeVsExpenses(
	v *validator.Validator,
	userID int64,
	startDate, endDate *time.Time,
	f filters.Filters,
) (map[string]float64, error) {
	summary, err := s.GetFinancialSummary(v, userID, startDate, endDate, f)
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"income":   summary.TotalIncome,
		"expenses": summary.TotalExpenses,
		"balance":  summary.Balance,
	}, nil
}

func (s *ReportService) GetTopCategories(
	v *validator.Validator,
	userID int64,
	startDate, endDate *time.Time,
	limit int,
	categoryType model.TypeCategoria,
	f filters.Filters,
) ([]model.CategorySummary, error) {
	allCategories, err := s.GetCategoryReport(v, userID, startDate, endDate, f)
	if err != nil {
		return nil, err
	}

	filtered := make([]model.CategorySummary, 0)
	for _, cat := range allCategories {
		if model.TypeCategoriaFromString(*cat.Category.Type) == categoryType {
			filtered = append(filtered, cat)
		}
	}

	for i := 0; i < len(filtered)-1; i++ {
		for j := i + 1; j < len(filtered); j++ {
			if filtered[i].Total < filtered[j].Total {
				filtered[i], filtered[j] = filtered[j], filtered[i]
			}
		}
	}

	if limit > 0 && limit < len(filtered) {
		filtered = filtered[:limit]
	}

	return filtered, nil
}
