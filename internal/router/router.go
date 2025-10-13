package router

import (
	"database/sql"
	"financas/internal/handler"
	"financas/internal/jsonlog"
	"financas/utils/errors"
	"net/http"

	"github.com/go-chi/chi"
)

type Router struct {
	User    *UserRouter
	errResp errors.ErrorResponseInterface
}

func NewRouter(db *sql.DB, logger *jsonlog.Logger) *Router {
	e := errors.NewErrorResponse(logger)
	h := handler.NewHandler(db, e)

	return &Router{
		User:    NewUserRouter(h.User),
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
	})
	return r
}
