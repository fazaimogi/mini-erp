package repositories

import (
	"sort"
	"strings"
	"sync"

	"github.com/fazasuny/erp-system/internal/models"
)

// InMemoryVoucherRepository backs the voucher unit tests.
type InMemoryVoucherRepository struct {
	mu       sync.Mutex
	vouchers []models.Voucher
	nextID   uint
}

func NewInMemoryVoucherRepository() *InMemoryVoucherRepository {
	return &InMemoryVoucherRepository{nextID: 1}
}

func (r *InMemoryVoucherRepository) ListPaged(limit, offset int) ([]models.Voucher, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	sorted := make([]models.Voucher, len(r.vouchers))
	copy(sorted, r.vouchers)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID > sorted[j].ID })

	total := int64(len(sorted))
	if offset >= len(sorted) {
		return []models.Voucher{}, total, nil
	}
	end := offset + limit
	if end > len(sorted) {
		end = len(sorted)
	}
	return sorted[offset:end], total, nil
}

func (r *InMemoryVoucherRepository) GetByID(id uint) (*models.Voucher, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.vouchers {
		if r.vouchers[i].ID == id {
			return &r.vouchers[i], nil
		}
	}
	return nil, ErrVoucherNotFound
}

func (r *InMemoryVoucherRepository) GetByCode(code string) (*models.Voucher, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	wanted := strings.ToUpper(strings.TrimSpace(code))
	for i := range r.vouchers {
		if r.vouchers[i].Code == wanted {
			return &r.vouchers[i], nil
		}
	}
	return nil, ErrVoucherNotFound
}

func (r *InMemoryVoucherRepository) Create(voucher *models.Voucher) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	voucher.ID = r.nextID
	r.nextID++
	r.vouchers = append(r.vouchers, *voucher)
	return nil
}

func (r *InMemoryVoucherRepository) Update(voucher *models.Voucher) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.vouchers {
		if r.vouchers[i].ID == voucher.ID {
			r.vouchers[i] = *voucher
			return nil
		}
	}
	return ErrVoucherNotFound
}

func (r *InMemoryVoucherRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.vouchers {
		if r.vouchers[i].ID == id {
			r.vouchers = append(r.vouchers[:i], r.vouchers[i+1:]...)
			return nil
		}
	}
	return ErrVoucherNotFound
}