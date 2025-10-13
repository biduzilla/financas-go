package handler

import (
	"financas/internal/service"
	"financas/utils"
	"financas/utils/errors"
	"financas/utils/validator"
	"net/http"
)

type AuthHandler struct {
	Auth          service.AuthServiceInterface
	ErrorResponse errors.ErrorResponseInterface
}

type AuthHandlerInterface interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
}

func NewAuthHandler(authService service.AuthServiceInterface, errResp errors.ErrorResponseInterface) *AuthHandler {
	return &AuthHandler{
		Auth:          authService,
		ErrorResponse: errResp,
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		h.ErrorResponse.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	token, err := h.Auth.Login(v, input.Email, input.Password)
	if err != nil {
		h.ErrorResponse.HandlerErrorResponse(w, r, err, v)
		return
	}

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"authentication_token": token}, nil)
	if err != nil {
		h.ErrorResponse.ServerErrorResponse(w, r, err)
	}
}
