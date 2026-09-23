package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success responds with the PRD section 24 success envelope:
// {"data": {}, "message": "success"}
func Success(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, gin.H{
		"data":    data,
		"message": message,
	})
}

// Collection responds with the PRD section 24 collection envelope.
func Collection(c *gin.Context, data interface{}, page, limit int, total int64) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// Error responds with the PRD section 24 error envelope:
// {"error": {"code": "...", "message": "..."}}
func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
