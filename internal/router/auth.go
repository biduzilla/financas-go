package router

import (
	"financas/internal/handler"

	"github.com/go-chi/chi"
)

type AuthRouter struct {
	Auth handler.AuthHandlerInterface
}

type AuthRoutesInterface interface {
	AuthRoutes(r chi.Router)
}

func NewAuthRouter(authHandler handler.AuthHandlerInterface) *AuthRouter {
	return &AuthRouter{
		Auth: authHandler,
	}
}

func (a *AuthRouter) AuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", a.Auth.LoginHandler)
	})
}
