package errors

import (
	"financas/internal/jsonlog"
	"financas/utils"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	logger *jsonlog.Logger
}

func NewErrorResponse(logger *jsonlog.Logger) *ErrorResponse {
	return &ErrorResponse{logger: logger}
}

func (e *ErrorResponse) NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	e.errorResponse(w, r, http.StatusForbidden, message)
}

func (e *ErrorResponse) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (e *ErrorResponse) InactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	e.errorResponse(w, r, http.StatusForbidden, message)
}

func (e *ErrorResponse) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *ErrorResponse) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	e.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (e *ErrorResponse) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceed"
	e.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (e *ErrorResponse) logError(r *http.Request, err error) {
	e.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (e *ErrorResponse) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := utils.Envelope{"error": message}
	err := utils.WriteJSON(w, status, env, nil)
	if err != nil {
		e.logError(r, err)
		w.WriteHeader(500)
	}
}

func (e *ErrorResponse) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	e.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (e *ErrorResponse) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	e.errorResponse(w, r, http.StatusNotFound, message)
}

func (e *ErrorResponse) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	e.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (e *ErrorResponse) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	e.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (e *ErrorResponse) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	e.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (e *ErrorResponse) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	e.errorResponse(w, r, http.StatusConflict, message)
}
