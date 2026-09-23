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
	List() ([]models.Product, error)
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

func (r *InMemoryProductRepository) List() ([]models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]uint, 0, len(r.products))
	for id := range r.products {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	products := make([]models.Product, 0, len(r.products))
	for _, id := range ids {
		products = append(products, *r.products[id])
	}
	return products, nil
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
