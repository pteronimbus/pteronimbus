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
	GetUserGuilds(ctx context.Context, accessToken string) ([]models.DiscordGuild, error)
	GetGuildRoles(ctx context.Context, botToken, guildID string) ([]models.DiscordRole, error)
	GetGuildMembers(ctx context.Context, botToken, guildID string, limit int) ([]models.DiscordMember, error)
	GetGuildMember(ctx context.Context, botToken, guildID, userID string) (*models.DiscordMember, error)
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
	ParseTokenClaims(accessToken string) (*models.JWTClaims, error)
	Logout(ctx context.Context, accessToken string) error
}

// TenantServiceInterface defines the interface for tenant service operations
type TenantServiceInterface interface {
	CreateTenant(ctx context.Context, discordGuild *models.DiscordGuild, ownerID string) (*models.Tenant, error)
	GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error)
	GetTenantByDiscordServerID(ctx context.Context, discordServerID string) (*models.Tenant, error)
	GetUserTenants(ctx context.Context, userID string) ([]models.Tenant, error)
	AddUserToTenant(ctx context.Context, userID, tenantID string, roles []string, permissions []string) error
	RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error
	SyncDiscordRoles(ctx context.Context, tenantID, botToken string) error
	SyncDiscordUsers(ctx context.Context, tenantID, botToken string) error
	UpdateTenantConfig(ctx context.Context, tenantID string, config models.TenantConfig) error
	DeleteTenant(ctx context.Context, tenantID string) error
	HasPermission(ctx context.Context, userID, tenantID, permission string) (bool, error)
	CheckManageServerPermission(discordGuild *models.DiscordGuild) bool
}