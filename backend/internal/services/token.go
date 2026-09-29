package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims is the JWT payload carried by every authenticated request.
type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager issues and verifies the JWTs described in PRD section 5.
type TokenManager struct {
	secret    []byte
	expiresIn time.Duration
}

func NewTokenManager(secret string, expiresIn time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), expiresIn: expiresIn}
}

func (m *TokenManager) Generate(userID uint, email, role string) (string, error) {
	now := time.Now()

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiresIn)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

// Parse validates the signature and expiry, returning the token claims.
func (m *TokenManager) Parse(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid || claims.UserID == 0 {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
