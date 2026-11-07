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
	GetDetail(ctx context.Context, id int64) (*GetDetailCustomerResult, error)
	Delete(ctx context.Context, id int64) error
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

func (s *service) GetDetail(ctx context.Context, id int64) (*GetDetailCustomerResult, error) {
	customer, err := s.customerRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := ToGetDetailCustomerResult(customer)
	return result, nil
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.customerRepository.DeleteByID(ctx, id)
}
