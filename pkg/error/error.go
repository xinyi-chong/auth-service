package apperrors

import (
	locale "auth-service/pkg/i18n"
	"errors"
	"net/http"
)

var (
	ErrUnauthorized             = New("unauthorized", http.StatusUnauthorized)
	ErrRequestFailed            = New("request_failed", http.StatusBadGateway)
	ErrSessionExpired           = New("session_expired", http.StatusUnauthorized)
	ErrBadRequest               = New("bad_request", http.StatusBadRequest)
	ErrInternalServerError      = New("internal_server_error", http.StatusInternalServerError)
	ErrTooManyRequests          = New("too_many_requests", http.StatusTooManyRequests)
	ErrRegistrationFailed       = New("registration_failed", http.StatusInternalServerError)
	ErrInvalidEmail             = New("invalid_email", http.StatusBadRequest)
	ErrIncorrectEmailOrPassword = New("incorrect_email_or_password", http.StatusUnauthorized)
	ErrIncorrectPassword        = New("incorrect_password", http.StatusUnauthorized)
	ErrEmailConflict            = New("email_already_exists", http.StatusConflict)
	ErrUserNotFound             = New("user_not_found", http.StatusNotFound)
	ErrUserConflict             = New("user_already_exists", http.StatusConflict)
	ErrUsernameConflict         = New("username_already_exists", http.StatusConflict)
)

type Error struct {
	MessageKey   string              // i18n key (e.g., "user_not_found")
	HTTPStatus   int                 // HTTP status code
	Err          error               // Wrapped error, if any
	TemplateData locale.TemplateData // For dynamic i18n translation
	Op           string              // Operation name for logging or debugging
}

func New(msgKey string, status int) *Error {
	return &Error{MessageKey: msgKey, HTTPStatus: status} //, Err: errors.New(code)
}

func Is(err error, target *Error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.MessageKey == target.MessageKey
	}
	return false
}

func (e *Error) Error() string {
	var opPrefix string
	if e.Op != "" {
		opPrefix = "[" + e.Op + "] "
	}

	errStr := e.MessageKey

	if e.Err != nil {
		var appErr *Error
		if !errors.As(e.Err, &appErr) {
			errStr = e.Err.Error()
		}
	}

	return opPrefix + errStr
}

func (e *Error) WithOp(op string) *Error {
	e.Op = op
	return e
}

func (e *Error) Wrap(err error) *Error {
	e.Err = err
	return e
}
