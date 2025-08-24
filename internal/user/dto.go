package user

type (
	CreateUserParam struct {
		Username *string `json:"username"`
		Email    string  `json:"email"`
		Password string  `json:"password"`
	}

	UpdateUserParam struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}
)
