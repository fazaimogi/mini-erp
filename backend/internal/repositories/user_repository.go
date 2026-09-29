package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
)

// ErrUserNotFound is returned when no user matches the given id or email.
var ErrUserNotFound = errors.New("user not found")

// ErrDuplicateEmail is returned when the email is already registered.
var ErrDuplicateEmail = errors.New("email already exists")

// UserRepository abstracts user persistence. Day 3 only needs lookup by id and
// email plus creation; full CRUD arrives with user management (PRD section 15).
type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	List(limit, offset int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
	GetRoleByName(name string) (*models.Role, error)
	GetRoleByID(id uint) (*models.Role, error)
}
