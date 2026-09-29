package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormInventoryRepository struct {
	db *gorm.DB
}

func NewGormInventoryRepository(db *gorm.DB) *GormInventoryRepository {
	return &GormInventoryRepository{db: db}
}

// ApplyStockChange moves stock and writes the history row inside one database
// transaction, so a failed insert leaves the product untouched (PRD section
// 14.3 and 26.2).
//
// The product row is read with SELECT ... FOR UPDATE: two concurrent movements
// on the same product then serialise instead of both reading the same
// previous_stock and one silently overwriting the other.
func (r *GormInventoryRepository) ApplyStockChange(change StockChange) (*models.InventoryTransaction, error) {
	var record models.InventoryTransaction

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var product models.Product

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&product, change.ProductID).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrProductNotFound
			}
			return err
		}

		previous := product.Stock
		current := previous + change.Delta
		if current < 0 {
			return ErrInsufficientStock
		}

		if err := tx.Model(&models.Product{}).
			Where("id = ?", change.ProductID).
			Update("stock", current).
			Error; err != nil {
			return err
		}

		quantity := change.Delta
		if quantity < 0 {
			quantity = -quantity
		}

		record = models.InventoryTransaction{
			ProductID:     change.ProductID,
			Type:          change.Type,
			Quantity:      quantity,
			PreviousStock: previous,
			CurrentStock:  current,
			Reference:     change.Reference,
			CreatedBy:     change.CreatedBy,
		}

		return tx.Create(&record).Error
	})

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *GormInventoryRepository) ListTransactions(productID uint, limit, offset int) ([]models.InventoryTransaction, int64, error) {
	// Rebuilt per call: reusing one *gorm.DB for Count and Find leaks the
	// count's conditions into the select.
	query := func() *gorm.DB {
		q := r.db.Model(&models.InventoryTransaction{})
		if productID != 0 {
			q = q.Where("product_id = ?", productID)
		}
		return q
	}

	var total int64
	if err := query().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []models.InventoryTransaction
	err := query().
		Order("id DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).
		Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
