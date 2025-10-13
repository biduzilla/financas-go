package handler

import (
	"database/sql"
	"financas/configuration"
	"financas/internal/service"
	"financas/utils/errors"
)

type Handler struct {
	User    UserHandlerInterface
	Auth    AuthHandlerInterface
	errResp errors.ErrorResponseInterface
}

func NewHandler(db *sql.DB, errResp errors.ErrorResponseInterface, config *configuration.Conf) *Handler {
	service := service.NewService(db, config)

	return &Handler{
		User:    NewUserHandler(service.User, errResp),
		Auth:    NewAuthHandler(service.Auth, errResp),
		errResp: errResp,
	}
}
