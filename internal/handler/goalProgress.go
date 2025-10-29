package handler

import (
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"net/http"
)

type GoalProgressHandler struct {
	gP             service.GoalProgressServiceInterface
	errRsp         e.ErrorResponseInterface
	contextGetUser func(r *http.Request) *model.User
}

type GoalProgressHandlerInterface interface {
	GetByGoalID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

func NewGoalProgressHandler(
	gP service.GoalProgressServiceInterface,
	errRsp e.ErrorResponseInterface,
	contextGetUser func(r *http.Request) *model.User,
) *GoalProgressHandler {
	return &GoalProgressHandler{
		gP:             gP,
		errRsp:         errRsp,
		contextGetUser: contextGetUser,
	}
}

func (h *GoalProgressHandler) GetByGoalID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.errRsp)
	if !ok {
		return
	}
	user := h.contextGetUser(r)
	gPs, err := h.gP.GetGoalProgressIDGoal(user.ID, id)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	gPsDTO := make([]*model.GoalProgressDTO, len(gPs))
	for i, gP := range gPs {
		gPsDTO[i] = gP.ToDTO()
	}

	respond(w, r, http.StatusOK, utils.Envelope{"goal_progress": gPsDTO}, nil, h.errRsp)
}

func (h *GoalProgressHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto model.GoalProgressDTO
	if err := utils.ReadJSON(w, r, &dto); err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	user := h.contextGetUser(r)
	gP := dto.ToModel()

	if gP.Goal != nil {
		gP.Goal.User = user
	}

	err := h.gP.Insert(v, gP, user.ID)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(w, r, http.StatusCreated, utils.Envelope{"goal_progress": gP.ToDTO()}, nil, h.errRsp)
}

func (h *GoalProgressHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto model.GoalProgressDTO
	if err := utils.ReadJSON(w, r, &dto); err != nil {
		h.errRsp.BadRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	user := h.contextGetUser(r)
	gP := dto.ToModel()

	if gP.Goal != nil {
		gP.Goal.User = user
	}

	err := h.gP.Update(v, gP, user.ID)
	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, v)
		return
	}

	respond(w, r, http.StatusOK, utils.Envelope{"goal_progress": gP.ToDTO()}, nil, h.errRsp)
}

func (h *GoalProgressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, h.errRsp)
	if !ok {
		return
	}
	user := h.contextGetUser(r)
	v := validator.New()
	err := h.gP.Delete(v, id, user.ID)

	if err != nil {
		h.errRsp.HandlerErrorResponse(w, r, err, nil)
		return
	}

	respond(w, r, http.StatusNoContent, nil, nil, h.errRsp)
}
