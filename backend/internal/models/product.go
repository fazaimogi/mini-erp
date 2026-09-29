package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Product represents the product entity (PRD section 7.1).
type Product struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	Name         string          `json:"name" gorm:"column:name;not null"`
	SKU          string          `json:"sku" gorm:"column:sku;unique;not null"`
	CategoryID   uint            `json:"category_id" gorm:"column:category_id;not null"`
	Description  string          `json:"description" gorm:"column:description"`
	Price        decimal.Decimal `json:"price" gorm:"column:price;type:numeric(15,2);not null"`
	Stock        int             `json:"stock" gorm:"column:stock;not null"`
	MinimumStock int             `json:"minimum_stock" gorm:"column:minimum_stock;not null"`
	CreatedAt    time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"column:updated_at"`
}
