package success

import (
	locale "auth-service/pkg/i18n"
	"net/http"
)

var (
	SessionRefreshed = New("session_refreshed", "Session successfully refreshed", http.StatusOK)
	Registered       = New("registered", "User successfully registered", http.StatusCreated)
	LoggedIn         = New("logged_in", "User successfully logged in", http.StatusOK)
	LoggedOut        = New("logged_out", "User successfully logged out", http.StatusOK)
	PasswordChanged  = New("password_changed", "Password successfully changed", http.StatusOK)
	PasswordReset    = New("password_reset", "Password successfully reset", http.StatusOK)
	UserUpdated      = New("user_updated", "User information successfully updated", http.StatusOK)
	UserDeleted      = New("user_deleted", "User successfully deleted", http.StatusOK)
	UserFound        = New("user_found", "User successfully found", http.StatusFound)
)

type Success struct {
	Code         string
	DefaultMsg   string
	HTTPStatus   int
	TemplateData locale.TemplateData
}

func New(code, defaultMsg string, httpStatus int) *Success {
	return &Success{Code: code, DefaultMsg: defaultMsg, HTTPStatus: httpStatus}
}
