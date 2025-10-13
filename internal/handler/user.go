package handler

import (
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"net/http"
)

type UserHandler struct {
	user          service.UserServiceInterface
	errorResponse e.ErrorResponseInterface
	validator     *validator.Validator
}

type UserHandlerInterface interface {
	ActivateUserHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(userService service.UserServiceInterface, errResp e.ErrorResponseInterface, v *validator.Validator) *UserHandler {
	return &UserHandler{
		user:          userService,
		errorResponse: errResp,
		validator:     v,
	}
}

func (h *UserHandler) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Cod   int    `json:"cod"`
		Email string `json:"email"`
	}

	err := utils.ReadJSON(w, r, &input)

	if err != nil {
		h.errorResponse.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.user.ActivateUser(input.Cod, input.Email)

	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, h.validator)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}, nil)
	if err != nil {
		h.errorResponse.ServerErrorResponse(w, r, err)
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO model.UserSaveDTO
	err := utils.ReadJSON(w, r, &userDTO)
	if err != nil {
		h.errorResponse.BadRequestResponse(w, r, err)
		return
	}

	user, err := userDTO.ToModel()
	if err != nil {
		h.errorResponse.ServerErrorResponse(w, r, err)
		return
	}

	err = h.user.RegisterUserHandler(user)
	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, h.validator)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}, nil)
	if err != nil {
		h.errorResponse.ServerErrorResponse(w, r, err)
	}
}
