package customer

import (
	"context"
	queryinput "user-service/internal/domain/query"
	"user-service/internal/query"
	customerrepository "user-service/internal/repository/customer"
)

type service struct {
	customerRepository customerrepository.Repository
}

type Service interface {
	GetAll(ctx context.Context, params queryinput.GetAllCustomerQuery) ([]GetAllCustomerResult, *query.Meta, error)
}

func NewService(customerRepository customerrepository.Repository) Service {
	return &service{
		customerRepository: customerRepository,
	}
}

func (s *service) GetAll(ctx context.Context, params queryinput.GetAllCustomerQuery) ([]GetAllCustomerResult, *query.Meta, error) {
	customers, total, err := s.customerRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	metadata := query.Metadata(params.Page, params.Size, total)
	result := ToGetAllCustomerResult(customers)
	return result, metadata, nil
}
