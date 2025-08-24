package auth

import "github.com/google/uuid"

type (
	RegisterParam struct {
		Username *string `json:"username" validate:"omitempty"`
		Email    string  `json:"email" validate:"required,email"`
		Password string  `json:"password" validate:"required,min=6"`
	}

	LoginParam struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	LoginResponse struct {
		Tokens
		User UserClaims `json:"user"`
	}

	UserClaims struct {
		UserID   uuid.UUID `json:"user_id"`
		Email    *string   `json:"email,omitempty"`
		Username *string   `json:"username,omitempty"`
	}

	ChangePasswordParam struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=6"`
	}
)
