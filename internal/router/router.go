package router

import (
	"database/sql"
	"financas/configuration"
	"financas/internal/handler"
	"financas/internal/jsonlog"
	"financas/utils/errors"
	"net/http"

	"github.com/go-chi/chi"
)

type Router struct {
	User    UserRoutesInterface
	Auth    AuthRoutesInterface
	errResp errors.ErrorResponseInterface
}

func NewRouter(db *sql.DB, logger *jsonlog.Logger, config *configuration.Conf) *Router {
	e := errors.NewErrorResponse(logger)
	h := handler.NewHandler(db, e, config)

	return &Router{
		User:    NewUserRouter(h.User),
		Auth:    NewAuthRouter(h.Auth),
		errResp: e,
	}
}

func (router *Router) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		router.errResp.NotFoundResponse(w, req)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		router.errResp.MethodNotAllowedResponse(w, req)
	})

	r.Route("/v1", func(r chi.Router) {
		router.User.UserRoutes(r)
		router.Auth.AuthRoutes(r)
	})
	return r
}
