package repositories

import (
	"github.com/fazasuny/erp-system/internal/models"
	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	Create(log *models.ActivityLog) error
	List(page, limit int, filters map[string]interface{}) ([]*models.ActivityLog, int64, error)
}

type GormActivityLogRepository struct {
	db *gorm.DB
}

func NewGormActivityLogRepository(db *gorm.DB) *GormActivityLogRepository {
	return &GormActivityLogRepository{db: db}
}

func (r *GormActivityLogRepository) Create(log *models.ActivityLog) error {
	return r.db.Create(log).Error
}

func (r *GormActivityLogRepository) List(page, limit int, filters map[string]interface{}) ([]*models.ActivityLog, int64, error) {
	var logs []*models.ActivityLog
	var total int64

	query := r.db.Model(&models.ActivityLog{})

	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if module, ok := filters["module"]; ok {
		query = query.Where("module = ?", module)
	}
	if action, ok := filters["action"]; ok {
		query = query.Where("action = ?", action)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error

	return logs, total, err
}
