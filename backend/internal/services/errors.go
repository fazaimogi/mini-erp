package services

import (
	"fmt"
	"net/http"
)

// AppError carries an error code, human readable message and the HTTP status
// that fits the PRD section 24 error response standard.
type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAppError(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

var (
	ErrProductNotFound = NewAppError("PRODUCT_NOT_FOUND", "product not found", http.StatusNotFound)
	ErrDuplicateSKU    = NewAppError("DUPLICATE_SKU", "sku already exists", http.StatusConflict)
	ErrInvalidInput    = NewAppError("INVALID_INPUT", "invalid product input", http.StatusUnprocessableEntity)
)
