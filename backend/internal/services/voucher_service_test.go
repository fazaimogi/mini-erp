package services

import (
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

func setupVoucherService() VoucherService {
	return NewVoucherService(repositories.NewInMemoryVoucherRepository())
}

func pctVoucher(t *testing.T, svc VoucherService, value string) *models.Voucher {
	t.Helper()
	v, err := svc.Create(&VoucherInput{
		Code:          "promo10",
		DiscountType:  models.DiscountTypePercentage,
		DiscountValue: decimal.RequireFromString(value),
	})
	if err != nil {
		t.Fatalf("seed voucher: %v", err)
	}
	return v
}

func TestCreateVoucherNormalizesCodeAndDefaultsActive(t *testing.T) {
	svc := setupVoucherService()
	v := pctVoucher(t, svc, "10")

	if v.Code != "PROMO10" {
		t.Errorf("expected uppercased code, got %s", v.Code)
	}
	if !v.IsActive {
		t.Error("expected new voucher active by default")
	}
}

func TestCreateVoucherRejectsDuplicateCodeCaseInsensitive(t *testing.T) {
	svc := setupVoucherService()
	pctVoucher(t, svc, "10")

	_, err := svc.Create(&VoucherInput{
		Code:          "PROMO10",
		DiscountType:  models.DiscountTypePercentage,
		DiscountValue: decimal.NewFromInt(5),
	})
	if !errors.Is(err, ErrDuplicateVoucherCode) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestCreateVoucherRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name  string
		input VoucherInput
	}{
		{"empty code", VoucherInput{Code: "  ", DiscountType: models.DiscountTypeFixed, DiscountValue: decimal.NewFromInt(5)}},
		{"bad type", VoucherInput{Code: "X", DiscountType: "bogus", DiscountValue: decimal.NewFromInt(5)}},
		{"zero value", VoucherInput{Code: "X", DiscountType: models.DiscountTypeFixed, DiscountValue: decimal.Zero}},
		{"pct over 100", VoucherInput{Code: "X", DiscountType: models.DiscountTypePercentage, DiscountValue: decimal.NewFromInt(101)}},
		{"negative limit", VoucherInput{Code: "X", DiscountType: models.DiscountTypeFixed, DiscountValue: decimal.NewFromInt(5), UsageLimit: -1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := setupVoucherService()
			if _, err := svc.Create(&tc.input); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestValidatePercentageDiscount(t *testing.T) {
	svc := setupVoucherService()
	pctVoucher(t, svc, "10")

	result, err := svc.Validate("promo10", decimal.NewFromInt(250))
	if err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
	if !result.DiscountAmount.Equal(decimal.NewFromInt(25)) {
		t.Errorf("expected 25 discount, got %s", result.DiscountAmount)
	}
	if !result.FinalAmount.Equal(decimal.NewFromInt(225)) {
		t.Errorf("expected 225 final, got %s", result.FinalAmount)
	}
}

func TestValidateFixedDiscountClampsToTotal(t *testing.T) {
	svc := setupVoucherService()
	_, err := svc.Create(&VoucherInput{
		Code:          "FLAT50",
		DiscountType:  models.DiscountTypeFixed,
		DiscountValue: decimal.NewFromInt(50),
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	// total 30 < discount 50 => discount clamped to 30, final 0 (never negative)
	result, err := svc.Validate("FLAT50", decimal.NewFromInt(30))
	if err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
	if !result.DiscountAmount.Equal(decimal.NewFromInt(30)) {
		t.Errorf("expected clamped discount 30, got %s", result.DiscountAmount)
	}
	if !result.FinalAmount.Equal(decimal.Zero) {
		t.Errorf("expected zero final, got %s", result.FinalAmount)
	}
}

func TestValidateRejectsMinPurchaseNotMet(t *testing.T) {
	svc := setupVoucherService()
	_, err := svc.Create(&VoucherInput{
		Code:          "BIGSPEND",
		DiscountType:  models.DiscountTypeFixed,
		DiscountValue: decimal.NewFromInt(20),
		MinPurchase:   decimal.NewFromInt(100),
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	_, err = svc.Validate("BIGSPEND", decimal.NewFromInt(50))
	if !errors.Is(err, ErrVoucherMinPurchase) {
		t.Fatalf("expected min purchase error, got %v", err)
	}
}

func TestValidateRejectsExpiredInactiveAndExhausted(t *testing.T) {
	svc := setupVoucherService()

	past := time.Now().Add(-time.Hour)
	expired, _ := svc.Create(&VoucherInput{
		Code: "OLD", DiscountType: models.DiscountTypeFixed,
		DiscountValue: decimal.NewFromInt(5), ExpiresAt: &past,
	})
	if _, err := svc.Validate(expired.Code, decimal.NewFromInt(100)); !errors.Is(err, ErrVoucherExpired) {
		t.Errorf("expected expired error, got %v", err)
	}

	off, _ := svc.Create(&VoucherInput{
		Code: "OFF", DiscountType: models.DiscountTypeFixed,
		DiscountValue: decimal.NewFromInt(5), IsActive: boolPtr(false),
	})
	if _, err := svc.Validate(off.Code, decimal.NewFromInt(100)); !errors.Is(err, ErrVoucherInactive) {
		t.Errorf("expected inactive error, got %v", err)
	}

	repo := repositories.NewInMemoryVoucherRepository()
	limited := NewVoucherService(repo)
	seed, _ := limited.Create(&VoucherInput{
		Code: "ONCE", DiscountType: models.DiscountTypeFixed,
		DiscountValue: decimal.NewFromInt(5), UsageLimit: 1,
	})
	seed.UsedCount = 1
	if err := repo.Update(seed); err != nil {
		t.Fatalf("update: %v", err)
	}
	if _, err := limited.Validate("ONCE", decimal.NewFromInt(100)); !errors.Is(err, ErrVoucherUsageLimit) {
		t.Errorf("expected usage limit error, got %v", err)
	}
}

func TestValidateUnknownCodeReturnsNotFound(t *testing.T) {
	svc := setupVoucherService()
	if _, err := svc.Validate("NOPE", decimal.NewFromInt(10)); !errors.Is(err, ErrVoucherNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestUpdateRejectsCodeTakenByAnotherVoucher(t *testing.T) {
	svc := setupVoucherService()
	first := pctVoucher(t, svc, "10")
	second, _ := svc.Create(&VoucherInput{
		Code: "SECOND", DiscountType: models.DiscountTypeFixed, DiscountValue: decimal.NewFromInt(5),
	})

	_, err := svc.Update(second.ID, &VoucherInput{
		Code: first.Code, DiscountType: models.DiscountTypeFixed, DiscountValue: decimal.NewFromInt(5),
	})
	if !errors.Is(err, ErrDuplicateVoucherCode) {
		t.Fatalf("expected duplicate, got %v", err)
	}
}

func TestDeleteRemovesVoucher(t *testing.T) {
	svc := setupVoucherService()
	v := pctVoucher(t, svc, "10")

	if err := svc.Delete(v.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.GetByID(v.ID); !errors.Is(err, ErrVoucherNotFound) {
		t.Fatalf("expected gone, got %v", err)
	}
}

func boolPtr(b bool) *bool { return &b }