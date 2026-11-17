package category

import (
	"context"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/model"
	categoryrepository "product-service/internal/repository/category"
)

type service struct {
	categoryRepository categoryrepository.Repository
}

type Service interface {
	GetAll(ctx context.Context) ([]GetAllCategoriesResult, error)
	Create(ctx context.Context, input CreateCategoryInput) error
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

func (s *service) Create(ctx context.Context, input CreateCategoryInput) error {
	exist, err := s.categoryRepository.IsSlugExist(ctx, input.Slug)
	if err != nil {
		return err
	}
	if exist {
		return errConstant.ErrCategoryAlreadyExists
	}

	var parentCategory *model.Category
	if input.ParentSlug != nil {
		parentCategory, err = s.categoryRepository.FindBySlug(ctx, input.ParentSlug)
		if err != nil {
			return err
		}
	}

	input := model.Category{
		ParentID: parentCategory.ID,
		Name:        input.Name,
		Icon:        input.Icon,
		Status:      input.Status.Bool(),
		Slug:        input.Slug,
		Description: input.Description,
	}

	err = s.categoryRepository.Create(ctx, ToCreateCategoryModel(input))
	if err != nil {
		return err
	}
	return nil
}
