package handler

import (
	"database/sql"
	"financas/internal/jsonlog"
	"financas/internal/service"
	"financas/utils/errors"
	"financas/utils/validator"
)

type Handler struct {
	User    UserHandlerInterface
	errResp errors.ErrorResponseInterface
}

func NewHandler(db *sql.DB, v *validator.Validator, logger *jsonlog.Logger) *Handler {
	service := service.NewService(db, v)
	errorResponse := errors.NewErrorResponse(logger)

	return &Handler{
		User:    NewUserHandler(service.User, errorResponse, v),
		errResp: errorResponse,
	}
}
