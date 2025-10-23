package handler

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"fmt"
	"net/http"
	"time"
)

type TransactionHandler struct {
	transaction    service.TransactionServiceInterface
	category       service.CategoryServiceInterface
	errRsp         e.ErrorResponseInterface
	contextGetUser func(r *http.Request) *model.User
}

type TransactionHandlerInterface interface {
	GetAllByUserAndCategory(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
}

func NewTransactionHandler(
	t service.TransactionServiceInterface,
	errRsp e.ErrorResponseInterface,
	contextGetUser func(r *http.Request) *model.User,
	c service.CategoryServiceInterface,
) *TransactionHandler {
	return &TransactionHandler{
		transaction:    t,
		errRsp:         errRsp,
		contextGetUser: contextGetUser,
		category:       c,
	}
}

func (h *TransactionHandler) GetAllByUserAndCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.errRsp)
	if !ok {
		return
	}

	var input struct {
		Name       string
		CategoryID int
		StartDate  *time.Time
		EndDate    *time.Time
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Name = utils.ReadString(qs, "description", "")
	input.StartDate = utils.ReadDate(qs, "start", "2006-01-02")
	input.EndDate = utils.ReadDate(qs, "end", "2006-01-02")
	input.Filters.Page = utils.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = utils.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "description", "-id", "-description"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		h.errRsp.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user := h.contextGetUser(r)

	t, m, err := h.transaction.GetAllByUserAndCategory(v, input.Name, user.ID, id, input.StartDate, input.EndDate, input.Filters)

	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	transactionsDTO := []*model.TransactionDTO{}

	for _, t := range t {
		err = h.prepareTransactionForResponse(t, user)
		if err != nil {
			h.errRsp.ServerErrorResponse(w, r, err)
			return
		}
		transactionsDTO = append(transactionsDTO, t.ToDTO())
	}

	respond(w, r, http.StatusOK, utils.Envelope{"transactions": transactionsDTO, "metadata": m}, nil, h.errRsp)
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.errRsp)
	if !ok {
		return
	}

	user := h.contextGetUser(r)
	t, err := h.transaction.GetByID(id, user.ID)

	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	h.prepareTransactionForResponse(t, user)

	respond(w, r, http.StatusOK, utils.Envelope{"transaction": t.ToDTO()}, nil, h.errRsp)
}

func (h *TransactionHandler) Save(w http.ResponseWriter, r *http.Request) {
	var dto = &model.TransactionDTO{}
	if err := utils.ReadJSON(w, r, dto); err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	t := dto.ToModel()

	err := h.transaction.Save(v, t)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	user := h.contextGetUser(r)

	h.prepareTransactionForResponse(t, user)

	headers := http.Header{"Location": {fmt.Sprintf("/v1/transactions/%d", t.ID)}}
	respond(w, r, http.StatusCreated, utils.Envelope{"category": t.ToDTO()}, headers, h.errRsp)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto = &model.TransactionDTO{}
	if err := utils.ReadJSON(w, r, dto); err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	t := dto.ToModel()
	user := h.contextGetUser(r)

	err := h.transaction.Update(v, t, user.ID)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	h.prepareTransactionForResponse(t, user)

	respond(w, r, http.StatusOK, utils.Envelope{"category": t.ToDTO()}, nil, h.errRsp)
}

func (h *TransactionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.errRsp)
	if !ok {
		return
	}

	user := h.contextGetUser(r)
	if err := h.transaction.Delete(id, user.ID); err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}
	respond(w, r, http.StatusNoContent, utils.Envelope{"message": "transaction successfully deleted"}, nil, h.errRsp)
}

func (h *TransactionHandler) prepareTransactionForResponse(transaction *model.Transaction, user *model.User) error {
	transaction.User = user

	category, err := h.category.GetByID(transaction.Category.ID, user.ID)

	if err != nil {
		return err
	}

	if category != nil {
		transaction.Category = category
	}

	return nil
}
