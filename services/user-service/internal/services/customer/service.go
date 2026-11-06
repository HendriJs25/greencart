package customer

import (
	"context"
	customerrepository "user-service/internal/repository/customer"
)

type service struct {
	customerRepository customerrepository.Repository
}

type Service interface {
	GetAll(ctx context.Context) ([]GetAllCustomerResult, error)
}

func NewService(customerRepository customerrepository.Repository) Service {
	return &service{
		customerRepository: customerRepository,
	}
}

func (s *service) GetAll(ctx context.Context) ([]GetAllCustomerResult, error) {
	customers, err := s.customerRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := ToGetAllCustomerResult(customers)
	return result, nil
}
