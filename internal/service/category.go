package service

import (
	"financas/internal/model"
	"financas/internal/model/filters"
	"financas/internal/repository"
	e "financas/utils/errors"
	"financas/utils/validator"
)

type CategoryService struct {
	CategoryRepository repository.CategoryRepositoryIntercafe
}

type CategoryServiceInterface interface {
	GetByID(id int64, userID int64) (*model.Category, error)
	GetAll(name string, userID int64, f filters.Filters, v *validator.Validator) ([]*model.Category, filters.Metadata, error)
	Insert(category *model.Category, v *validator.Validator) error
	Update(category *model.Category, userID int64, v *validator.Validator) error
	Delete(id int64, userID int64) error
}

func NewCategoryService(c repository.CategoryRepositoryIntercafe) *CategoryService {
	return &CategoryService{
		CategoryRepository: c,
	}
}

func (s *CategoryService) GetByID(id int64, userID int64) (*model.Category, error) {
	user, err := s.CategoryRepository.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *CategoryService) GetAll(name string, userID int64, f filters.Filters, v *validator.Validator) ([]*model.Category, filters.Metadata, error) {
	if filters.ValidateFilters(v, f); !v.Valid() {
		return nil, filters.Metadata{}, e.ErrInvalidData
	}

	categories, metadata, err := s.CategoryRepository.GetAll(
		name,
		userID,
		f,
	)

	if err != nil {
		return nil, filters.Metadata{}, err
	}

	return categories, metadata, nil
}

func (s *CategoryService) Insert(category *model.Category, v *validator.Validator) error {
	if category.ValidateCategory(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.CategoryRepository.Insert(category)

	if err != nil {
		return err
	}

	return nil

}

func (s *CategoryService) Update(category *model.Category, userID int64, v *validator.Validator) error {
	if category.ValidateCategory(v); !v.Valid() {
		return e.ErrInvalidData
	}

	err := s.CategoryRepository.Update(category, userID)

	if err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) Delete(id int64, userID int64) error {
	err := s.CategoryRepository.Delete(id, userID)

	if err != nil {
		return err
	}

	return nil
}
