package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
)

// ErrCustomerNotFound is returned when no customer matches the given id.
var ErrCustomerNotFound = errors.New("customer not found")

// CustomerRepository abstracts customer persistence.
type CustomerRepository interface {
	ListPaged(limit, offset int) ([]models.Customer, int64, error)
	GetByID(id uint) (*models.Customer, error)
	GetByEmail(email string) (*models.Customer, error)
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
	Delete(id uint) error
	SearchPaged(query string, limit, offset int) ([]models.Customer, int64, error)
}
