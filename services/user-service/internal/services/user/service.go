package user

import (
	"context"
	userrepository "user-service/internal/repository/user"
	"uuid"
)

type service struct {
	userRepository userrepository.Repository
}

type Service interface {
	GetProfile(ctx context.Context, userUUID uuid.UUID) (*GetProfileResult, error)
}

func NewService(userRepository userrepository.Repository) Service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) GetProfile(ctx context.Context, userUUID uuid.UUID) (*GetProfileResult, error) {
	user, err := s.userRepository.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	result := ToGetProfileResult(user)
	return result, nil
}
