package router

import (
	"financas/internal/handler"
	"financas/internal/middleware"

	"github.com/go-chi/chi"
)

type GoalRouter struct {
	goal handler.GoalHandlerInterface
	m    middleware.MiddlewareInterface
}

func NewGoalRouter(h handler.GoalHandlerInterface, m middleware.MiddlewareInterface) *GoalRouter {
	return &GoalRouter{
		goal: h,
		m:    m,
	}
}

type GoalRouterInterface interface {
	GoalRoutes(r chi.Router)
}

func (g *GoalRouter) GoalRoutes(r chi.Router) {
	r.Route("/goals", func(r chi.Router) {
		r.Use(g.m.RequireActivatedUser)

		r.Get("/{id}", g.goal.GetById)
		r.Get("/", g.goal.GetAll)
		r.Post("/", g.goal.Create)
		r.Put("/{id}", g.goal.Update)
		r.Delete("/{id}", g.goal.Delete)
	})
}
