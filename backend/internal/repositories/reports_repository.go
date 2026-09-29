package repositories

import (
	"time"

	"github.com/shopspring/decimal"
)

// ReportsRepository defines queries for dashboard and reports (PRD section 6, 16).
type ReportsRepository interface {
	// Dashboard metrics
	CountProducts() (int64, error)
	CountCustomers() (int64, error)
	CountSales() (int64, error)
	TotalRevenue() (decimal.Decimal, error)
	CountLowStock() (int64, error) // stock <= minimum_stock

	// Sales by period
	SalesByPeriod(startDate, endDate time.Time) (int64, decimal.Decimal, int64, error) // count, revenue, items
	SalesToday() (int64, decimal.Decimal, int64, error)
	SalesThisWeek() (int64, decimal.Decimal, int64, error)
	SalesThisMonth() (int64, decimal.Decimal, int64, error)

	// Product metrics
	TopSellingProducts(limit int) ([]ProductSalesInfo, error)
	LowStockProducts(limit int) ([]LowStockInfo, error)
}

type ProductSalesInfo struct {
	ID       uint
	Name     string
	SKU      string
	Quantity int64
	Revenue  decimal.Decimal
}

type LowStockInfo struct {
	ID           uint
	Name         string
	SKU          string
	Stock        int
	MinimumStock int
}
