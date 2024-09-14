package payloads

type UserRegister struct {
	Username string `json:"username" validate:"required,gte=2,lte=16"`
	Password string `json:"password" validate:"required,gte=12,lte=32"`
	Email    string `json:"email" validate:"required,email"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserUpdate struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,gte=12,lte=32"`
	//Role     *enums.UserRole `json:"role,omitempty" validate:"omitempty,gte=0,lte=6"`
}
