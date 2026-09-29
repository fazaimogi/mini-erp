package services

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

func setupAuthService() (AuthService, repositories.UserRepository, *TokenManager) {
	users := repositories.NewInMemoryUserRepository()
	tokens := NewTokenManager("test-secret", time.Hour)
	return NewAuthService(users, tokens), users, tokens
}

func validRegisterInput() *RegisterInput {
	return &RegisterInput{
		Name:     "Budi Santoso",
		Email:    "Budi@Example.com ",
		Password: "secret12345",
	}
}

func TestRegisterHashesPasswordAndAssignsStaffRole(t *testing.T) {
	service, _, _ := setupAuthService()

	user, err := service.Register(validRegisterInput())
	if err != nil {
		t.Fatalf("expected register success, got: %v", err)
	}

	if user.PasswordHash == "secret12345" || user.PasswordHash == "" {
		t.Fatal("password stored in plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("secret12345")); err != nil {
		t.Errorf("stored hash does not match password: %v", err)
	}
	if user.Email != "budi@example.com" {
		t.Errorf("expected normalised email, got: %s", user.Email)
	}
	if user.Role != "staff" || user.RoleID == 0 {
		t.Errorf("expected staff role, got role=%q role_id=%d", user.Role, user.RoleID)
	}
	if user.Status != "active" {
		t.Errorf("expected active status, got: %s", user.Status)
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	service, _, _ := setupAuthService()

	if _, err := service.Register(validRegisterInput()); err != nil {
		t.Fatalf("seed register failed: %v", err)
	}

	duplicate := validRegisterInput()
	duplicate.Name = "Someone Else"
	duplicate.Email = "budi@example.com"

	_, err := service.Register(duplicate)

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got: %v", err)
	}
	if appErr.Code != "EMAIL_ALREADY_EXISTS" || appErr.HTTPStatus != 409 {
		t.Errorf("unexpected error: %+v", appErr)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	cases := map[string]func(*RegisterInput){
		"blank name":     func(in *RegisterInput) { in.Name = "   " },
		"invalid email":  func(in *RegisterInput) { in.Email = "not-an-email" },
		"short password": func(in *RegisterInput) { in.Password = "short" },
		"long password":  func(in *RegisterInput) { in.Password = string(make([]byte, 73)) },
		"overlong name":  func(in *RegisterInput) { in.Name = string(make([]byte, 256)) },
	}

	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			service, _, _ := setupAuthService()

			input := validRegisterInput()
			mutate(input)

			_, err := service.Register(input)

			var appErr *AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got: %v", err)
			}
			if appErr.HTTPStatus != 422 {
				t.Errorf("expected 422, got %d", appErr.HTTPStatus)
			}
		})
	}
}

func TestLoginReturnsVerifiableToken(t *testing.T) {
	service, _, tokens := setupAuthService()

	if _, err := service.Register(validRegisterInput()); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user, token, err := service.Login(&LoginInput{Email: "BUDI@example.com", Password: "secret12345"})
	if err != nil {
		t.Fatalf("expected login success, got: %v", err)
	}
	if token == "" {
		t.Fatal("expected token")
	}

	claims, err := tokens.Parse(token)
	if err != nil {
		t.Fatalf("issued token failed verification: %v", err)
	}
	if claims.UserID != user.ID || claims.Email != "budi@example.com" || claims.Role != "staff" {
		t.Errorf("unexpected claims: %+v", claims)
	}
}

func TestLoginRejectsWrongPasswordAndUnknownEmail(t *testing.T) {
	service, _, _ := setupAuthService()

	if _, err := service.Register(validRegisterInput()); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	cases := map[string]*LoginInput{
		"wrong password": {Email: "budi@example.com", Password: "wrong-password"},
		"unknown email":  {Email: "ghost@example.com", Password: "secret12345"},
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, err := service.Login(input)

			var appErr *AppError
			if !errors.As(err, &appErr) {
				t.Fatalf("expected AppError, got: %v", err)
			}
			if appErr.Code != "INVALID_CREDENTIALS" || appErr.HTTPStatus != 401 {
				t.Errorf("unexpected error: %+v", appErr)
			}
		})
	}
}

// inactiveUserRepo returns an inactive account so the login guard can be
// exercised; the in-memory repository has no update path and PRD section 15
// (user management) is not implemented yet.
type inactiveUserRepo struct {
	repositories.UserRepository
	user models.User
}

func (r *inactiveUserRepo) GetByEmail(string) (*models.User, error) {
	found := r.user
	return &found, nil
}

func TestLoginRejectsInactiveAccount(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret12345"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	repo := &inactiveUserRepo{
		UserRepository: repositories.NewInMemoryUserRepository(),
		user: models.User{
			ID:           3,
			Email:        "budi@example.com",
			PasswordHash: string(hash),
			RoleID:       2,
			Role:         "staff",
			Status:       "inactive",
		},
	}

	service := NewAuthService(repo, NewTokenManager("test-secret", time.Hour))

	_, _, err = service.Login(&LoginInput{Email: "budi@example.com", Password: "secret12345"})

	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got: %v", err)
	}
	if appErr.Code != "ACCOUNT_INACTIVE" || appErr.HTTPStatus != 403 {
		t.Errorf("unexpected error: %+v", appErr)
	}
}

func TestMeReturnsUserAndRole(t *testing.T) {
	service, _, _ := setupAuthService()

	registered, err := service.Register(validRegisterInput())
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user, err := service.Me(registered.ID)
	if err != nil {
		t.Fatalf("expected me success, got: %v", err)
	}
	if user.ID != registered.ID || user.Role != "staff" {
		t.Errorf("unexpected user: %+v", user)
	}

	if _, err := service.Me(9999); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestTokenManagerRejectsTamperedAndExpiredTokens(t *testing.T) {
	tokens := NewTokenManager("test-secret", time.Hour)

	token, err := tokens.Generate(7, "budi@example.com", "staff")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	if _, err := tokens.Parse(token + "x"); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken for tampered token, got: %v", err)
	}

	other := NewTokenManager("different-secret", time.Hour)
	if _, err := other.Parse(token); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken for foreign secret, got: %v", err)
	}

	expired := NewTokenManager("test-secret", -time.Minute)
	expiredToken, err := expired.Generate(7, "budi@example.com", "staff")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if _, err := tokens.Parse(expiredToken); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("expected ErrInvalidToken for expired token, got: %v", err)
	}
}
