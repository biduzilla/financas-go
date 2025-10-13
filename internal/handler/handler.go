package handler

import (
	"database/sql"
	"financas/internal/service"
	"financas/utils/errors"
)

type Handler struct {
	User    UserHandlerInterface
	errResp errors.ErrorResponseInterface
}

func NewHandler(db *sql.DB, errResp errors.ErrorResponseInterface) *Handler {
	service := service.NewService(db)

	return &Handler{
		User:    NewUserHandler(service.User, errResp),
		errResp: errResp,
	}
}
