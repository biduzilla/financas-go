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
}

type UserHandlerInterface interface {
	ActivateUserHandler(w http.ResponseWriter, r *http.Request)
	CreateUserHandler(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(userService service.UserServiceInterface, errResp e.ErrorResponseInterface) *UserHandler {
	return &UserHandler{
		user:          userService,
		errorResponse: errResp,
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

	v := validator.New()

	user, err := h.user.ActivateUser(input.Cod, input.Email, v)

	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, v)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user.ToDTO()}, nil)
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

	v := validator.New()

	err = h.user.RegisterUserHandler(user, v)
	if err != nil {
		h.errorResponse.HandlerErrorResponse(w, r, err, v)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user.ToDTO()}, nil)
	if err != nil {
		h.errorResponse.ServerErrorResponse(w, r, err)
	}
}
