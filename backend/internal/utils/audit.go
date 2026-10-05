package utils

import (
	"github.com/fazasuny/erp-system/internal/services"
	"github.com/gin-gonic/gin"
)

// AuditLogger wraps ActivityLogService untuk kemudahan logging
type AuditLogger struct {
	service services.ActivityLogService
}

func NewAuditLogger(service services.ActivityLogService) *AuditLogger {
	return &AuditLogger{service: service}
}

// LogFromContext logs activity dari Gin context (ambil user info dari JWT)
func (a *AuditLogger) LogFromContext(c *gin.Context, action, module, description string) {
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")
	
	var uid *uint
	if id, ok := userID.(uint); ok {
		uid = &id
	}
	
	name := "Unknown"
	if n, ok := userName.(string); ok && n != "" {
		name = n
	}
	
	ip := c.ClientIP()
	
	_ = a.service.Log(uid, name, action, module, description, ip)
}

// Log logs activity langsung tanpa context
func (a *AuditLogger) Log(userID *uint, userName, action, module, description, ipAddress string) {
	_ = a.service.Log(userID, userName, action, module, description, ipAddress)
}
