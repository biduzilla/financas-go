package handler

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"net/http"
	"time"
)

type ReportHandler struct {
	report         service.ReportServiceInterface
	errRsp         e.ErrorResponseInterface
	contextGetUser func(r *http.Request) *model.User
}

type ReportHandlerInterface interface {
	GetFinancialSummaryHandler(w http.ResponseWriter, r *http.Request)
	GetCategoryReportHandler(w http.ResponseWriter, r *http.Request)
	GetTopCategoriesHandler(w http.ResponseWriter, r *http.Request)
	GetIncomeVsExpensesHandler(w http.ResponseWriter, r *http.Request)
}

func NewReportHandler(report service.ReportServiceInterface, errResp e.ErrorResponseInterface, contextGetUser func(r *http.Request) *model.User) *ReportHandler {
	return &ReportHandler{
		report:         report,
		errRsp:         errResp,
		contextGetUser: contextGetUser,
	}
}

func (h *ReportHandler) GetIncomeVsExpensesHandler(w http.ResponseWriter, r *http.Request) {
	user := h.contextGetUser(r)
	v := validator.New()
	qs := r.URL.Query()

	var input struct {
		StartDate *time.Time
		EndDate   *time.Time
		filters.Filters
	}

	input.StartDate = utils.ReadDate(qs, "start", "2006-01-02")
	input.EndDate = utils.ReadDate(qs, "end", "2006-01-02")
	input.Filters.Page = utils.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = utils.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "description", "-id", "-description"}

	filters.ValidateFilters(v, input.Filters)
	if !v.Valid() {
		h.errRsp.FailedValidationResponse(w, r, v.Errors)
		return
	}

	incomeVsExpenses, err := h.report.GetIncomeVsExpenses(v, user.ID, input.StartDate, input.EndDate, input.Filters)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{
		"income_vs_expenses": incomeVsExpenses,
	}, nil, h.errRsp)
}

func (h *ReportHandler) GetFinancialSummaryHandler(w http.ResponseWriter, r *http.Request) {
	user := h.contextGetUser(r)
	v := validator.New()
	qs := r.URL.Query()

	var input struct {
		StartDate *time.Time
		EndDate   *time.Time
		filters.Filters
	}

	input.StartDate = utils.ReadDate(qs, "start", "2006-01-02")
	input.EndDate = utils.ReadDate(qs, "end", "2006-01-02")
	input.Filters.Page = 1
	input.Filters.PageSize = 9999
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "description", "-id", "-description"}

	summary, err := h.report.GetFinancialSummary(v, user.ID, input.StartDate, input.EndDate, input.Filters)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"summary": summary}, nil, h.errRsp)
}

func (h *ReportHandler) GetCategoryReportHandler(w http.ResponseWriter, r *http.Request) {
	user := h.contextGetUser(r)
	v := validator.New()
	qs := r.URL.Query()

	var input struct {
		StartDate *time.Time
		EndDate   *time.Time
		filters.Filters
	}

	input.StartDate = utils.ReadDate(qs, "start", "2006-01-02")
	input.EndDate = utils.ReadDate(qs, "end", "2006-01-02")
	input.Filters.Page = 1
	input.Filters.PageSize = 9999
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "description", "-id", "-description"}

	report, err := h.report.GetCategoryReport(v, user.ID, input.StartDate, input.EndDate, input.Filters)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"report": report}, nil, h.errRsp)
}

func (h *ReportHandler) GetTopCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	user := h.contextGetUser(r)
	v := validator.New()
	qs := r.URL.Query()

	var input struct {
		StartDate *time.Time
		EndDate   *time.Time
		Limit     int
		typeStr   string
		filters.Filters
	}

	input.StartDate = utils.ReadDate(qs, "start", "2006-01-02")
	input.EndDate = utils.ReadDate(qs, "end", "2006-01-02")
	input.Filters.Page = 1
	input.Filters.PageSize = 9999
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "description", "-id", "-description"}
	input.Limit = utils.ReadInt(qs, "limit", 5, v)

	if !v.Valid() {
		h.errRsp.FailedValidationResponse(w, r, v.Errors)
		return
	}

	input.typeStr = utils.ReadString(qs, "type", "RECEITA")
	categoryType := model.TypeCategoriaFromString(input.typeStr)

	topCategories, err := h.report.GetTopCategories(v, user.ID, input.StartDate, input.EndDate, input.Limit, categoryType, input.Filters)
	if err != nil {
		h.errRsp.ServerErrorResponse(w, r, err)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"topCategories": topCategories}, nil, h.errRsp)
}
