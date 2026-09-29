package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

const (
	TransactionTypeSale = "OUT"
)

var (
	ErrSaleNotFound     = errors.New("sale not found")
	ErrInvalidSaleInput = errors.New("invalid sale input")
)

type CreateSaleInput struct {
	CustomerID uint                `json:"customer_id" binding:"required"`
	Items      []CreateSaleItemInput `json:"items" binding:"required,min=1"`
}

type CreateSaleItemInput struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

type SaleService interface {
	GetByID(id uint) (*models.Sale, error)
	List(page, limit int) ([]models.Sale, int64, error)
	Create(input *CreateSaleInput, userID uint) (*models.Sale, error)
}

type saleService struct {
	repo      repositories.SaleRepository
	products  repositories.ProductRepository
	inventory repositories.InventoryRepository
	db        *gorm.DB
}

func NewSaleService(
	repo repositories.SaleRepository,
	products repositories.ProductRepository,
	inventory repositories.InventoryRepository,
	db *gorm.DB,
) SaleService {
	return &saleService{
		repo:      repo,
		products:  products,
		inventory: inventory,
		db:        db,
	}
}

func (s *saleService) GetByID(id uint) (*models.Sale, error) {
	return s.repo.GetByID(id)
}

func (s *saleService) List(page, limit int) ([]models.Sale, int64, error) {
	return s.repo.ListPaged(limit, (page-1)*limit)
}

// Create orchestrates the sale creation with stock validation, deduction, and
// inventory transaction in a single DB transaction (PRD section 14.3).
func (s *saleService) Create(input *CreateSaleInput, userID uint) (*models.Sale, error) {
	if err := validateSaleInput(input); err != nil {
		return nil, err
	}

	// Pre-validate stock for all items (fail fast before transaction).
	stockMap := make(map[uint]int)
	itemMap := make(map[uint]*CreateSaleItemInput)
	var totalAmount decimal.Decimal

	for i := range input.Items {
		item := &input.Items[i]
		product, err := s.products.GetByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("%w for product %d", ErrInsufficientStock, item.ProductID)
		}

		stockMap[item.ProductID] = product.Stock
		itemMap[item.ProductID] = item
		subtotal := decimal.NewFromInt(int64(item.Quantity)).Mul(product.Price)
		totalAmount = totalAmount.Add(subtotal)
	}

	// Transaction: create sale, items, deduct stock, record inventory txns.
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	sale := &models.Sale{
		InvoiceNumber: generateInvoiceNumber(),
		CustomerID:    input.CustomerID,
		UserID:        userID,
		TotalAmount:   totalAmount,
		Status:        "completed",
	}

	if err := tx.Create(sale).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	saleItems := make([]models.SaleItem, len(input.Items))
	for i, item := range input.Items {
		product, _ := s.products.GetByID(item.ProductID)
		subtotal := decimal.NewFromInt(int64(item.Quantity)).Mul(product.Price)

		saleItems[i] = models.SaleItem{
			SaleID:    sale.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Subtotal:  subtotal,
		}

		if err := tx.Create(&saleItems[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Deduct stock: update product and record inventory transaction.
		newStock := product.Stock - item.Quantity
		if err := tx.Model(&models.Product{}).Where("id = ?", item.ProductID).Update("stock", newStock).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		invTxn := &models.InventoryTransaction{
			ProductID:     item.ProductID,
			Type:          TransactionTypeSale,
			Quantity:      item.Quantity,
			PreviousStock: product.Stock,
			CurrentStock:  newStock,
			Reference:     fmt.Sprintf("Sale #%s", sale.InvoiceNumber),
			CreatedBy:     &userID,
			CreatedAt:     time.Now(),
		}

		if err := tx.Create(invTxn).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with items and relations.
	sale.Items = saleItems
	return sale, nil
}

func generateInvoiceNumber() string {
	return fmt.Sprintf("INV-%d", time.Now().UnixNano())
}

func validateSaleInput(input *CreateSaleInput) error {
	if input.CustomerID == 0 {
		return fmt.Errorf("%w: customer_id required", ErrInvalidSaleInput)
	}
	if len(input.Items) == 0 {
		return fmt.Errorf("%w: at least one item required", ErrInvalidSaleInput)
	}
	for _, item := range input.Items {
		if item.ProductID == 0 {
			return fmt.Errorf("%w: product_id required", ErrInvalidSaleInput)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("%w: quantity must be > 0", ErrInvalidSaleInput)
		}
	}
	return nil
}
