package handler

import (
	"database/sql"
	"financas/internal/config"
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils"
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

func parseID(
	w http.ResponseWriter,
	r *http.Request,
	errRsp errors.ErrorResponseInterface,
) (int64, bool) {
	id, err := utils.ReadIntParam(r, "id")
	if err != nil {
		errRsp.BadRequestResponse(w, r, err)
		return 0, false
	}
	return id, true
}

func respond(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data utils.Envelope,
	headers http.Header,
	errRsp errors.ErrorResponseInterface,
) {
	err := utils.WriteJSON(w, status, data, headers)
	if err != nil {
		errRsp.ServerErrorResponse(w, r, err)
	}
}
