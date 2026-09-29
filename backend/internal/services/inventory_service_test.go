package services

import (
	"errors"
	"testing"

	"github.com/fazasuny/erp-system/internal/repositories"
)

type inventoryFixture struct {
	service   InventoryService
	products  *repositories.InMemoryProductRepository
	inventory *repositories.InMemoryInventoryRepository
}

// Seat 1 is "Laptop Asus XYZ" with stock 10 and minimum stock 3; seat 3 is
// "Kursi Kantor Ergonomis" with stock 2 (already low stock).
func setupInventory() inventoryFixture {
	products := repositories.NewInMemoryProductRepository()
	inventory := repositories.NewInMemoryInventoryRepository(products)

	return inventoryFixture{
		service:   NewInventoryService(NewProductService(products), inventory),
		products:  products,
		inventory: inventory,
	}
}

func stockOf(t *testing.T, products *repositories.InMemoryProductRepository, id uint) int {
	t.Helper()

	product, err := products.GetByID(id)
	if err != nil {
		t.Fatalf("lookup product %d failed: %v", id, err)
	}
	return product.Stock
}

func TestStockInAddsStockAndRecordsHistory(t *testing.T) {
	f := setupInventory()

	record, err := f.service.StockIn(&StockMovementInput{ProductID: 1, Quantity: 5, Reference: "PO-001"}, 7)
	if err != nil {
		t.Fatalf("expected stock in success, got: %v", err)
	}

	// PRD section 10.3: 10 -> 15, type IN, quantity 5.
	if record.Type != "IN" || record.Quantity != 5 {
		t.Errorf("unexpected transaction: %+v", record)
	}
	if record.PreviousStock != 10 || record.CurrentStock != 15 {
		t.Errorf("expected 10 -> 15, got %d -> %d", record.PreviousStock, record.CurrentStock)
	}
	if record.CreatedBy == nil || *record.CreatedBy != 7 {
		t.Errorf("expected created_by 7, got %v", record.CreatedBy)
	}
	if got := stockOf(t, f.products, 1); got != 15 {
		t.Errorf("expected product stock 15, got %d", got)
	}

	history, total, err := f.service.History(1, 1, 20)
	if err != nil {
		t.Fatalf("history failed: %v", err)
	}
	if total != 1 || len(history) != 1 || history[0].ID != record.ID {
		t.Errorf("expected history to contain the new record, got total=%d len=%d", total, len(history))
	}
}

func TestStockOutDeductsStockAndRecordsHistory(t *testing.T) {
	f := setupInventory()

	if _, err := f.service.StockIn(&StockMovementInput{ProductID: 1, Quantity: 5, Reference: "PO-001"}, 1); err != nil {
		t.Fatalf("setup stock in failed: %v", err)
	}

	record, err := f.service.StockOut(&StockMovementInput{ProductID: 1, Quantity: 3, Reference: "INV-100"}, 7)
	if err != nil {
		t.Fatalf("expected stock out success, got: %v", err)
	}

	// PRD section 10.3: 15 -> 12, type OUT, quantity 3.
	if record.Type != "OUT" || record.Quantity != 3 {
		t.Errorf("unexpected transaction: %+v", record)
	}
	if record.PreviousStock != 15 || record.CurrentStock != 12 {
		t.Errorf("expected 15 -> 12, got %d -> %d", record.PreviousStock, record.CurrentStock)
	}
	if got := stockOf(t, f.products, 1); got != 12 {
		t.Errorf("expected product stock 12, got %d", got)
	}
}

func TestStockOutBeyondAvailableIsRejected(t *testing.T) {
	f := setupInventory()

	// PRD section 9.3: current stock 10, requested 11 must be refused.
	_, err := f.service.StockOut(&StockMovementInput{ProductID: 1, Quantity: 11}, 7)

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got: %v", err)
	}
	if appErr.Code != "INSUFFICIENT_STOCK" || appErr.HTTPStatus != 409 {
		t.Errorf("unexpected error: %+v", appErr)
	}

	if got := stockOf(t, f.products, 1); got != 10 {
		t.Errorf("stock must be unchanged after a rejected movement, got %d", got)
	}

	_, total, err := f.service.History(1, 1, 20)
	if err != nil {
		t.Fatalf("history failed: %v", err)
	}
	if total != 0 {
		t.Errorf("rejected movement must not be recorded, history total=%d", total)
	}
}

