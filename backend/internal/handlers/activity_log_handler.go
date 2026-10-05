package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

type ActivityLogHandler struct {
	service services.ActivityLogService
}

func NewActivityLogHandler(service services.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{service: service}
}

// List godoc
// GET /api/activity-logs?page=1&limit=20&user_id=1&module=Products&action=CREATE
// Admin only. Returns paginated list of activity logs with optional filters.
func (h *ActivityLogHandler) List(c *gin.Context) {
	page, limit := pagination(c)

	filters := make(map[string]interface{})
	
	if userID := c.Query("user_id"); userID != "" {
		if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
			filters["user_id"] = uint(id)
		}
	}
	
	if module := c.Query("module"); module != "" {
		filters["module"] = module
	}
	
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}

	logs, total, err := h.service.List(page, limit, filters)
	if err != nil {
		respondInternalError(c, err)
		return
	}

	Collection(c, logs, page, limit, total)
}
