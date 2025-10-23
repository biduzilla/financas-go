package router

import (
	"financas/internal/handler"
	"financas/internal/middleware"

	"github.com/go-chi/chi"
)

type TransactionRouter struct {
	transaction handler.TransactionHandlerInterface
	m           middleware.MiddlewareInterface
}

type TransactionRouterInterface interface {
	TransactionRoutes(r chi.Router)
}

func NewTransactionRouter(
	transaction handler.TransactionHandlerInterface,
	m middleware.MiddlewareInterface,
) *TransactionRouter {
	return &TransactionRouter{
		transaction: transaction,
		m:           m,
	}
}

func (router *TransactionRouter) TransactionRoutes(r chi.Router) {
	r.Route("/transactions", func(r chi.Router) {
		r.Use(router.m.RequireActivatedUser)

		r.Get("/{id}", router.transaction.GetByID)
		r.Get("/category/{id}", router.transaction.GetAllByUserAndCategory)
		r.Post("/", router.transaction.Save)
		r.Put("/", router.transaction.Update)
		r.Delete("/{id}", router.transaction.DeleteByID)
	})
}
