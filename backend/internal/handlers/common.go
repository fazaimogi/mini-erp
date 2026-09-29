package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fazasuny/erp-system/internal/services"
)

const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

// pagination reads the PRD section 24 collection params and clamps them.
func pagination(c *gin.Context) (int, int) {
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	limit := parsePositiveInt(c.DefaultQuery("limit", strconv.Itoa(defaultPageLimit)), defaultPageLimit)
	if limit > maxPageLimit {
		limit = defaultPageLimit
	}
	return page, limit
}

func parseIDParam(c *gin.Context) (uint, error) {
	return parseUint(c.Param("id"))
}

// parseUint rejects zero so a caller cannot address a non-existent row id 0.
func parseUint(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func parsePositiveInt(value string, fallback int) int {
	number, err := strconv.Atoi(value)
	if err != nil || number < 1 {
		return fallback
	}
	return number
}

func respondServiceError(c *gin.Context, err error) {
	var appErr *services.AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
		return
	}
	respondInternalError(c, err)
}

func respondInternalError(c *gin.Context, err error) {
	// Recorded so gin's logger prints the cause; the client still only sees the
	// generic PRD section 24 envelope.
	_ = c.Error(err)
	Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
}
