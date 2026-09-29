package repositories

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/fazasuny/erp-system/internal/models"
)

type GormReportsRepository struct {
	db *gorm.DB
}

func NewGormReportsRepository(db *gorm.DB) ReportsRepository {
	return &GormReportsRepository{db: db}
}

// Dashboard metrics

func (r *GormReportsRepository) CountProducts() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormReportsRepository) CountCustomers() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Customer{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormReportsRepository) CountSales() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Sale{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *GormReportsRepository) TotalRevenue() (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.Model(&models.Sale{}).Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&total); err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (r *GormReportsRepository) CountLowStock() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Product{}).
		Where("stock <= minimum_stock").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Sales by period

func (r *GormReportsRepository) SalesByPeriod(startDate, endDate time.Time) (int64, decimal.Decimal, int64, error) {
	var count int64
	var revenue decimal.Decimal
	var items int64

	if err := r.db.Model(&models.Sale{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Count(&count).Error; err != nil {
		return 0, decimal.Zero, 0, err
	}

	if err := r.db.Model(&models.Sale{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Select("COALESCE(SUM(total_amount), 0)").
		Row().Scan(&revenue); err != nil {
		return 0, decimal.Zero, 0, err
	}

	if err := r.db.Model(&models.SaleItem{}).
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Where("sales.created_at >= ? AND sales.created_at < ?", startDate, endDate).
		Select("COALESCE(SUM(quantity), 0)").
		Row().Scan(&items); err != nil {
		return 0, decimal.Zero, 0, err
	}

	return count, revenue, items, nil
}

func (r *GormReportsRepository) SalesToday() (int64, decimal.Decimal, int64, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.Add(24 * time.Hour)
	return r.SalesByPeriod(startDate, endDate)
}

func (r *GormReportsRepository) SalesThisWeek() (int64, decimal.Decimal, int64, error) {
	now := time.Now()
	weekday := now.Weekday()
	startDate := now.AddDate(0, 0, int(time.Monday-weekday))
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate := startDate.AddDate(0, 0, 7)
	return r.SalesByPeriod(startDate, endDate)
}

func (r *GormReportsRepository) SalesThisMonth() (int64, decimal.Decimal, int64, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0)
	return r.SalesByPeriod(startDate, endDate)
}

// Product metrics

func (r *GormReportsRepository) TopSellingProducts(limit int) ([]ProductSalesInfo, error) {
	var results []ProductSalesInfo
	err := r.db.Table("sale_items").
		Select("products.id, products.name, products.sku, SUM(sale_items.quantity) as quantity, SUM(sale_items.subtotal) as revenue").
		Joins("JOIN products ON products.id = sale_items.product_id").
		Group("products.id, products.name, products.sku").
		Order("quantity DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}

func (r *GormReportsRepository) LowStockProducts(limit int) ([]LowStockInfo, error) {
	var results []LowStockInfo
	err := r.db.Model(&models.Product{}).
		Where("stock <= minimum_stock").
		Select("id, name, sku, stock, minimum_stock").
		Order("stock ASC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}
