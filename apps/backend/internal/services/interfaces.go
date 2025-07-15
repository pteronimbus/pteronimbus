package services

import (
	"context"
	"time"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// DiscordServiceInterface defines the interface for Discord service operations
type DiscordServiceInterface interface {
	GetAuthURL(state string) string
	ExchangeCodeForToken(ctx context.Context, code string) (*models.DiscordTokenResponse, error)
	GetUserInfo(ctx context.Context, accessToken string) (*models.DiscordUser, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.DiscordTokenResponse, error)
}

// JWTServiceInterface defines the interface for JWT service operations
type JWTServiceInterface interface {
	GenerateAccessToken(user *models.User, sessionID string) (string, time.Time, error)
	GenerateRefreshToken(user *models.User, sessionID string) (string, time.Time, error)
	ValidateToken(tokenString string) (*models.JWTClaims, error)
}

// RedisServiceInterface defines the interface for Redis service operations
type RedisServiceInterface interface {
	StoreSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, sessionID string) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

// AuthServiceInterface defines the interface for authentication service operations
type AuthServiceInterface interface {
	GetAuthURL(state string) string
	HandleCallback(ctx context.Context, code string) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*models.AuthResponse, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (*models.User, error)
	Logout(ctx context.Context, accessToken string) error
}