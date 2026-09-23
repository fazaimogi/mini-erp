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

func (r *GormProductRepository) List() ([]models.Product, error) {
	var products []models.Product

	err := r.db.
		Order("id ASC").
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	return products, nil
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
