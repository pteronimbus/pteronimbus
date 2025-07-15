package services

import (
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

// GenerateAccessToken generates a new access token
func (j *JWTService) GenerateAccessToken(user *models.User, sessionID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(j.accessTokenTTL)
	
	claims := &models.JWTClaims{
		UserID:        user.ID,
		DiscordUserID: user.DiscordUserID,
		Username:      user.Username,
		SessionID:     sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
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
	expiresAt := time.Now().Add(j.refreshTokenTTL)
	
	claims := &models.JWTClaims{
		UserID:        user.ID,
		DiscordUserID: user.DiscordUserID,
		Username:      user.Username,
		SessionID:     sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
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