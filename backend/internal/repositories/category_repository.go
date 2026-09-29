package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
)

var ErrCategoryNotFound = errors.New("category not found")

// CategoryRepository abstracts category persistence.
type CategoryRepository interface {
	ListPaged(limit, offset int) ([]models.Category, int64, error)
	GetByID(id uint) (*models.Category, error)
	Create(category *models.Category) error
	Update(category *models.Category) error
	Delete(id uint) error
}
