package auth

import (
	"context"
	"errors"
	"user-service/internal/common/helper"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/logger"
	userrepository "user-service/internal/repository/user"
	jwtservice "user-service/internal/services/jwt"
)

type service struct {
	jwtService     jwtservice.Service
	userRepository userrepository.Repository
}

type Service interface {
	SignIn(ctx context.Context, input SignInInput) (*SignInResult, error)
}

func NewService(jwtService jwtservice.Service, userRepository userrepository.Repository) Service {
	return &service{
		jwtService:     jwtService,
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

	csrfToken, err := helper.GenerateCSRFToken()
	if err != nil {
		return nil, err
	}
	logger.Log.Debug("csrftoken")

	accessToken, err := s.jwtService.GenerateAccessToken(user.UUID)
	if err != nil {
		return nil, err
	}
	logger.Log.Debug("accesstoken")

	return ToSignInResult(user, accessToken, csrfToken), nil
}
