package auth

import (
	"context"
	"errors"
	"time"
	"user-service/internal/common/helper"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/model"
	sessionrepository "user-service/internal/repository/session"
	userrepository "user-service/internal/repository/user"
	jwtservice "user-service/internal/services/jwt"
	"uuid"
)

type service struct {
	jwtService        jwtservice.Service
	userRepository    userrepository.Repository
	sessionRepository sessionrepository.Repository
}

type Service interface {
	SignIn(ctx context.Context, input SignInInput) (*SignInResult, error)
}

func NewService(jwtService jwtservice.Service, userRepository userrepository.Repository, sessionRepository sessionrepository.Repository) Service {
	return &service{
		jwtService:        jwtService,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
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

	hashedCSRFToken := helper.HashToken(csrfToken)

	jti := uuid.NewV4()

	accessToken, err := s.jwtService.GenerateAccessToken(user.UUID, jti)
	if err != nil {
		return nil, err
	}

	var userRoles []int64
	for _, role := range user.Roles {
		userRoles = append(userRoles, role.ID)
	}

	sessionCreatedAt := time.Now().UTC()

	sessionTTL := accessToken.ExpiresAt.Sub(sessionCreatedAt)

	session := model.Session{
		UUID:          user.UUID,
		CSRFTokenHash: hashedCSRFToken,
		Email:         user.Email,
		RoleID:        userRoles,
		CreatedAt:     sessionCreatedAt,
	}

	err = s.sessionRepository.Set(ctx, jti.String(), session, sessionTTL)
	if err != nil {
		return nil, err
	}

	return ToSignInResult(user, accessToken, csrfToken), nil
}
