package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/fazasuny/erp-system/internal/models"
)

type gormVoucherRepository struct {
	db *gorm.DB
}

func NewGormVoucherRepository(db *gorm.DB) VoucherRepository {
	return &gormVoucherRepository{db: db}
}

func (r *gormVoucherRepository) ListPaged(limit, offset int) ([]models.Voucher, int64, error) {
	var vouchers []models.Voucher
	var total int64

	if err := r.db.Model(&models.Voucher{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&vouchers).Error; err != nil {
		return nil, 0, err
	}

	return vouchers, total, nil
}

func (r *gormVoucherRepository) GetByID(id uint) (*models.Voucher, error) {
	var voucher models.Voucher
	if err := r.db.First(&voucher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVoucherNotFound
		}
		return nil, err
	}
	return &voucher, nil
}

func (r *gormVoucherRepository) GetByCode(code string) (*models.Voucher, error) {
	var voucher models.Voucher
	if err := r.db.Where("code = ?", code).First(&voucher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVoucherNotFound
		}
		return nil, err
	}
	return &voucher, nil
}

func (r *gormVoucherRepository) Create(voucher *models.Voucher) error {
	return r.db.Create(voucher).Error
}

func (r *gormVoucherRepository) Update(voucher *models.Voucher) error {
	return r.db.Save(voucher).Error
}

func (r *gormVoucherRepository) Delete(id uint) error {
	return r.db.Delete(&models.Voucher{}, id).Error
}