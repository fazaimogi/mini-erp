package services

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/fazasuny/erp-system/internal/repositories"
)

type ReportsService interface {
	Dashboard() (DashboardMetrics, error)
	SalesReport(period string, startDate, endDate *time.Time) (SalesReport, error)
	ProductsReport() (ProductsReport, error)
}

type reportsService struct {
	repo repositories.ReportsRepository
}

func NewReportsService(repo repositories.ReportsRepository) ReportsService {
	return &reportsService{repo: repo}
}

// DashboardMetrics aggregates key metrics (PRD section 6.1).
type DashboardMetrics struct {
	TotalProducts   int64           `json:"total_products"`
	TotalCustomers  int64           `json:"total_customers"`
	TotalSales      int64           `json:"total_sales"`
	TotalRevenue    decimal.Decimal `json:"total_revenue"`
	LowStockCount   int64           `json:"low_stock_count"`
}

// SalesReport shows sales by period (PRD section 6.2, 16.1).
type SalesReport struct {
	Period       string          `json:"period"`
	TotalOrders  int64           `json:"total_orders"`
	TotalRevenue decimal.Decimal `json:"total_revenue"`
	TotalItems   int64           `json:"total_items"`
}

// ProductsReport shows top selling and low stock (PRD section 16.2).
type ProductsReport struct {
	TopSelling []ProductSaleInfo  `json:"top_selling"`
	LowStock   []ProductLowStockInfo `json:"low_stock"`
}

type ProductSaleInfo struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	SKU      string          `json:"sku"`
	Quantity int64           `json:"quantity"`
	Revenue  decimal.Decimal `json:"revenue"`
}

type ProductLowStockInfo struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SKU          string `json:"sku"`
	Stock        int    `json:"stock"`
	MinimumStock int    `json:"minimum_stock"`
}

func (s *reportsService) Dashboard() (DashboardMetrics, error) {
	totalProducts, err := s.repo.CountProducts()
	if err != nil {
		return DashboardMetrics{}, err
	}

	totalCustomers, err := s.repo.CountCustomers()
	if err != nil {
		return DashboardMetrics{}, err
	}

	totalSales, err := s.repo.CountSales()
	if err != nil {
		return DashboardMetrics{}, err
	}

	totalRevenue, err := s.repo.TotalRevenue()
	if err != nil {
		return DashboardMetrics{}, err
	}

	lowStockCount, err := s.repo.CountLowStock()
	if err != nil {
		return DashboardMetrics{}, err
	}

	return DashboardMetrics{
		TotalProducts:   totalProducts,
		TotalCustomers:  totalCustomers,
		TotalSales:      totalSales,
		TotalRevenue:    totalRevenue,
		LowStockCount:   lowStockCount,
	}, nil
}

func (s *reportsService) SalesReport(period string, startDate, endDate *time.Time) (SalesReport, error) {
	var orders int64
	var revenue decimal.Decimal
	var items int64
	var err error
	var periodLabel string

	switch period {
	case "today":
		orders, revenue, items, err = s.repo.SalesToday()
		periodLabel = "today"
	case "week":
		orders, revenue, items, err = s.repo.SalesThisWeek()
		periodLabel = "this_week"
	case "month":
		orders, revenue, items, err = s.repo.SalesThisMonth()
		periodLabel = "this_month"
	case "custom":
		if startDate == nil || endDate == nil {
			return SalesReport{}, ErrInvalidDateRange
		}
		orders, revenue, items, err = s.repo.SalesByPeriod(*startDate, *endDate)
		periodLabel = "custom"
	default:
		return SalesReport{}, ErrInvalidPeriod
	}

	if err != nil {
		return SalesReport{}, err
	}

	return SalesReport{
		Period:       periodLabel,
		TotalOrders:  orders,
		TotalRevenue: revenue,
		TotalItems:   items,
	}, nil
}

func (s *reportsService) ProductsReport() (ProductsReport, error) {
	topSelling, err := s.repo.TopSellingProducts(10)
	if err != nil {
		return ProductsReport{}, err
	}

	lowStock, err := s.repo.LowStockProducts(10)
	if err != nil {
		return ProductsReport{}, err
	}

	topSalesInfo := make([]ProductSaleInfo, len(topSelling))
	for i, p := range topSelling {
		topSalesInfo[i] = ProductSaleInfo{
			ID:       p.ID,
			Name:     p.Name,
			SKU:      p.SKU,
			Quantity: p.Quantity,
			Revenue:  p.Revenue,
		}
	}

	lowStockInfo := make([]ProductLowStockInfo, len(lowStock))
	for i, p := range lowStock {
		lowStockInfo[i] = ProductLowStockInfo{
			ID:           p.ID,
			Name:         p.Name,
			SKU:          p.SKU,
			Stock:        p.Stock,
			MinimumStock: p.MinimumStock,
		}
	}

	return ProductsReport{
		TopSelling: topSalesInfo,
		LowStock:   lowStockInfo,
	}, nil
}
