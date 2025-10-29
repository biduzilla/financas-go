package handler

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"net/http"
)

type GoalHandler struct {
	goal           service.GoalServiceInterface
	errRsp         e.ErrorResponseInterface
	contextGetUser func(r *http.Request) *model.User
}

type GoalHandlerInterface interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewGoalHandler(
	goal service.GoalServiceInterface,
	errRsp e.ErrorResponseInterface,
	contextGetUser func(r *http.Request) *model.User,
) *GoalHandler {
	return &GoalHandler{
		goal:           goal,
		errRsp:         errRsp,
		contextGetUser: contextGetUser,
	}
}

func (h *GoalHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var input struct {
		name string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.name = utils.ReadString(qs, "name", "")
	input.Filters.Page = utils.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = utils.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = utils.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		h.errRsp.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user := h.contextGetUser(r)
	goals, metadata, err := h.goal.GetAllByUserId(input.name, user.ID, input.Filters, v)

	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	goalsDTO := make([]*model.GoalDTO, len(goals))

	for _, g := range goals {
		goalsDTO = append(goalsDTO, g.ToDTO())
	}

	respond(
		w,
		r,
		http.StatusOK,
		utils.Envelope{
			"goals":    goalsDTO,
			"metadata": metadata,
		},
		nil,
		h.errRsp,
	)
}

func (h *GoalHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.errRsp)
	if !ok {
		return
	}
	v := validator.New()
	user := h.contextGetUser(r)
	goal, err := h.goal.GetById(v, id, user.ID)

	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	respond(
		w,
		r,
		http.StatusOK,
		utils.Envelope{"goal": goal.ToDTO()},
		nil,
		h.errRsp,
	)
}

func (h *GoalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto model.GoalDTO
	if err := utils.ReadJSON(w, r, &dto); err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	user := h.contextGetUser(r)
	goal := dto.ToModel()
	goal.User = user
	if err := h.goal.Create(v, goal); err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	respond(
		w,
		r,
		http.StatusCreated,
		utils.Envelope{"goal": goal.ToDTO()},
		nil,
		h.errRsp,
	)
}

func (h *GoalHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto model.GoalDTO
	if err := utils.ReadJSON(w, r, &dto); err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	user := h.contextGetUser(r)
	goal := dto.ToModel()
	goal.User = user
	if err := h.goal.Update(v, goal, user.ID); err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	respond(
		w,
		r,
		http.StatusOK,
		utils.Envelope{"goal": goal.ToDTO()},
		nil,
		h.errRsp,
	)
}

func (h *GoalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.errRsp)
	if !ok {
		return
	}
	user := h.contextGetUser(r)
	if err := h.goal.Delete(id, user.ID); err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	respond(
		w,
		r,
		http.StatusNoContent,
		nil,
		nil,
		h.errRsp,
	)
}
