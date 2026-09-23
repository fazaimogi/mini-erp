package services

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/fazasuny/erp-system/internal/repositories"
)

func setupService() ProductService {
	return NewProductService(repositories.NewInMemoryProductRepository())
}

func validInput() *ProductInput {
	return &ProductInput{
		Name:         "Mechanical Keyboard",
		SKU:          "KEY-001",
		CategoryID:   1,
		Description:  "Blue switch keyboard",
		Price:        decimal.NewFromInt(250000),
		Stock:        20,
		MinimumStock: 5,
	}
}

func TestCreateProduct(t *testing.T) {
	service := setupService()

	product, err := service.Create(validInput())
	if err != nil {
		t.Fatalf("expected create success, got error: %v", err)
	}
	if product.ID == 0 {
		t.Error("expected assigned id")
	}
	if product.SKU != "KEY-001" {
		t.Errorf("unexpected sku: %s", product.SKU)
	}
	if product.CreatedAt.IsZero() || product.UpdatedAt.IsZero() {
		t.Error("expected timestamps to be set")
	}
}

func TestCreateProduct_DuplicateSKU(t *testing.T) {
	service := setupService()

	if _, err := service.Create(validInput()); err != nil {
		t.Fatalf("seed product failed: %v", err)
	}

	duplicate := validInput()
	duplicate.Name = "Another Keyboard"
	_, err := service.Create(duplicate)

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got: %v", err)
	}
	if appErr != ErrDuplicateSKU {
		t.Errorf("expected ErrDuplicateSKU, got %v", err)
	}
}

func TestCreateProduct_InvalidInput(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*ProductInput)
	}{
		{"empty name", func(i *ProductInput) { i.Name = "  " }},
		{"empty sku", func(i *ProductInput) { i.SKU = "" }},
		{"missing category", func(i *ProductInput) { i.CategoryID = 0 }},
		{"negative price", func(i *ProductInput) { i.Price = decimal.NewFromInt(-1) }},
		{"negative stock", func(i *ProductInput) { i.Stock = -1 }},
		{"negative minimum stock", func(i *ProductInput) { i.MinimumStock = -1 }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := validInput()
			tc.mutate(input)

			if _, err := setupService().Create(input); err != ErrInvalidInput {
				t.Errorf("expected ErrInvalidInput, got %v", err)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	service := setupService()

	product, err := service.Create(validInput())
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	found, err := service.GetByID(product.ID)
	if err != nil {
		t.Fatalf("expected found product, got: %v", err)
	}
	if found.SKU != product.SKU {
		t.Errorf("unexpected product: %+v", found)
	}

	if _, err := service.GetByID(9999); err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestList_Pagination(t *testing.T) {
	service := setupService()

	for i := 0; i < 5; i++ {
		input := validInput()
		input.SKU = "KEY-00" + string(rune('2'+i))
		if _, err := service.Create(input); err != nil {
			t.Fatalf("create failed: %v", err)
		}
	}

	products, total, err := service.List(1, 2)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	// 3 seeded products + 5 created here.
	if total != 8 {
		t.Errorf("expected total 8, got %d", total)
	}
	if len(products) != 2 {
		t.Errorf("expected page size 2, got %d", len(products))
	}

	products, _, err = service.List(4, 2)
	if err != nil {
		t.Fatalf("list page 4 failed: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected last page size 2, got %d", len(products))
	}

	products, _, err = service.List(5, 2)
	if err != nil {
		t.Fatalf("out of range page failed: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("expected empty result, got %d", len(products))
	}
}

func TestUpdate(t *testing.T) {
	service := setupService()

	product, err := service.Create(validInput())
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	input := validInput()
	input.Name = "Updated Keyboard"
	input.Price = decimal.NewFromInt(300000)
	updated, err := service.Update(product.ID, input)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Updated Keyboard" || !updated.Price.Equal(decimal.NewFromInt(300000)) {
		t.Errorf("unexpected product: %+v", updated)
	}
	if !updated.CreatedAt.Equal(product.CreatedAt) {
		t.Error("expected created_at to be preserved")
	}

	if _, err := service.Update(9999, validInput()); err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
}

func TestUpdate_DuplicateSKUFromOtherProduct(t *testing.T) {
	service := setupService()

	first := validInput()
	first.SKU = "KEY-A"
	if _, err := service.Create(first); err != nil {
		t.Fatalf("create first failed: %v", err)
	}

	second := validInput()
	second.SKU = "KEY-B"
	created, err := service.Create(second)
	if err != nil {
		t.Fatalf("create second failed: %v", err)
	}

	update := validInput()
	update.SKU = "KEY-A"
	if _, err := service.Update(created.ID, update); err != ErrDuplicateSKU {
		t.Errorf("expected ErrDuplicateSKU, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	service := setupService()

	product, err := service.Create(validInput())
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := service.Delete(product.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := service.GetByID(product.ID); err != ErrProductNotFound {
		t.Errorf("expected product to be deleted, got %v", err)
	}

	if err := service.Delete(product.ID); err != ErrProductNotFound {
		t.Errorf("expected ErrProductNotFound on repeat delete, got %v", err)
	}
}
