package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
	"gorm.io/gorm"
)

var ErrSaleNotFound = errors.New("sale not found")

type SaleRepository interface {
	GetByID(id uint) (*models.Sale, error)
	ListPaged(limit, offset int) ([]models.Sale, int64, error)
	Create(sale *models.Sale) error
	CreateWithItems(sale *models.Sale, items []models.SaleItem) error
}

type GormSaleRepository struct {
	db *gorm.DB
}

func NewGormSaleRepository(db *gorm.DB) *GormSaleRepository {
	return &GormSaleRepository{db: db}
}

func (r *GormSaleRepository) GetByID(id uint) (*models.Sale, error) {
	var sale models.Sale
	if err := r.db.Preload("Items").First(&sale, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSaleNotFound
		}
		return nil, err
	}
	return &sale, nil
}

func (r *GormSaleRepository) ListPaged(limit, offset int) ([]models.Sale, int64, error) {
	query := func() *gorm.DB {
		return r.db.Model(&models.Sale{})
	}

	var total int64
	if err := query().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var sales []models.Sale
	if err := query().
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("Items").
		Find(&sales).
		Error; err != nil {
		return nil, 0, err
	}

	return sales, total, nil
}

func (r *GormSaleRepository) Create(sale *models.Sale) error {
	return r.db.Create(sale).Error
}

func (r *GormSaleRepository) CreateWithItems(sale *models.Sale, items []models.SaleItem) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(sale).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range items {
		items[i].SaleID = sale.ID
		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
