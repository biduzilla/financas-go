package router

import (
	"database/sql"
	"expvar"
	"financas/internal/config"
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
	user           UserRoutesInterface
	auth           AuthRoutesInterface
	category       CategoryRouterInterface
	transaction    TransactionRouterInterface
	goal           GoalRouterInterface
	goalProgress   GoalProgressRouterInterface
	report         ReportRouterInterface
	ErrResp        errors.ErrorResponseInterface
	ContextGetUser func(r *http.Request) *model.User
	ContextSetUser func(r *http.Request, user *model.User) *http.Request
	Handler        *handler.Handler
	Config         config.Config
	m              middleware.MiddlewareInterface
}

func NewRouter(
	db *sql.DB,
	logger *jsonlog.Logger,
	contextGetUser func(r *http.Request) *model.User,
	contextSetUser func(r *http.Request, user *model.User) *http.Request,
	config config.Config,
) *Router {
	e := errors.NewErrorResponse(logger)
	h := handler.NewHandler(db, e, config, contextGetUser)
	m := middleware.New(
		e,
		contextGetUser,
		contextSetUser,
		h.Service.Auth,
		h.Service.User,
		config,
	)
	return &Router{
		ErrResp:        e,
		ContextGetUser: contextGetUser,
		ContextSetUser: contextSetUser,
		Handler:        h,
		m:              m,
		user:           NewUserRouter(h.User),
		auth:           NewAuthRouter(h.Auth),
		category:       NewCategoryRouter(h.Category, m),
		transaction:    NewTransactionRouter(h.Transaction, m),
		report:         NewReportRouter(h.Report, m),
		goal:           NewGoalRouter(h.Goal, m),
	}
}

func (router *Router) RegisterRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(router.m.RecoverPanic)
	r.Use(router.m.Metrics)
	r.Use(router.m.RateLimit)
	r.Use(router.m.EnableCORS)
	r.Use(router.m.Authenticate)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		router.ErrResp.NotFoundResponse(w, req)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		router.ErrResp.MethodNotAllowedResponse(w, req)
	})

	r.Route("/v1", func(r chi.Router) {
		r.Mount("/debug/vars", expvar.Handler())
		router.user.UserRoutes(r)
		router.auth.AuthRoutes(r)
		router.category.CategoryRoutes(r)
		router.transaction.TransactionRoutes(r)
		router.report.ReportRoutes(r)
		router.goal.GoalRoutes(r)
		router.goalProgress.GoalProgressRoutes(r)

		r.Route("/healthcheck", func(r chi.Router) {
			r.Use(router.m.RequireActivatedUser)
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
