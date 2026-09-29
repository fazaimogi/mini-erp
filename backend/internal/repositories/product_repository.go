package repositories

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/fazasuny/erp-system/internal/models"
)

// ErrProductNotFound is returned when no product matches the given id.
var ErrProductNotFound = errors.New("product not found")

// ProductRepository abstracts product persistence.
// The in-memory implementation is temporary; it will be replaced by a
// GORM/PostgreSQL implementation on Day 2 without touching the service layer.
type ProductRepository interface {
	// ListPaged returns one page of products plus the total row count.
	ListPaged(limit, offset int) ([]models.Product, int64, error)
	// ListLowStockPaged returns products where stock <= minimum_stock.
	ListLowStockPaged(limit, offset int) ([]models.Product, int64, error)
	GetByID(id uint) (*models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
}

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[uint]*models.Product
	nextID   uint
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	repo := &InMemoryProductRepository{
		products: make(map[uint]*models.Product),
		nextID:   1,
	}
	repo.seed()
	return repo
}

func (r *InMemoryProductRepository) seed() {
	now := time.Now()
	items := []models.Product{
		{Name: "Laptop Asus XYZ", SKU: "LAPTOP-XYZ", CategoryID: 1, Description: "Laptop untuk kantor", Price: decimal.NewFromInt(7500000), Stock: 10, MinimumStock: 3},
		{Name: "Mouse Logitech M185", SKU: "MOUSE-M185", CategoryID: 2, Description: "Mouse wireless", Price: decimal.NewFromInt(150000), Stock: 50, MinimumStock: 10},
		{Name: "Kursi Kantor Ergonomis", SKU: "KURSI-001", CategoryID: 3, Description: "Kursi kerja nyaman", Price: decimal.NewFromInt(1200000), Stock: 2, MinimumStock: 5},
	}
	for _, item := range items {
		product := item
		product.ID = r.nextID
		product.CreatedAt = now
		product.UpdatedAt = now
		r.products[r.nextID] = &product
		r.nextID++
	}
}

// ListPaged returns a page of products plus the total count. The SQL
// implementation pushes LIMIT/OFFSET to the database; this one slices in Go so
// both satisfy the same contract.
func (r *InMemoryProductRepository) ListPaged(limit, offset int) ([]models.Product, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]uint, 0, len(r.products))
	for id := range r.products {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	total := int64(len(ids))
	if offset >= len(ids) {
		return []models.Product{}, total, nil
	}

	end := offset + limit
	if end > len(ids) {
		end = len(ids)
	}

	products := make([]models.Product, 0, end-offset)
	for _, id := range ids[offset:end] {
		products = append(products, *r.products[id])
	}

	return products, total, nil
}

// ListLowStockPaged returns products at or below their minimum stock
// (PRD section 6.3).
func (r *InMemoryProductRepository) ListLowStockPaged(limit, offset int) ([]models.Product, int64, error) {
	all, _, err := r.ListPaged(len(r.products)+offset+1, 0)
	if err != nil {
		return nil, 0, err
	}

	low := make([]models.Product, 0, len(all))
	for _, product := range all {
		if product.Stock <= product.MinimumStock {
			low = append(low, product)
		}
	}

	total := int64(len(low))
	if offset >= len(low) {
		return []models.Product{}, total, nil
	}

	end := offset + limit
	if end > len(low) {
		end = len(low)
	}

	return low[offset:end], total, nil
}

// adjustStock applies a signed delta and reports the stock before and after.
// It takes the write lock so concurrent movements cannot interleave and both
// read the same previous stock.
func (r *InMemoryProductRepository) adjustStock(id uint, delta int) (int, int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, ok := r.products[id]
	if !ok {
		return 0, 0, ErrProductNotFound
	}

	previous := product.Stock
	current := previous + delta
	if current < 0 {
		return 0, 0, ErrInsufficientStock
	}

	product.Stock = current
	product.UpdatedAt = time.Now()

	return previous, current, nil
}

func (r *InMemoryProductRepository) GetByID(id uint) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}
	copy := *product
	return &copy, nil
}

func (r *InMemoryProductRepository) GetBySKU(sku string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, product := range r.products {
		if product.SKU == sku {
			copy := *product
			return &copy, nil
		}
	}
	return nil, ErrProductNotFound
}

func (r *InMemoryProductRepository) Create(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	product.ID = r.nextID
	product.CreatedAt = now
	product.UpdatedAt = now
	r.products[product.ID] = product
	r.nextID++

	return nil
}

func (r *InMemoryProductRepository) Update(product *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.products[product.ID]
	if !ok {
		return ErrProductNotFound
	}

	product.CreatedAt = existing.CreatedAt
	product.UpdatedAt = time.Now()
	r.products[product.ID] = product

	return nil
}

func (r *InMemoryProductRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return ErrProductNotFound
	}
	delete(r.products, id)

	return nil
}
