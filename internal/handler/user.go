package handler

import (
	"errors"
	"financas/internal/service"
	"financas/utils"
	e "financas/utils/errors"
	"financas/utils/validator"
	"net/http"
)

type UserHandler struct {
	userService   service.UserServiceInterface
	errorResponse *e.ErrorResponse
	validator     *validator.Validator
}

func NewUserHandler(userService service.UserServiceInterface, errResp *e.ErrorResponse, v *validator.Validator) *UserHandler {
	return &UserHandler{
		userService:   userService,
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

	user, err := h.userService.ActivateUser(input.Cod, input.Email)

	if err != nil {
		switch {
		case errors.Is(err, validator.ErrInvalidData):
			h.errorResponse.FailedValidationResponse(w, r, h.validator.Errors)
		default:
			h.errorResponse.ServerErrorResponse(w, r, err)
		}
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}, nil)
	if err != nil {
		h.errorResponse.ServerErrorResponse(w, r, err)
	}
}
