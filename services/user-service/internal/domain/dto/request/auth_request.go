package request

import authservice "user-service/internal/services/auth"

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func ToSignInInput(req SignInRequest) authservice.SignInInput {
	return authservice.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	}
}
