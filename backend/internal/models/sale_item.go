package models

import "github.com/shopspring/decimal"

type SaleItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	SaleID    uint    `json:"sale_id" gorm:"column:sale_id;not null"`
	ProductID uint    `json:"product_id" gorm:"column:product_id;not null"`
	Quantity  int             `json:"quantity" gorm:"column:quantity;not null;check:quantity > 0"`
	Price     decimal.Decimal `json:"price" gorm:"column:price;type:numeric(15,2);not null"`
	Subtotal  decimal.Decimal `json:"subtotal" gorm:"column:subtotal;type:numeric(15,2);not null"`
}
