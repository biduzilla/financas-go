package middleware

import (
	"expvar"
	"financas/internal/config"
	"financas/internal/model"
	"financas/internal/service"
	"financas/utils/errors"
	"financas/utils/validator"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	totalRequestsReceived           = expvar.NewInt("total_requests_received")
	totalResponsesSent              = expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
	totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")
)

type Middleware struct {
	ErrResp        errors.ErrorResponseInterface
	ContextGetUser func(r *http.Request) *model.User
	ContextSetUser func(r *http.Request, user *model.User) *http.Request
	AuthService    service.AuthServiceInterface
	UserService    service.UserServiceInterface
	Config         config.Config
}

func New(
	errResp errors.ErrorResponseInterface,
	contextGetUser func(r *http.Request) *model.User,
	contextSetUser func(r *http.Request, user *model.User) *http.Request,
	authService service.AuthServiceInterface,
	userService service.UserServiceInterface,
	Config config.Config,
) *Middleware {
	return &Middleware{
		ErrResp:        errResp,
		ContextGetUser: contextGetUser,
		ContextSetUser: contextSetUser,
		AuthService:    authService,
		UserService:    userService,
		Config:         Config,
	}
}

type metricsResponseWriter struct {
	wrapped       http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}

func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)
	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.wrapped.Write(b)
}

func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}

func (m *Middleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		totalRequestsReceived.Add(1)

		mw := newMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		totalResponsesSent.Add(1)
		totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)

		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}

func (m *Middleware) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")
		if origin != "" {
			for i := range m.Config.CORS.TrustedOrigins {
				if origin == m.Config.CORS.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

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

func (m *Middleware) RateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					if time.Since(client.lastSeen) > 3*time.Minute {
						delete(clients, ip)
					}
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.Config.Limiter.Enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				m.ErrResp.ServerErrorResponse(w, r, err)
				return
			}
			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(m.Config.Limiter.RPS), m.Config.Limiter.Burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				m.ErrResp.RateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				w.Header().Set("Connection", "close")

				m.ErrResp.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
