package repositories

import (
	"github.com/fazasuny/erp-system/internal/models"
	"gorm.io/gorm"
)

type gormCategoryRepository struct {
	db *gorm.DB
}

func NewGormCategoryRepository(db *gorm.DB) CategoryRepository {
	return &gormCategoryRepository{db: db}
}

func (r *gormCategoryRepository) ListPaged(limit, offset int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *gormCategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *gormCategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *gormCategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *gormCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}
