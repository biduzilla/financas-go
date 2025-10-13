package router

import (
	"financas/internal/handler"

	"github.com/go-chi/chi"
)

type UserRouter struct {
	User handler.UserHandlerInterface
}

func NewUserRouter(userHandler handler.UserHandlerInterface) *UserRouter {
	return &UserRouter{
		User: userHandler,
	}
}

func (u *UserRouter) UserRoutes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/activate", u.User.ActivateUserHandler)
		r.Post("/", u.User.CreateUserHandler)
	})
}
