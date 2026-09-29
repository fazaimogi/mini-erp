package repositories

import (
	"github.com/fazasuny/erp-system/internal/models"
	"gorm.io/gorm"
)

type GormCustomerRepository struct {
	db *gorm.DB
}

func NewGormCustomerRepository(db *gorm.DB) *GormCustomerRepository {
	return &GormCustomerRepository{
		db: db,
	}
}

func (r *GormCustomerRepository) ListPaged(limit, offset int) ([]models.Customer, int64, error) {
	query := func() *gorm.DB {
		return r.db.Model(&models.Customer{})
	}

	var total int64
	if err := query().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var customers []models.Customer
	err := query().
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&customers).
		Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *GormCustomerRepository) GetByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("id = ?", id).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCustomerNotFound
		}
		return nil, err
	}
	return &customer, nil
}

func (r *GormCustomerRepository) GetByEmail(email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("email = ?", email).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *GormCustomerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *GormCustomerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *GormCustomerRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Customer{}, id)
	if result.RowsAffected == 0 {
		return ErrCustomerNotFound
	}
	return result.Error
}

// SearchPaged searches customers by name, email, or phone (PRD section 11.3).
func (r *GormCustomerRepository) SearchPaged(query string, limit, offset int) ([]models.Customer, int64, error) {
	q := func() *gorm.DB {
		return r.db.Model(&models.Customer{}).
			Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")
	}

	var total int64
	if err := q().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var customers []models.Customer
	err := q().
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&customers).
		Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}
