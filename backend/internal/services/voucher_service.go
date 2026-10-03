package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

// VoucherInput is the accepted payload for create and update.
type VoucherInput struct {
	Code          string          `json:"code" binding:"required"`
	DiscountType  string          `json:"discount_type" binding:"required"`
	DiscountValue decimal.Decimal `json:"discount_value"`
	MinPurchase   decimal.Decimal `json:"min_purchase"`
	UsageLimit    int             `json:"usage_limit"`
	ExpiresAt     *time.Time      `json:"expires_at"`
	IsActive      *bool           `json:"is_active"`
}

// VoucherValidation is what /vouchers/validate returns for a valid code.
type VoucherValidation struct {
	Code           string          `json:"code"`
	DiscountType   string          `json:"discount_type"`
	DiscountValue  decimal.Decimal `json:"discount_value"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	FinalAmount    decimal.Decimal `json:"final_amount"`
}

type VoucherService interface {
	List(page, limit int) ([]models.Voucher, int64, error)
	GetByID(id uint) (*models.Voucher, error)
	Create(input *VoucherInput) (*models.Voucher, error)
	Update(id uint, input *VoucherInput) (*models.Voucher, error)
	Delete(id uint) error
	Validate(code string, totalAmount decimal.Decimal) (*VoucherValidation, error)
}

type voucherService struct {
	repo repositories.VoucherRepository
}

func NewVoucherService(repo repositories.VoucherRepository) VoucherService {
	return &voucherService{repo: repo}
}

func (s *voucherService) List(page, limit int) ([]models.Voucher, int64, error) {
	return s.repo.ListPaged(limit, (page-1)*limit)
}

func (s *voucherService) GetByID(id uint) (*models.Voucher, error) {
	voucher, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrVoucherNotFound {
			return nil, ErrVoucherNotFound
		}
		return nil, err
	}
	return voucher, nil
}

func (s *voucherService) Create(input *VoucherInput) (*models.Voucher, error) {
	if err := validateVoucherInput(input); err != nil {
		return nil, err
	}

	code := normalizeVoucherCode(input.Code)
	if existing, err := s.repo.GetByCode(code); err == nil && existing != nil {
		return nil, ErrDuplicateVoucherCode
	} else if err != nil && err != repositories.ErrVoucherNotFound {
		return nil, err
	}

	active := true
	if input.IsActive != nil {
		active = *input.IsActive
	}

	voucher := &models.Voucher{
		Code:          code,
		DiscountType:  input.DiscountType,
		DiscountValue: input.DiscountValue,
		MinPurchase:   input.MinPurchase,
		UsageLimit:    input.UsageLimit,
		ExpiresAt:     input.ExpiresAt,
		IsActive:      active,
	}

	if err := s.repo.Create(voucher); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *voucherService) Update(id uint, input *VoucherInput) (*models.Voucher, error) {
	if err := validateVoucherInput(input); err != nil {
		return nil, err
	}

	voucher, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	code := normalizeVoucherCode(input.Code)
	if existing, err := s.repo.GetByCode(code); err == nil && existing != nil && existing.ID != id {
		return nil, ErrDuplicateVoucherCode
	} else if err != nil && err != repositories.ErrVoucherNotFound {
		return nil, err
	}

	voucher.Code = code
	voucher.DiscountType = input.DiscountType
	voucher.DiscountValue = input.DiscountValue
	voucher.MinPurchase = input.MinPurchase
	voucher.UsageLimit = input.UsageLimit
	voucher.ExpiresAt = input.ExpiresAt
	if input.IsActive != nil {
		voucher.IsActive = *input.IsActive
	}

	if err := s.repo.Update(voucher); err != nil {
		return nil, err
	}

	return voucher, nil
}

func (s *voucherService) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// Validate resolves a code and returns the discount it yields for the given
// total. Invalid codes surface as AppError so the client gets a reason.
func (s *voucherService) Validate(code string, totalAmount decimal.Decimal) (*VoucherValidation, error) {
	voucher, err := s.repo.GetByCode(normalizeVoucherCode(code))
	if err != nil {
		if err == repositories.ErrVoucherNotFound {
			return nil, ErrVoucherNotFound
		}
		return nil, err
	}

	if err := checkVoucherUsable(voucher, totalAmount); err != nil {
		return nil, err
	}

	discount := ComputeVoucherDiscount(voucher, totalAmount)
	return &VoucherValidation{
		Code:           voucher.Code,
		DiscountType:   voucher.DiscountType,
		DiscountValue:  voucher.DiscountValue,
		DiscountAmount: discount,
		FinalAmount:    totalAmount.Sub(discount),
	}, nil
}

// ComputeVoucherDiscount applies a voucher to a total and clamps the result so
// it never exceeds the total (no negative payable amount).
func ComputeVoucherDiscount(voucher *models.Voucher, total decimal.Decimal) decimal.Decimal {
	var discount decimal.Decimal
	if voucher.DiscountType == models.DiscountTypePercentage {
		discount = total.Mul(voucher.DiscountValue).Div(decimal.NewFromInt(100))
	} else {
		discount = voucher.DiscountValue
	}

	discount = discount.Round(2)
	if discount.GreaterThan(total) {
		discount = total
	}
	if discount.IsNegative() {
		discount = decimal.Zero
	}
	return discount
}

// checkVoucherUsable enforces active flag, expiry, usage limit and minimum
// purchase. Exported logic lives here so sale creation can reuse it.
func checkVoucherUsable(voucher *models.Voucher, totalAmount decimal.Decimal) error {
	if !voucher.IsActive {
		return ErrVoucherInactive
	}
	if voucher.ExpiresAt != nil && voucher.ExpiresAt.Before(time.Now()) {
		return ErrVoucherExpired
	}
	if voucher.UsageLimit > 0 && voucher.UsedCount >= voucher.UsageLimit {
		return ErrVoucherUsageLimit
	}
	if totalAmount.LessThan(voucher.MinPurchase) {
		return ErrVoucherMinPurchase
	}
	return nil
}

func normalizeVoucherCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func validateVoucherInput(input *VoucherInput) error {
	if normalizeVoucherCode(input.Code) == "" {
		return NewAppError("INVALID_INPUT", "voucher code is required", http.StatusUnprocessableEntity)
	}
	if input.DiscountType != models.DiscountTypePercentage && input.DiscountType != models.DiscountTypeFixed {
		return NewAppError("INVALID_INPUT", "discount_type must be 'percentage' or 'fixed'", http.StatusUnprocessableEntity)
	}
	if input.DiscountValue.LessThanOrEqual(decimal.Zero) {
		return NewAppError("INVALID_INPUT", "discount_value must be greater than zero", http.StatusUnprocessableEntity)
	}
	if input.DiscountType == models.DiscountTypePercentage && input.DiscountValue.GreaterThan(decimal.NewFromInt(100)) {
		return NewAppError("INVALID_INPUT", "percentage discount cannot exceed 100", http.StatusUnprocessableEntity)
	}
	if input.MinPurchase.IsNegative() || input.UsageLimit < 0 {
		return NewAppError("INVALID_INPUT", "min_purchase and usage_limit cannot be negative", http.StatusUnprocessableEntity)
	}
	return nil
}