package api

import (
	"context"
	"financas/internal/model"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) ContextSetUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) ContextGetUser(r *http.Request) *model.User {
	user, ok := r.Context().Value(userContextKey).(*model.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
