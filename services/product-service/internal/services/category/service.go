package category

import (
	"context"
	categoryrepository "product-service/internal/repository/category"
)

type service struct {
	categoryRepository categoryrepository.Repository
}

type Service interface {
	GetAll(ctx context.Context) ([]GetAllCategoriesResult, error)
}

func NewService(categoryRepository categoryrepository.Repository) Service {
	return &service{
		categoryRepository: categoryRepository,
	}
}

func (s *service) GetAll(ctx context.Context) ([]GetAllCategoriesResult, error) {
	categories, err := s.categoryRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := ToGetAllCategoriesResult(categories)
	return result, nil
}
