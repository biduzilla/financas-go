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

type UserRoutesInterface interface {
	UserRoutes(r chi.Router)
}

func (u *UserRouter) UserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/activate", u.User.ActivateUserHandler)
		r.Post("/", u.User.CreateUserHandler)
	})
}
