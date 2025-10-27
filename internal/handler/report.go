package handler

import (
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"net/http"
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
}

func NewReportHandler(report service.ReportServiceInterface, errResp e.ErrorResponseInterface, contextGetUser func(r *http.Request) *model.User) *ReportHandler {
	return &ReportHandler{
		report:         report,
		errRsp:         errResp,
		contextGetUser: contextGetUser,
	}
}

func (h *ReportHandler) GetFinancialSummaryHandler(w http.ResponseWriter, r *http.Request) {
	user := h.contextGetUser(r)
	v := validator.New()
	qs := r.URL.Query()

	startDate := utils.ReadDate(qs, "start", "2006-01-02")
	endDate := utils.ReadDate(qs, "end", "2006-01-02")

	summary, err := h.report.GetFinancialSummary(v, user.ID, startDate, endDate)
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

	startDate := utils.ReadDate(qs, "start", "2006-01-02")
	endDate := utils.ReadDate(qs, "end", "2006-01-02")

	report, err := h.report.GetCategoryReport(v, user.ID, startDate, endDate)
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

	startDate := utils.ReadDate(qs, "start", "2006-01-02")
	endDate := utils.ReadDate(qs, "end", "2006-01-02")
	limit := utils.ReadInt(qs, "limit", 5, v)
	categoryType := model.TypeCategoriaFromString(utils.ReadString(qs, "type", "RECEITA"))

	topCategories, err := h.report.GetTopCategories(v, user.ID, startDate, endDate, limit, categoryType)
	if err != nil {
		h.errRsp.ServerErrorResponse(w, r, err)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"topCategories": topCategories}, nil, h.errRsp)
}
