package apperrors

import (
	locale "auth-service/pkg/i18n"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized             = New("unauthorized", "Unauthorized access", http.StatusUnauthorized)
	ErrRequestFailed            = New("request_failed", "Request failed", http.StatusBadGateway)
	ErrSessionExpired           = New("session_expired", "Session has expired", http.StatusUnauthorized)
	ErrBadRequest               = New("bad_request", "Bad request", http.StatusBadRequest)
	ErrInternalServerError      = New("internal_server_error", "Something went wrong", http.StatusInternalServerError)
	ErrTooManyRequests          = New("too_many_requests", "Too many requests, please try again later", http.StatusTooManyRequests)
	ErrRegistrationFailed       = New("registration_failed", "Registration failed", http.StatusInternalServerError)
	ErrInvalidEmail             = New("invalid_email", "Invalid email format", http.StatusBadRequest)
	ErrIncorrectEmailOrPassword = New("incorrect_email_or_password", "Incorrect email or password", http.StatusUnauthorized)
	ErrIncorrectPassword        = New("incorrect_password", "Incorrect password", http.StatusUnauthorized)
	ErrEmailExists              = New("email_already_exists", "Email already registered", http.StatusConflict)
	ErrUserNotFound             = New("user_not_found", "User not found", http.StatusNotFound)
	ErrUsernameExists           = New("username_already_exists", "Username already exists", http.StatusConflict)
)

type Error struct {
	Code         string              // i18n key (e.g., "user_not_found")
	DefaultMsg   string              // Fallback
	HTTPStatus   int                 // HTTP status code
	Err          error               // Wrapped error, if any
	TemplateData locale.TemplateData // For dynamic i18n translation
	Tags         []string
}

func New(code, defaultMsg string, status int) *Error {
	return &Error{Code: code, DefaultMsg: defaultMsg, HTTPStatus: status, Err: errors.New(code)}
}

func Is(err error, target *Error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == target.Code
	}
	return false
}

func (e *Error) Error() string {
	base := e.Code
	if e.Err != nil {
		base = e.Err.Error()
	}
	if len(e.Tags) > 0 {
		return "[" + strings.Join(e.Tags, "-") + "] " + base
	}
	return base
}

func (e *Error) WithTags(tags ...string) *Error {
	if len(tags) > 0 {
		e.Tags = append(e.Tags, tags...)
	}
	return e
}

func (e *Error) Wrap(err error) *Error {
	e.Err = err
	return e
}

func IsDbErrNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func HandleDbNotFoundError(err error, notFoundErr *Error) *Error {
	if IsDbErrNotFound(err) {
		return notFoundErr
	}
	return ErrInternalServerError.Wrap(err)
}
