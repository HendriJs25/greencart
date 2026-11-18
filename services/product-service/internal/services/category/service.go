package category

import (
	"context"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/enum"
	"product-service/internal/domain/model"
	categoryrepository "product-service/internal/repository/category"
)

type service struct {
	categoryRepository categoryrepository.Repository
}

type Service interface {
	GetAll(ctx context.Context) ([]GetAllCategoriesResult, error)
	GetByID(ctx context.Context, id int64) (*GetDetailCategoryResult, error)
	GetAllSlugForOptions(ctx context.Context) ([]string, error)
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

func (s *service) GetByID(ctx context.Context, id int64) (*GetDetailCategoryResult, error) {
	category, err := s.categoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var parentSlug *string
	if category.Parent != nil {
		parentSlug = &category.Parent.Slug
	}

	result := &GetDetailCategoryResult{
		ID:          category.ID,
		ParentSlug:  parentSlug,
		Name:        category.Name,
		Icon:        category.Icon,
		Status:      enum.CategoryStatusFromBool(category.Status),
		Slug:        category.Slug,
		Description: category.Description,
	}

	return result, nil
}

func (s *service) GetAllSlugForOptions(ctx context.Context) ([]string, error) {
	return s.categoryRepository.FindAllSlugForOptions(ctx)
}

func (s *service) Create(ctx context.Context, input CreateCategoryInput) error {
	exist, err := s.categoryRepository.IsSlugExist(ctx, input.Slug)
	if err != nil {
		return err
	}
	if exist {
		return errConstant.ErrCategoryAlreadyExists
	}

	var parentCategoryID *int64
	if input.ParentSlug != nil {
		parentCategory, err := s.categoryRepository.FindBySlug(ctx, input.ParentSlug)
		if err != nil {
			return err
		}
		parentCategoryID = &parentCategory.ID
	}

	inputModel := &model.Category{
		ParentID:    parentCategoryID,
		Name:        input.Name,
		Icon:        input.Icon,
		Status:      input.Status.Bool(),
		Slug:        input.Slug,
		Description: input.Description,
	}

	err = s.categoryRepository.Create(ctx, inputModel)
	if err != nil {
		return err
	}
	return nil
}
