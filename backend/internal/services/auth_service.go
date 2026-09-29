package services

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/fazasuny/erp-system/internal/models"
	"github.com/fazasuny/erp-system/internal/repositories"
)

const (
	// Self-registration always lands on this role; promoting a user to admin is
	// an administrative action (PRD section 15), not a signup choice.
	defaultSignupRole = "staff"

	activeStatus = "active"

	minPasswordLength = 8
	// bcrypt silently truncates beyond 72 bytes, so reject instead of hashing
	// a password the user did not actually set.
	maxPasswordLength = 72
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// RegisterInput is the signup payload for POST /api/auth/register.
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginInput is the credential payload for POST /api/auth/login.
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthService interface {
	Register(input *RegisterInput) (*models.User, error)
	Login(input *LoginInput) (*models.User, string, error)
	Me(userID uint) (*models.User, error)
}

type authService struct {
	users  repositories.UserRepository
	tokens *TokenManager
}

func NewAuthService(users repositories.UserRepository, tokens *TokenManager) AuthService {
	return &authService{users: users, tokens: tokens}
}

func (s *authService) Register(input *RegisterInput) (*models.User, error) {
	email := normalizeEmail(input.Email)
	name := strings.TrimSpace(input.Name)

	if err := validateRegisterInput(name, email, input.Password); err != nil {
		return nil, err
	}

	role, err := s.users.GetRoleByName(defaultSignupRole)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		RoleID:       role.ID,
		Role:         role.Name,
		Status:       activeStatus,
	}

	if err := s.users.Create(user); err != nil {
		if err == repositories.ErrDuplicateEmail {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(input *LoginInput) (*models.User, string, error) {
	email := normalizeEmail(input.Email)

	user, err := s.users.GetByEmail(email)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			// Same error as a wrong password so the endpoint does not confirm
			// which emails are registered.
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if user.Status != activeStatus {
		return nil, "", ErrAccountInactive
	}

	if err := s.attachRole(user); err != nil {
		return nil, "", err
	}

	token, err := s.tokens.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Me(userID uint) (*models.User, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := s.attachRole(user); err != nil {
		return nil, err
	}

	return user, nil
}

// attachRole resolves role_id into the role name returned by the API, so the
// frontend can gate menus without a second request.
func (s *authService) attachRole(user *models.User) error {
	if user.Role != "" {
		return nil
	}

	role, err := s.users.GetRoleByID(user.RoleID)
	if err != nil {
		return err
	}

	user.Role = role.Name
	return nil
}

func validateRegisterInput(name, email, password string) error {
	if name == "" || len(name) > maxNameLength {
		return ErrInvalidInput
	}
	if !emailPattern.MatchString(email) || len(email) > maxNameLength {
		return ErrInvalidInput
	}
	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return ErrInvalidInput
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
