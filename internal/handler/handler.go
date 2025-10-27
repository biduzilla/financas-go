package handler

import (
	"database/sql"
	"financas/internal/config"
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils"
	"financas/utils/errors"
	"financas/utils/validator"
	"net/http"
	"time"
)

type Handler struct {
	User        UserHandlerInterface
	Auth        AuthHandlerInterface
	Category    CategoryHandlerInterface
	Report      ReportHandlerInterface
	Transaction TransactionHandlerInterface
	errResp     errors.ErrorResponseInterface
	Service     *service.Service
}

func NewHandler(db *sql.DB, errResp errors.ErrorResponseInterface, config config.Config, ContextGetUser func(r *http.Request) *model.User) *Handler {
	service := service.NewService(db, config)

	return &Handler{
		User:        NewUserHandler(service.User, errResp),
		Auth:        NewAuthHandler(service.Auth, errResp),
		Category:    NewCategoryHandler(service.Category, ContextGetUser, errResp),
		Transaction: NewTransactionHandler(service.Transaction, errResp, ContextGetUser, service.Category),
		Report:      NewReportHandler(service.Report, errResp, ContextGetUser),
		errResp:     errResp,
		Service:     service,
	}
}

func parseID(
	w http.ResponseWriter,
	r *http.Request,
	errRsp errors.ErrorResponseInterface,
) (int64, bool) {
	id, err := utils.ReadIntParam(r, "id")
	if err != nil {
		errRsp.BadRequestResponse(w, r, err)
		return 0, false
	}
	return id, true
}

func respond(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data utils.Envelope,
	headers http.Header,
	errRsp errors.ErrorResponseInterface,
) {
	err := utils.WriteJSON(w, status, data, headers)
	if err != nil {
		errRsp.ServerErrorResponse(w, r, err)
	}
}

func parseDateRange(r *http.Request, v *validator.Validator) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	qs := r.URL.Query()

	endDatePtr := utils.ReadDate(qs, "end_date", "2006-01-02")
	startDatePtr := utils.ReadDate(qs, "start_date", "2006-01-02")

	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
	if endDatePtr != nil {
		endDate = time.Date(endDatePtr.Year(), endDatePtr.Month(), endDatePtr.Day(), 23, 59, 59, 0, time.UTC)
	}

	startDate := endDate.AddDate(0, 0, -30) // 30 dias antes da data final
	if startDatePtr != nil {
		startDate = time.Date(startDatePtr.Year(), startDatePtr.Month(), startDatePtr.Day(), 0, 0, 0, 0, time.UTC)
	}

	if startDate.After(endDate) {
		v.AddError("date_range", "start_date cannot be after end_date")
		return time.Time{}, time.Time{}, errors.ErrInvalidData
	}

	return startDate, endDate, nil
}
