package repositories

import (
	"strings"
	"sync"

	"github.com/fazasuny/erp-system/internal/models"
)

// InMemoryUserRepository backs the auth unit tests. It mirrors the role seed
// from migration 007 so register can resolve "staff" without a database.
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	nextID uint
	users  []models.User
	roles  []models.Role
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		nextID: 1,
		roles: []models.Role{
			{ID: 1, Name: "admin"},
			{ID: 2, Name: "staff"},
		},
	}
}

func (r *InMemoryUserRepository) Create(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.users {
		if strings.EqualFold(existing.Email, user.Email) {
			return ErrDuplicateEmail
		}
	}

	user.ID = r.nextID
	r.nextID++
	r.users = append(r.users, *user)

	return nil
}

func (r *InMemoryUserRepository) GetByEmail(email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			found := user
			return &found, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) GetByID(id uint) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			found := user
			return &found, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) GetRoleByID(id uint) (*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, role := range r.roles {
		if role.ID == id {
			found := role
			return &found, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) GetRoleByName(name string) (*models.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, role := range r.roles {
		if role.Name == name {
			found := role
			return &found, nil
		}
	}

	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}

	return ErrUserNotFound
}

func (r *InMemoryUserRepository) Update(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, existing := range r.users {
		if existing.ID == user.ID {
			r.users[i] = *user
			return nil
		}
	}

	return ErrUserNotFound
}

func (r *InMemoryUserRepository) List(limit, offset int) ([]models.User, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := int64(len(r.users))
	start := offset
	end := offset + limit
	if end > len(r.users) {
		end = len(r.users)
	}
	if start > len(r.users) {
		start = len(r.users)
	}

	result := make([]models.User, 0)
	if start < end {
		result = r.users[start:end]
	}
	return result, total, nil
}
