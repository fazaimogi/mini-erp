package services

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

// UserInput defines the accepted payload for create and update operations.
// PRD section 15.2: name, email, role_id required; password for create only.
type UserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   uint   `json:"role_id" binding:"required"`
	Status   string `json:"status"`
}

// UserUpdateInput for PATCH/PUT without requiring password.
type UserUpdateInput struct {
	Name   string `json:"name"`
	Email  string `json:"email" binding:"email"`
	RoleID uint   `json:"role_id"`
	Status string `json:"status"`
}

type UserService interface {
	List(page, limit int) ([]models.User, int64, error)
	GetByID(id uint) (*models.User, error)
	Create(input *UserInput) (*models.User, error)
	Update(id uint, input *UserUpdateInput) (*models.User, error)
	Delete(id uint) error
	ChangeRole(id uint, roleID uint) (*models.User, error)
	ChangeStatus(id uint, status string) (*models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) List(page, limit int) ([]models.User, int64, error) {
	return s.repo.List(limit, (page-1)*limit)
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if err.Error() == repositories.ErrUserNotFound.Error() {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Create(input *UserInput) (*models.User, error) {
	if err := validateUserInput(input); err != nil {
		return nil, err
	}

	// Check email not already used
	existing, _ := s.repo.GetByEmail(input.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, NewAppError("HASH_ERROR", "failed to hash password", 500)
	}

	user := &models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		RoleID:       input.RoleID,
		Status:       "active",
	}

	if err := s.repo.Create(user); err != nil {
		if err.Error() == repositories.ErrDuplicateEmail.Error() {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(id uint, input *UserUpdateInput) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if err.Error() == repositories.ErrUserNotFound.Error() {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.RoleID != 0 {
		user.RoleID = input.RoleID
	}
	if input.Status != "" {
		if input.Status != "active" && input.Status != "inactive" {
			return nil, NewAppError("INVALID_STATUS", "status must be 'active' or 'inactive'", 422)
		}
		user.Status = input.Status
	}

	if err := s.repo.Update(user); err != nil {
		if err.Error() == repositories.ErrDuplicateEmail.Error() {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id uint) error {
	// Check user exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err.Error() == repositories.ErrUserNotFound.Error() {
			return ErrUserNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *userService) ChangeRole(id uint, roleID uint) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if err.Error() == repositories.ErrUserNotFound.Error() {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.RoleID = roleID

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ChangeStatus(id uint, status string) (*models.User, error) {
	if status != "active" && status != "inactive" {
		return nil, NewAppError("INVALID_STATUS", "status must be 'active' or 'inactive'", 422)
	}

	user, err := s.repo.GetByID(id)
	if err != nil {
		if err.Error() == repositories.ErrUserNotFound.Error() {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Status = status

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func validateUserInput(input *UserInput) error {
	if input.Name == "" {
		return ErrInvalidInput
	}
	if input.Email == "" {
		return ErrInvalidInput
	}
	if input.Password == "" || len(input.Password) < 6 {
		return NewAppError("INVALID_PASSWORD", "password must be at least 6 characters", 422)
	}
	if input.RoleID == 0 {
		return ErrInvalidInput
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
