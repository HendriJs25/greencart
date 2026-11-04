package request

import authservice "user-service/internal/services/auth"

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignUpRequest struct {
	Name                 string `json:"name" validate:"required,notblank"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,validpassword"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func ToSignInInput(req SignInRequest) authservice.SignInInput {
	return authservice.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

func ToSignUpInput(req SignUpRequest) authservice.SignUpInput {
	return authservice.SignUpInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}
