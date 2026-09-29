package services

import (
	"errors"
	"strings"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

// Movement types recorded in inventory_transactions (PRD section 10.2).
const (
	TransactionTypeIn  = "IN"
	TransactionTypeOut = "OUT"
)

const maxReferenceLength = 255

// StockMovementInput is the payload for both POST /api/inventory/stock-in and
// POST /api/inventory/stock-out. PRD section 23 keeps the product id out of the
// path, and the sign of the movement comes from the endpoint, not the body.
type StockMovementInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	// Quantity deliberately has no binding tag: gin's "required" rejects 0
	// before the service can, which would report VALIDATION_ERROR instead of the
	// more specific INVALID_QUANTITY.
	Quantity  int    `json:"quantity"`
	Reference string `json:"reference"`
}

type InventoryService interface {
	// List and GetByID expose stock, which lives on products (PRD section 7.1).
	List(page, limit int) ([]models.Product, int64, error)
	GetByID(id uint) (*models.Product, error)
	StockIn(input *StockMovementInput, userID uint) (*models.InventoryTransaction, error)
	StockOut(input *StockMovementInput, userID uint) (*models.InventoryTransaction, error)
	History(productID uint, page, limit int) ([]models.InventoryTransaction, int64, error)
}

type inventoryService struct {
	products  ProductService
	inventory repositories.InventoryRepository
}

func NewInventoryService(products ProductService, inventory repositories.InventoryRepository) InventoryService {
	return &inventoryService{products: products, inventory: inventory}
}

func (s *inventoryService) List(page, limit int) ([]models.Product, int64, error) {
	return s.products.List(page, limit)
}

func (s *inventoryService) GetByID(id uint) (*models.Product, error) {
	return s.products.GetByID(id)
}

// StockIn adds stock and records an IN transaction (PRD section 9.1).
func (s *inventoryService) StockIn(input *StockMovementInput, userID uint) (*models.InventoryTransaction, error) {
	return s.apply(input, userID, TransactionTypeIn, 1)
}

// StockOut removes stock and records an OUT transaction. The repository refuses
// the change when it would push stock below zero (PRD section 9.3).
func (s *inventoryService) StockOut(input *StockMovementInput, userID uint) (*models.InventoryTransaction, error) {
	return s.apply(input, userID, TransactionTypeOut, -1)
}

func (s *inventoryService) History(productID uint, page, limit int) ([]models.InventoryTransaction, int64, error) {
	return s.inventory.ListTransactions(productID, limit, (page-1)*limit)
}

func (s *inventoryService) apply(input *StockMovementInput, userID uint, movementType string, sign int) (*models.InventoryTransaction, error) {
	if input.ProductID == 0 {
		return nil, ErrInvalidInput
	}

	reference := strings.TrimSpace(input.Reference)
	if input.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if len(reference) > maxReferenceLength {
		return nil, ErrInvalidInput
	}

	createdBy := userID
	record, err := s.inventory.ApplyStockChange(repositories.StockChange{
		ProductID: input.ProductID,
		Delta:     sign * input.Quantity,
		Type:      movementType,
		Reference: reference,
		CreatedBy: &createdBy,
	})
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrProductNotFound):
			return nil, ErrProductNotFound
		case errors.Is(err, repositories.ErrInsufficientStock):
			return nil, ErrInsufficientStock
		default:
			return nil, err
		}
	}

	return record, nil
}
