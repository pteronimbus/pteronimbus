package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// JWTService handles JWT token operations
type JWTService struct {
	secret           []byte
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	issuer           string
	rbacService      *RBACService
}

// NewJWTService creates a new JWT service
func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		secret:           []byte(cfg.JWT.Secret),
		accessTokenTTL:   cfg.JWT.AccessTokenTTL,
		refreshTokenTTL:  cfg.JWT.RefreshTokenTTL,
		issuer:           cfg.JWT.Issuer,
	}
}

// NewJWTServiceWithRBAC creates a new JWT service with RBAC integration
func NewJWTServiceWithRBAC(cfg *config.Config, rbacService *RBACService) *JWTService {
	return &JWTService{
		secret:           []byte(cfg.JWT.Secret),
		accessTokenTTL:   cfg.JWT.AccessTokenTTL,
		refreshTokenTTL:  cfg.JWT.RefreshTokenTTL,
		issuer:           cfg.JWT.Issuer,
		rbacService:      rbacService,
	}
}

// GenerateAccessToken generates a new access token
func (j *JWTService) GenerateAccessToken(user *models.User, sessionID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessTokenTTL)
	
	// Debug output
	fmt.Printf("JWT Debug: now=%v, expiresAt=%v, TTL=%v\n", now, expiresAt, j.accessTokenTTL)
	
	// Get user's system roles if RBAC service is available
	var systemRoles []string
	if j.rbacService != nil {
		// Get user's system roles
		userSystemRoles, err := j.rbacService.GetUserSystemRoles(context.Background(), user.ID)
		if err != nil {
			// Log error but don't fail token generation
			fmt.Printf("Warning: failed to get system roles for user %s: %v\n", user.ID, err)
		} else {
			systemRoles = make([]string, len(userSystemRoles))
			for i, role := range userSystemRoles {
				systemRoles[i] = role.Name
			}
		}
	}
	
	claims := &models.JWTClaims{
		UserID:        user.ID,
		DiscordUserID: user.DiscordUserID,
		Username:      user.Username,
		SessionID:     sessionID,
		SystemRoles:   systemRoles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a new refresh token
func (j *JWTService) GenerateRefreshToken(user *models.User, sessionID string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenTTL)
	
	// Get user's system roles if RBAC service is available
	var systemRoles []string
	if j.rbacService != nil {
		// Get user's system roles
		userSystemRoles, err := j.rbacService.GetUserSystemRoles(context.Background(), user.ID)
		if err != nil {
			// Log error but don't fail token generation
			fmt.Printf("Warning: failed to get system roles for user %s: %v\n", user.ID, err)
		} else {
			systemRoles = make([]string, len(userSystemRoles))
			for i, role := range userSystemRoles {
				systemRoles[i] = role.Name
			}
		}
	}
	
	claims := &models.JWTClaims{
		UserID:        user.ID,
		DiscordUserID: user.DiscordUserID,
		Username:      user.Username,
		SessionID:     sessionID,
		SystemRoles:   systemRoles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetAccessTokenTTL returns the access token TTL in seconds
func (j *JWTService) GetAccessTokenTTL() int64 {
	return int64(j.accessTokenTTL.Seconds())
}