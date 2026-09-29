package services

import (
	"net/http"
	"strings"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

// CategoryInput defines the accepted payload for create and update.
type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}

// CategoryService defines category business logic.
type CategoryService interface {
	List(page, limit int) ([]models.Category, int64, error)
	GetByID(id uint) (*models.Category, error)
	Create(input *CategoryInput) (*models.Category, error)
	Update(id uint, input *CategoryInput) (*models.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) List(page, limit int) ([]models.Category, int64, error) {
	return s.repo.ListPaged(limit, (page-1)*limit)
}

func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			return nil, NewAppError("CATEGORY_NOT_FOUND", "category not found", http.StatusNotFound)
		}
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Create(input *CategoryInput) (*models.Category, error) {
	if err := validateCategoryInput(input); err != nil {
		return nil, err
	}

	category := &models.Category{
		Name: strings.TrimSpace(input.Name),
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Update(id uint, input *CategoryInput) (*models.Category, error) {
	if err := validateCategoryInput(input); err != nil {
		return nil, err
	}

	category, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			return nil, NewAppError("CATEGORY_NOT_FOUND", "category not found", http.StatusNotFound)
		}
		return nil, err
	}

	category.Name = strings.TrimSpace(input.Name)
	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrCategoryNotFound {
			return NewAppError("CATEGORY_NOT_FOUND", "category not found", http.StatusNotFound)
		}
		return err
	}

	return s.repo.Delete(id)
}

func validateCategoryInput(input *CategoryInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return NewAppError("INVALID_INPUT", "category name is required", http.StatusUnprocessableEntity)
	}
	return nil
}
