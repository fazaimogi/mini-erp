package repositories

import (
	"errors"

	"github.com/fazasuny/erp-system/internal/models"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicateEmail
	}
	return err
}

func (r *GormUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *GormUserRepository) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role

	err := r.db.First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &role, nil
}

func (r *GormUserRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role

	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &role, nil
}

func (r *GormUserRepository) List(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	err := r.db.Model(&models.User{}).Count(&total).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *GormUserRepository) Update(user *models.User) error {
	// Check for email uniqueness if changed, excluding current user
	var existing models.User
	err := r.db.Where("email = ? AND id != ?", user.Email, user.ID).First(&existing).Error
	if err == nil {
		// Email exists and belongs to different user
		return ErrDuplicateEmail
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return r.db.Model(user).Updates(user).Error
}

func (r *GormUserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
