package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) *GormProductRepository {
	return &GormProductRepository{
		db: db,
	}
}

func (r *GormProductRepository) ListPaged(limit, offset int) ([]models.Product, int64, error) {
	// Rebuilt per call: reusing one *gorm.DB for Count and Find leaks the
	// count's conditions into the select.
	query := func() *gorm.DB {
		return r.db.Model(&models.Product{})
	}

	var total int64
	if err := query().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.Product
	err := query().
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ListLowStockPaged returns products at or below their minimum stock
// (PRD section 6.3).
func (r *GormProductRepository) ListLowStockPaged(limit, offset int) ([]models.Product, int64, error) {
	query := func() *gorm.DB {
		return r.db.Model(&models.Product{}).Where("stock <= minimum_stock")
	}

	var total int64
	if err := query().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.Product
	err := query().
		Order("stock ASC, id ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *GormProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product

	err := r.db.
		First(&product, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	return &product, nil
}

func (r *GormProductRepository) GetBySKU(sku string) (*models.Product, error) {
	var product models.Product

	err := r.db.
		Where("sku = ?", sku).
		First(&product).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	return &product, nil
}

func (r *GormProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *GormProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *GormProductRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Product{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}
