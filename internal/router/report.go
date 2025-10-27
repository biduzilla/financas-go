package router

import (
	"financas/internal/handler"
	"financas/internal/middleware"

	"github.com/go-chi/chi"
)

type ReportRouter struct {
	handler handler.ReportHandlerInterface
	m       middleware.MiddlewareInterface
}

func NewReportRouter(h handler.ReportHandlerInterface, m middleware.MiddlewareInterface) ReportRouterInterface {
	return &ReportRouter{
		handler: h,
		m:       m,
	}
}

type ReportRouterInterface interface {
	ReportRoutes(r chi.Router)
}

func (router *ReportRouter) ReportRoutes(r chi.Router) {
	r.Route("/reports", func(r chi.Router) {
		r.Use(router.m.RequireActivatedUser)

		r.Get("/summary", router.handler.GetFinancialSummaryHandler)
		r.Get("/categories", router.handler.GetCategoryReportHandler)
		r.Get("/top-categories", router.handler.GetTopCategoriesHandler)
	})
}
