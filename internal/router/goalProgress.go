package router

import (
	"financas/internal/handler"
	"financas/internal/middleware"

	"github.com/go-chi/chi"
)

type GoalProgressRouter struct {
	handler handler.GoalProgressHandlerInterface
	m       middleware.MiddlewareInterface
}

func NewGoalProgressRouter(
	h handler.GoalProgressHandlerInterface,
	m middleware.MiddlewareInterface,
) *GoalProgressRouter {
	return &GoalProgressRouter{
		handler: h,
		m:       m,
	}
}

type GoalProgressRouterInterface interface {
	GoalProgressRoutes(r chi.Router)
}

func (g *GoalProgressRouter) GoalProgressRoutes(r chi.Router) {
	r.Route("/goal_progress", func(r chi.Router) {
		r.Use(g.m.RequireActivatedUser)

		r.Get("/{id}", g.handler.GetByGoalID)
		r.Post("/", g.handler.Create)
		r.Put("/{id}", g.handler.Update)
		r.Delete("/{id}", g.handler.Delete)
	})
}
