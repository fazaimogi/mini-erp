package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	DiscountTypePercentage = "percentage"
	DiscountTypeFixed      = "fixed"
)

// Voucher is a promo code redeemable at checkout.
type Voucher struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	Code          string          `json:"code" gorm:"column:code;unique;not null"`
	DiscountType  string          `json:"discount_type" gorm:"column:discount_type;not null"`
	DiscountValue decimal.Decimal `json:"discount_value" gorm:"column:discount_value;type:numeric(15,2);not null"`
	MinPurchase   decimal.Decimal `json:"min_purchase" gorm:"column:min_purchase;type:numeric(15,2);not null"`
	UsageLimit    int             `json:"usage_limit" gorm:"column:usage_limit;not null"`
	UsedCount     int             `json:"used_count" gorm:"column:used_count;not null"`
	ExpiresAt     *time.Time      `json:"expires_at" gorm:"column:expires_at"`
	IsActive      bool            `json:"is_active" gorm:"column:is_active;not null"`
	CreatedAt     time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time       `json:"updated_at" gorm:"column:updated_at"`
}