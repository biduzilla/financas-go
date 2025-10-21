package handler

import (
	"database/sql"
	"financas/internal/config"
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils/errors"
	"net/http"
)

type Handler struct {
	User     UserHandlerInterface
	Auth     AuthHandlerInterface
	Category CategoryHandlerInterface
	errResp  errors.ErrorResponseInterface
	Service  *service.Service
}

func NewHandler(db *sql.DB, errResp errors.ErrorResponseInterface, config config.Config, ContextGetUser func(r *http.Request) *model.User) *Handler {
	service := service.NewService(db, config)

	return &Handler{
		User:     NewUserHandler(service.User, errResp),
		Auth:     NewAuthHandler(service.Auth, errResp),
		Category: NewCategoryHandler(service.Category, ContextGetUser, errResp),
		errResp:  errResp,
		Service:  service,
	}
}
