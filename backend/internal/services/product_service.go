package services

import (
	"strings"

	"github.com/shopspring/decimal"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

const (
	maxNameLength = 255
	maxSKULength  = 100
)

// ProductInput defines the accepted payload for create and update operations.
// PRD section 7.4: name, sku, category, price and minimum stock are required,
// sku must be unique and price must not be negative.
type ProductInput struct {
	Name         string          `json:"name" binding:"required"`
	SKU          string          `json:"sku" binding:"required"`
	CategoryID   uint            `json:"category_id" binding:"required,gt=0"`
	Description  string          `json:"description"`
	Price        decimal.Decimal `json:"price" binding:"required"`
	Stock        int             `json:"stock" binding:"min=0"`
	MinimumStock int             `json:"minimum_stock" binding:"min=0"`
}

type ProductService interface {
	List(page, limit int) ([]models.Product, int64, error)
	ListLowStock(page, limit int) ([]models.Product, int64, error)
	GetByID(id uint) (*models.Product, error)
	Create(input *ProductInput) (*models.Product, error)
	Update(id uint, input *ProductInput) (*models.Product, error)
	Delete(id uint) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) List(page, limit int) ([]models.Product, int64, error) {
	return s.repo.ListPaged(limit, (page-1)*limit)
}

// ListLowStock backs the dashboard's low stock metric and the product report
// (PRD section 6.3): stock <= minimum_stock.
func (s *productService) ListLowStock(page, limit int) ([]models.Product, int64, error) {
	return s.repo.ListLowStockPaged(limit, (page-1)*limit)
}

func (s *productService) GetByID(id uint) (*models.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrProductNotFound {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *productService) Create(input *ProductInput) (*models.Product, error) {
	if err := validateProductInput(input); err != nil {
		return nil, err
	}
	if err := s.ensureUniqueSKU(input.SKU, 0); err != nil {
		return nil, err
	}

	product := &models.Product{
		Name:         strings.TrimSpace(input.Name),
		SKU:          strings.ToUpper(strings.TrimSpace(input.SKU)),
		CategoryID:   input.CategoryID,
		Description:  input.Description,
		Price:        input.Price,
		Stock:        input.Stock,
		MinimumStock: input.MinimumStock,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Update(id uint, input *ProductInput) (*models.Product, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrProductNotFound {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if err := validateProductInput(input); err != nil {
		return nil, err
	}
	if err := s.ensureUniqueSKU(input.SKU, id); err != nil {
		return nil, err
	}

	product := &models.Product{
		ID:           id,
		Name:         strings.TrimSpace(input.Name),
		SKU:          strings.ToUpper(strings.TrimSpace(input.SKU)),
		CategoryID:   input.CategoryID,
		Description:  input.Description,
		Price:        input.Price,
		Stock:        input.Stock,
		MinimumStock: input.MinimumStock,
		// Carried over so the full-column Save keeps the original creation time.
		CreatedAt: existing.CreatedAt,
	}

	if err := s.repo.Update(product); err != nil {
		if err == repositories.ErrProductNotFound {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (s *productService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		if err == repositories.ErrProductNotFound {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

func (s *productService) ensureUniqueSKU(sku string, excludeID uint) error {
	existing, err := s.repo.GetBySKU(strings.ToUpper(strings.TrimSpace(sku)))
	if err != nil {
		if err == repositories.ErrProductNotFound {
			return nil
		}
		return err
	}
	if existing.ID == excludeID {
		return nil
	}
	return ErrDuplicateSKU
}

// validateProductInput performs business validation that gin binding tags
// cannot express (PRD section 7.4).
func validateProductInput(input *ProductInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" || len(name) > maxNameLength {
		return ErrInvalidInput
	}

	sku := strings.TrimSpace(input.SKU)
	if sku == "" || len(sku) > maxSKULength {
		return ErrInvalidInput
	}

	if input.CategoryID == 0 || input.Price.IsNegative() || input.Stock < 0 || input.MinimumStock < 0 {
		return ErrInvalidInput
	}

	return nil
}
