package handler

import (
	"database/sql"
	"financas/internal/config"
	"financas/internal/service"
	"financas/utils/errors"
)

type Handler struct {
	User    UserHandlerInterface
	Auth    AuthHandlerInterface
	errResp errors.ErrorResponseInterface
	Service *service.Service
}

func NewHandler(db *sql.DB, errResp errors.ErrorResponseInterface, config config.Config) *Handler {
	service := service.NewService(db, config)

	return &Handler{
		User:    NewUserHandler(service.User, errResp),
		Auth:    NewAuthHandler(service.Auth, errResp),
		errResp: errResp,
		Service: service,
	}
}
