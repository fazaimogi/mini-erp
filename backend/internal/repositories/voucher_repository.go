package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
)

var ErrVoucherNotFound = errors.New("voucher not found")

// VoucherRepository abstracts voucher persistence.
type VoucherRepository interface {
	ListPaged(limit, offset int) ([]models.Voucher, int64, error)
	GetByID(id uint) (*models.Voucher, error)
	GetByCode(code string) (*models.Voucher, error)
	Create(voucher *models.Voucher) error
	Update(voucher *models.Voucher) error
	Delete(id uint) error
}