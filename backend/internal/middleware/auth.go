package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

// Context keys populated by RequireAuth and read by downstream handlers.
const (
	ContextUserID = "user_id"
	ContextEmail  = "email"
	ContextRole   = "role"
)

// RequireAuth rejects requests without a valid Bearer token (PRD section 5).
func RequireAuth(tokens *services.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			abortError(c, services.ErrUnauthorized)
			return
		}

		scheme, token, found := strings.Cut(header, " ")
		if !found || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
			abortError(c, services.ErrUnauthorized)
			return
		}

		claims, err := tokens.Parse(strings.TrimSpace(token))
		if err != nil {
			abortError(c, services.ErrInvalidToken)
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextEmail, claims.Email)
		c.Set(ContextRole, claims.Role)

		c.Next()
	}
}

// RequireRole enforces PRD section 4: only the listed roles may proceed.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role := currentRole(c)
		if _, ok := allowed[role]; !ok {
			abortError(c, services.ErrForbidden)
			return
		}

		c.Next()
	}
}

// CurrentUserID returns the authenticated user id set by RequireAuth.
func CurrentUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(ContextUserID)
	if !exists {
		return 0, false
	}

	id, ok := value.(uint)
	return id, ok
}

func currentRole(c *gin.Context) string {
	value, exists := c.Get(ContextRole)
	if !exists {
		return ""
	}

	role, _ := value.(string)
	return role
}

// abortError writes the PRD section 24 error envelope directly instead of
// calling handlers.Error, because the handlers package imports this one for
// CurrentUserID.
func abortError(c *gin.Context, appErr *services.AppError) {
	c.AbortWithStatusJSON(appErr.HTTPStatus, gin.H{
		"error": gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		},
	})
}
