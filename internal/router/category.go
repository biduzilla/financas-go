package router

import (
	"financas/internal/handler"
	"financas/internal/middleware"

	"github.com/go-chi/chi"
)

type CategoryRouter struct {
	categoryHandler handler.CategoryHandlerInterface
	middleware      middleware.MiddlewareInterface
}

func NewCategoryRouter(h handler.CategoryHandlerInterface, middleware middleware.MiddlewareInterface) *CategoryRouter {
	return &CategoryRouter{
		categoryHandler: h,
		middleware:      middleware,
	}
}

type CategoryRouterIntercace interface {
	CategoryRoutes(r chi.Router)
}

func (c *CategoryRouter) CategoryRoutes(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Use(c.middleware.RequireActivatedUser)

		r.Get("/{id}", c.categoryHandler.GetById)
		r.Get("/", c.categoryHandler.GetAll)
		r.Post("/", c.categoryHandler.Insert)
		r.Put("/{id}", c.categoryHandler.Update)
		r.Delete("/{id}", c.categoryHandler.Delete)
	})
}
