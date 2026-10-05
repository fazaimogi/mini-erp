package services

import (
	"fmt"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

type ActivityLogService interface {
	Log(userID *uint, userName, action, module, description, ipAddress string) error
	List(page, limit int, filters map[string]interface{}) ([]*models.ActivityLog, int64, error)
}

type activityLogService struct {
	repo repositories.ActivityLogRepository
}

func NewActivityLogService(repo repositories.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

func (s *activityLogService) Log(userID *uint, userName, action, module, description, ipAddress string) error {
	log := &models.ActivityLog{
		UserID:      userID,
		UserName:    userName,
		Action:      action,
		Module:      module,
		Description: description,
		IPAddress:   ipAddress,
	}
	return s.repo.Create(log)
}

func (s *activityLogService) List(page, limit int, filters map[string]interface{}) ([]*models.ActivityLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	logs, total, err := s.repo.List(page, limit, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list activity logs: %w", err)
	}

	return logs, total, nil
}
