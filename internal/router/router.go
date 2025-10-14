package router

import (
	"database/sql"
	"financas/configuration"
	"financas/internal/handler"
	"financas/internal/jsonlog"
	"financas/internal/middleware"
	"financas/internal/model"
	"financas/utils"
	"financas/utils/errors"
	"net/http"

	"github.com/go-chi/chi"
)

type Router struct {
	User           UserRoutesInterface
	Auth           AuthRoutesInterface
	ErrResp        errors.ErrorResponseInterface
	ContextGetUser func(r *http.Request) *model.User
	ContextSetUser func(r *http.Request, user *model.User) *http.Request
	Handler        *handler.Handler
}

func NewRouter(
	db *sql.DB,
	logger *jsonlog.Logger,
	Config *configuration.Conf,
	contextGetUser func(r *http.Request) *model.User,
	contextSetUser func(r *http.Request, user *model.User) *http.Request,
) *Router {
	e := errors.NewErrorResponse(logger)
	h := handler.NewHandler(db, e, Config)

	return &Router{
		User:           NewUserRouter(h.User),
		Auth:           NewAuthRouter(h.Auth),
		ErrResp:        e,
		ContextGetUser: contextGetUser,
		ContextSetUser: contextSetUser,
		Handler:        h,
	}
}

func (router *Router) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()
	m := middleware.New(
		router.ErrResp,
		router.ContextGetUser,
		router.ContextSetUser,
		router.Handler.Service.Auth,
		router.Handler.Service.User,
	)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		router.ErrResp.NotFoundResponse(w, req)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		router.ErrResp.MethodNotAllowedResponse(w, req)
	})

	r.Route("/v1", func(r chi.Router) {
		router.User.UserRoutes(r)
		router.Auth.AuthRoutes(r)

		r.Route("/healthcheck", func(r chi.Router) {
			r.Use(m.Authenticate)
			r.Use(m.RequireActivatedUser)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				env := map[string]any{
					"status": "available",
					"system_info": map[string]string{
						"environment": "development",
						"version":     "v1",
					},
				}

				err := utils.WriteJSON(w, http.StatusOK, env, nil)
				if err != nil {
					router.ErrResp.ServerErrorResponse(w, r, err)
				}
			})
		})

	})

	return r
}