func TestStockOutExactlyAvailableIsAllowed(t *testing.T) {
	f := setupInventory()

	record, err := f.service.StockOut(&StockMovementInput{ProductID: 1, Quantity: 10}, 7)
	if err != nil {
		t.Fatalf("expected exact stock out to succeed, got: %v", err)
	}
	if record.CurrentStock != 0 {
		t.Errorf("expected stock 0, got %d", record.CurrentStock)
	}
	if got := stockOf(t, f.products, 1); got != 0 {
		t.Errorf("expected product stock 0, got %d", got)
	}
}

func TestMovementsRejectNonPositiveQuantity(t *testing.T) {
	f := setupInventory()

	for _, quantity := range []int{0, -5} {
		t.Run("quantity", func(t *testing.T) {
			for name, run := range map[string]func() error{
				"stock in": func() error {
					_, err := f.service.StockIn(&StockMovementInput{ProductID: 1, Quantity: quantity}, 7)
					return err
				},
				"stock out": func() error {
					_, err := f.service.StockOut(&StockMovementInput{ProductID: 1, Quantity: quantity}, 7)
					return err
				},
			} {
				err := run()

				var appErr *AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("%s: expected AppError, got: %v", name, err)
				}
				if appErr.Code != "INVALID_QUANTITY" || appErr.HTTPStatus != 422 {
					t.Errorf("%s: unexpected error: %+v", name, appErr)
				}
			}
		})
	}

	if got := stockOf(t, f.products, 1); got != 10 {
		t.Errorf("stock must be untouched, got %d", got)
	}
}

func TestMovementOnUnknownProductReturnsNotFound(t *testing.T) {
	f := setupInventory()

	_, err := f.service.StockIn(&StockMovementInput{ProductID: 9999, Quantity: 1}, 7)

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got: %v", err)
	}
	if appErr.Code != "PRODUCT_NOT_FOUND" || appErr.HTTPStatus != 404 {
		t.Errorf("unexpected error: %+v", appErr)
	}
}

func TestHistoryFiltersByProductAndPaginatesNewestFirst(t *testing.T) {
	f := setupInventory()

	for i := 1; i <= 3; i++ {
		if _, err := f.service.StockIn(&StockMovementInput{ProductID: 1, Quantity: i, Reference: "PO"}, 1); err != nil {
			t.Fatalf("stock in %d failed: %v", i, err)
		}
	}
	if _, err := f.service.StockIn(&StockMovementInput{ProductID: 2, Quantity: 1}, 1); err != nil {
		t.Fatalf("stock in on product 2 failed: %v", err)
	}

	page, total, err := f.service.History(1, 1, 2)
	if err != nil {
		t.Fatalf("history failed: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3 transactions for product 1, got %d", total)
	}
	if len(page) != 2 {
		t.Fatalf("expected 2 rows on page 1, got %d", len(page))
	}
	if page[0].Quantity != 3 || page[1].Quantity != 2 {
		t.Errorf("expected newest first (3, 2), got (%d, %d)", page[0].Quantity, page[1].Quantity)
	}

	second, _, err := f.service.History(1, 2, 2)
	if err != nil {
		t.Fatalf("history page 2 failed: %v", err)
	}
	if len(second) != 1 || second[0].Quantity != 1 {
		t.Errorf("expected the oldest row on page 2, got %+v", second)
	}

	all, total, err := f.service.History(0, 1, 20)
	if err != nil {
		t.Fatalf("history without filter failed: %v", err)
	}
	if total != 4 || len(all) != 4 {
		t.Errorf("expected 4 transactions overall, got total=%d len=%d", total, len(all))
	}
}

func TestInventoryListAndDetailExposeProductStock(t *testing.T) {
	f := setupInventory()

	products, total, err := f.service.List(1, 20)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 3 || len(products) != 3 {
		t.Errorf("expected 3 seeded products, got total=%d len=%d", total, len(products))
	}

	product, err := f.service.GetByID(1)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if product.Stock != 10 {
		t.Errorf("expected stock 10, got %d", product.Stock)
	}

	if _, err := f.service.GetByID(9999); !errors.Is(err, ErrProductNotFound) {
		t.Errorf("expected ErrProductNotFound, got: %v", err)
	}
}
