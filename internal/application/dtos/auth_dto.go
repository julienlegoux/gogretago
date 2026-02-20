package dtos

// RegisterInput contains the data needed for user registration
type RegisterInput struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,password"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

// LoginInput contains the data needed for user login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

// AuthResponse is returned after successful authentication
type AuthResponse struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}
