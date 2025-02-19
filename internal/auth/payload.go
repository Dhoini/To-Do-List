package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,min=10"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email,min=10"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
