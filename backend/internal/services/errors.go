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

	ErrUserNotFound       = NewAppError("USER_NOT_FOUND", "user not found", http.StatusNotFound)
	ErrEmailAlreadyExists = NewAppError("EMAIL_ALREADY_EXISTS", "email already registered", http.StatusConflict)
	ErrInvalidCredentials = NewAppError("INVALID_CREDENTIALS", "invalid email or password", http.StatusUnauthorized)
	ErrAccountInactive    = NewAppError("ACCOUNT_INACTIVE", "account is inactive", http.StatusForbidden)
	ErrUnauthorized       = NewAppError("UNAUTHORIZED", "authentication required", http.StatusUnauthorized)
	ErrInvalidToken       = NewAppError("INVALID_TOKEN", "invalid or expired token", http.StatusUnauthorized)
	ErrForbidden          = NewAppError("FORBIDDEN", "insufficient permissions", http.StatusForbidden)

	// INSUFFICIENT_STOCK is the exact code PRD section 14.1 specifies.
	ErrInsufficientStock = NewAppError("INSUFFICIENT_STOCK", "insufficient stock", http.StatusConflict)
	ErrInvalidQuantity   = NewAppError("INVALID_QUANTITY", "quantity must be greater than zero", http.StatusUnprocessableEntity)
	ErrInvalidPeriod     = NewAppError("INVALID_PERIOD", "invalid period parameter", http.StatusUnprocessableEntity)
	ErrInvalidDateRange  = NewAppError("INVALID_DATE_RANGE", "custom period requires start_date and end_date", http.StatusUnprocessableEntity)

	ErrVoucherNotFound       = NewAppError("VOUCHER_NOT_FOUND", "voucher not found", http.StatusNotFound)
	ErrDuplicateVoucherCode  = NewAppError("DUPLICATE_VOUCHER_CODE", "voucher code already exists", http.StatusConflict)
	ErrVoucherInactive       = NewAppError("VOUCHER_INACTIVE", "voucher is inactive", http.StatusUnprocessableEntity)
	ErrVoucherExpired        = NewAppError("VOUCHER_EXPIRED", "voucher has expired", http.StatusUnprocessableEntity)
	ErrVoucherUsageLimit     = NewAppError("VOUCHER_USAGE_LIMIT_REACHED", "voucher usage limit reached", http.StatusUnprocessableEntity)
	ErrVoucherMinPurchase    = NewAppError("VOUCHER_MIN_PURCHASE_NOT_MET", "minimum purchase not met for this voucher", http.StatusUnprocessableEntity)
)
