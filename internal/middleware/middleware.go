package middleware

import (
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils/errors"
	"financas/utils/validator"
	"net/http"
	"strings"
)

type Middleware struct {
	ErrResp        errors.ErrorResponseInterface
	ContextGetUser func(r *http.Request) *model.User
	ContextSetUser func(r *http.Request, user *model.User) *http.Request
	AuthService    service.AuthServiceInterface
	UserService    service.UserServiceInterface
}

func New(
	errResp errors.ErrorResponseInterface,
	contextGetUser func(r *http.Request) *model.User,
	contextSetUser func(r *http.Request, user *model.User) *http.Request,
	authService service.AuthServiceInterface,
	userService service.UserServiceInterface,
) *Middleware {
	return &Middleware{
		ErrResp:        errResp,
		ContextGetUser: contextGetUser,
		ContextSetUser: contextSetUser,
		AuthService:    authService,
		UserService:    userService,
	}
}

// RequireAuthenticatedUser garante que o usuário não é anônimo
func (m *Middleware) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.ContextGetUser(r)
		if user.IsAnonymous() {
			m.ErrResp.AuthenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireActivatedUser garante que o usuário está ativado e autenticado
func (m *Middleware) RequireActivatedUser(next http.Handler) http.Handler {
	return m.RequireAuthenticatedUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.ContextGetUser(r)
		if !user.Activated {
			m.ErrResp.InactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// Authenticate extrai o token do header e define o usuário no contexto
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = m.ContextSetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.ErrResp.InvalidCredentialsResponse(w, r)
			return
		}

		token := headerParts[1]
		username, err := m.AuthService.ExtractUsername(token)
		if err != nil {
			m.ErrResp.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		v := validator.New()
		user, err := m.UserService.GetUserByEmail(username, v)
		if err != nil {
			m.ErrResp.HandlerErrorResponse(w, r, err, v)
			return
		}

		r = m.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
