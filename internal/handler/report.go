package handler

import (
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"net/http"
	"strconv"
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

	startDate, endDate, err := parseDateRange(r, v)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	if !v.Valid() {
		h.errRsp.FailedValidationResponse(w, r, v.Errors)
		return
	}

	summary, err := h.report.GetFinancialSummary(v, user.ID, startDate, endDate)
	if err != nil {
		h.errRsp.ServerErrorResponse(w, r, err)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"summary": summary}, nil, h.errRsp)
}

func (h *ReportHandler) GetCategoryReportHandler(w http.ResponseWriter, r *http.Request) {
	user := h.contextGetUser(r)
	v := validator.New()

	startDate, endDate, err := parseDateRange(r, v)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

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

	startDate, endDate, err := parseDateRange(r, v)
	if err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	// Parse limit
	limit := 5
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Parse type
	typeStr := r.URL.Query().Get("type")
	var categoryType model.TypeCategoria
	if typeStr == "RECEITA" {
		categoryType = model.RECEITA
	} else {
		categoryType = model.DESPESA
	}

	topCategories, err := h.report.GetTopCategories(v, user.ID, startDate, endDate, limit, categoryType)
	if err != nil {
		h.errRsp.ServerErrorResponse(w, r, err)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"topCategories": topCategories}, nil, h.errRsp)
}
