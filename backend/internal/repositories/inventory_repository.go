package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
)

// ErrInsufficientStock is returned when a stock change would push the product
// below zero (PRD section 9.3 and 14.1). The PRD pins the error code string to
// INSUFFICIENT_STOCK.
var ErrInsufficientStock = errors.New("insufficient stock")

// StockChange is one atomic stock movement. Delta is signed: positive for stock
// in, negative for stock out.
type StockChange struct {
	ProductID uint
	Delta     int
	Type      string
	Reference string
	CreatedBy *uint
}

// InventoryRepository owns stock movements and their history.
//
// Stock itself lives on products.stock (PRD section 7.1), so browsing inventory
// is served by ProductRepository; this layer only records changes and keeps the
// history required by PRD section 10.
type InventoryRepository interface {
	ApplyStockChange(change StockChange) (*models.InventoryTransaction, error)
	ListTransactions(productID uint, limit, offset int) ([]models.InventoryTransaction, int64, error)
}
