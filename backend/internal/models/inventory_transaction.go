package models

import "time"

type InventoryTransaction struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ProductID     uint      `json:"product_id" gorm:"column:product_id;not null"`
	Type          string    `json:"type" gorm:"column:type;not null"`
	Quantity      int       `json:"quantity" gorm:"column:quantity;not null"`
	PreviousStock int       `json:"previous_stock" gorm:"column:previous_stock;not null"`
	CurrentStock  int       `json:"current_stock" gorm:"column:current_stock;not null"`
	Reference     string    `json:"reference" gorm:"column:reference"`
	CreatedBy     *uint     `json:"created_by" gorm:"column:created_by"`
	CreatedAt     time.Time `json:"created_at" gorm:"column:created_at"`
}
