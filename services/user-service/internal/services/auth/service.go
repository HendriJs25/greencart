package auth

import (
	"context"
	"errors"
	"user-service/internal/common/helper"
	errConstant "user-service/internal/constants/error"
	userrepository "user-service/internal/repository/user"
)

type service struct {
	userRepository userrepository.Repository
}

type Service interface {
	SignIn(ctx context.Context, input SignInInput) (*SignInResult, error)
}

func NewService(userRepository userrepository.Repository) Service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) SignIn(ctx context.Context, input SignInInput) (*SignInResult, error) {
	normalizedEmail := helper.NormalizeEmail(input.Email)
	user, err := s.userRepository.FindByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, errConstant.ErrNotFound) {
			return nil, errConstant.ErrUserNotFound
		}
		return nil, err
	}

	isVerify := helper.VerifyPassword(input.Password, user.PasswordHash)
	if !isVerify {
		return nil, errConstant.ErrPasswordIncorrect
	}

	CSRFToken, err := helper.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}

	return ToSignInResult(user, CSRFToken), nil
}
