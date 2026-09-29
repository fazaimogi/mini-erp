package repositories

import (
	"sort"
	"sync"
	"time"

	"github.com/fazasuny/erp-system/internal/models"
)

// InMemoryInventoryRepository backs the inventory unit tests. It mutates the
// shared InMemoryProductRepository so stock and history stay consistent, which
// is what a real transaction provides.
type InMemoryInventoryRepository struct {
	mu           sync.Mutex
	products     *InMemoryProductRepository
	transactions []models.InventoryTransaction
	nextID       uint
}

func NewInMemoryInventoryRepository(products *InMemoryProductRepository) *InMemoryInventoryRepository {
	return &InMemoryInventoryRepository{
		products: products,
		nextID:   1,
	}
}

func (r *InMemoryInventoryRepository) ApplyStockChange(change StockChange) (*models.InventoryTransaction, error) {
	previous, current, err := r.products.adjustStock(change.ProductID, change.Delta)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	quantity := change.Delta
	if quantity < 0 {
		quantity = -quantity
	}

	record := models.InventoryTransaction{
		ID:            r.nextID,
		ProductID:     change.ProductID,
		Type:          change.Type,
		Quantity:      quantity,
		PreviousStock: previous,
		CurrentStock:  current,
		Reference:     change.Reference,
		CreatedBy:     change.CreatedBy,
		CreatedAt:     time.Now(),
	}

	r.transactions = append(r.transactions, record)
	r.nextID++

	return &record, nil
}

func (r *InMemoryInventoryRepository) ListTransactions(productID uint, limit, offset int) ([]models.InventoryTransaction, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	matching := make([]models.InventoryTransaction, 0, len(r.transactions))
	for _, transaction := range r.transactions {
		if productID != 0 && transaction.ProductID != productID {
			continue
		}
		matching = append(matching, transaction)
	}

	// Newest first, mirroring the SQL implementation's ORDER BY id DESC.
	sort.Slice(matching, func(i, j int) bool { return matching[i].ID > matching[j].ID })

	total := int64(len(matching))
	if offset >= len(matching) {
		return []models.InventoryTransaction{}, total, nil
	}

	end := offset + limit
	if end > len(matching) {
		end = len(matching)
	}

	return matching[offset:end], total, nil
}
