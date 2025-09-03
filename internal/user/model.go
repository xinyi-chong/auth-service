package user

import (
	"github.com/google/uuid"
	"github.com/xinyi-chong/common-lib/filters"
	"time"
)

type User struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Username           *string    `json:"username,omitempty" db:"username"`
	Email              *string    `json:"email,omitempty" db:"email"`
	EmailVerified      bool       `json:"email_verified" db:"email_verified"`
	PasswordHash       *string    `json:"-" db:"password_hash"` // never expose in JSON
	LastLogin          *time.Time `json:"last_login,omitempty" db:"last_login"`
	IsActive           bool       `json:"is_active" db:"is_active"`
	PasswordChangedAt  *time.Time `json:"password_changed_at,omitempty" db:"password_changed_at"`
	AccountLockedUntil *time.Time `json:"account_locked_until,omitempty" db:"account_locked_until"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type Response struct {
	ID                 uuid.UUID  `json:"id"`
	Username           *string    `json:"username,omitempty"`
	Email              *string    `json:"email,omitempty"`
	EmailVerified      bool       `json:"email_verified"`
	IsActive           bool       `json:"is_active"`
	LastLogin          *time.Time `json:"last_login,omitempty"`
	AccountLockedUntil *time.Time `json:"account_locked_until,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type Filter struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	filters.Pagination
}

func (u *User) Response() *Response {
	return &Response{
		ID:                 u.ID,
		Username:           u.Username,
		Email:              u.Email,
		EmailVerified:      u.EmailVerified,
		IsActive:           u.IsActive,
		LastLogin:          u.LastLogin,
		AccountLockedUntil: u.AccountLockedUntil,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
	}
}
