package services

import (
	"net/http"
	"strings"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

// CustomerInput defines the accepted payload for create and update operations.
// PRD section 11: name is required, email must be unique.
type CustomerInput struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type CustomerService interface {
	List(page, limit int) ([]models.Customer, int64, error)
	GetByID(id uint) (*models.Customer, error)
	Create(input *CustomerInput) (*models.Customer, error)
	Update(id uint, input *CustomerInput) (*models.Customer, error)
	Delete(id uint) error
	Search(query string, page, limit int) ([]models.Customer, int64, error)
}

type customerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) List(page, limit int) ([]models.Customer, int64, error) {
	return s.repo.ListPaged(limit, (page-1)*limit)
}

func (s *customerService) GetByID(id uint) (*models.Customer, error) {
	customer, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrCustomerNotFound {
			return nil, NewAppError("CUSTOMER_NOT_FOUND", "customer not found", http.StatusNotFound)
		}
		return nil, err
	}
	return customer, nil
}

func (s *customerService) Create(input *CustomerInput) (*models.Customer, error) {
	if err := validateCustomerInput(input); err != nil {
		return nil, err
	}

	if err := s.ensureUniqueEmail(input.Email, 0); err != nil {
		return nil, err
	}

	customer := &models.Customer{
		Name:    strings.TrimSpace(input.Name),
		Email:   strings.TrimSpace(input.Email),
		Phone:   strings.TrimSpace(input.Phone),
		Address: strings.TrimSpace(input.Address),
	}

	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) Update(id uint, input *CustomerInput) (*models.Customer, error) {
	if err := validateCustomerInput(input); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(id)
	if err != nil {
		if err == repositories.ErrCustomerNotFound {
			return nil, NewAppError("CUSTOMER_NOT_FOUND", "customer not found", http.StatusNotFound)
		}
		return nil, err
	}

	// Only check email uniqueness if email changed
	if strings.TrimSpace(input.Email) != existing.Email {
		if err := s.ensureUniqueEmail(input.Email, id); err != nil {
			return nil, err
		}
	}

	existing.Name = strings.TrimSpace(input.Name)
	existing.Email = strings.TrimSpace(input.Email)
	existing.Phone = strings.TrimSpace(input.Phone)
	existing.Address = strings.TrimSpace(input.Address)

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *customerService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *customerService) Search(query string, page, limit int) ([]models.Customer, int64, error) {
	return s.repo.SearchPaged(query, limit, (page-1)*limit)
}

func (s *customerService) ensureUniqueEmail(email string, excludeID uint) error {
	if email == "" {
		return nil
	}

	existing, err := s.repo.GetByEmail(email)
	if err != nil {
		return err
	}

	if existing != nil && existing.ID != excludeID {
		return NewAppError("DUPLICATE_EMAIL", "email already exists", http.StatusConflict)
	}

	return nil
}

func validateCustomerInput(input *CustomerInput) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return NewAppError("INVALID_INPUT", "name is required", http.StatusUnprocessableEntity)
	}

	return nil
}
