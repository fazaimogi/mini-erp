package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Sale struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	InvoiceNumber string          `json:"invoice_number" gorm:"column:invoice_number;unique;not null"`
	CustomerID    uint            `json:"customer_id" gorm:"column:customer_id;not null"`
	UserID        uint            `json:"user_id" gorm:"column:user_id;not null"`
	TotalAmount   decimal.Decimal `json:"total_amount" gorm:"column:total_amount;type:numeric(15,2);not null"`
	DiscountAmount decimal.Decimal `json:"discount_amount" gorm:"column:discount_amount;type:numeric(15,2);not null"`
	VoucherID     *uint           `json:"voucher_id" gorm:"column:voucher_id"`
	Status        string          `json:"status" gorm:"column:status;not null"`
	Items         []SaleItem      `json:"items" gorm:"foreignKey:SaleID"`
	CreatedAt     time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     time.Time       `json:"updated_at" gorm:"column:updated_at"`
}
