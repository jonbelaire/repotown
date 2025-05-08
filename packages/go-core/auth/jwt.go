package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Common errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// TokenClaims represents the claims in a JWT token
type TokenClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

// JWTService provides methods for JWT token handling
type JWTService struct {
	secretKey     []byte
	tokenDuration time.Duration
	issuer        string
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, tokenDuration time.Duration, issuer string) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		tokenDuration: tokenDuration,
		issuer:        issuer,
	}
}

// GenerateToken creates a new JWT token for a user
func (s *JWTService) GenerateToken(userID, userRole string) (string, error) {
	tokenID := uuid.New().String()
	now := time.Now()
	
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenDuration)),
		},
		UserID:   userID,
		UserRole: userRole,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates a token and returns its claims
func (s *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secretKey, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}