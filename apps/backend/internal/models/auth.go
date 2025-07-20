package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// User represents a Discord user
type User struct {
	ID            string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DiscordUserID string         `json:"discord_user_id" gorm:"uniqueIndex;not null"`
	Username      string         `json:"username" gorm:"not null"`
	Avatar        string         `json:"avatar"`
	Email         string         `json:"email"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Tenants  []UserTenant `json:"tenants,omitempty" gorm:"foreignKey:UserID"`
	Sessions []Session    `json:"sessions,omitempty" gorm:"foreignKey:UserID"`
}

// Session represents a user session
type Session struct {
	ID                  string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID              string         `json:"user_id" gorm:"not null;index"`
	AccessToken         string         `json:"access_token" gorm:"not null"`
	RefreshToken        string         `json:"refresh_token" gorm:"not null;uniqueIndex"`
	DiscordAccessToken  string         `json:"discord_access_token" gorm:"not null"`
	DiscordRefreshToken string         `json:"discord_refresh_token" gorm:"not null"`
	ExpiresAt           time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID        string   `json:"user_id"`
	DiscordUserID string   `json:"discord_user_id"`
	Username      string   `json:"username"`
	SessionID     string   `json:"session_id"`
	SystemRoles   []string `json:"system_roles,omitempty"`
	jwt.RegisteredClaims
}

// DiscordUser represents Discord user data from API
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
	Discriminator string `json:"discriminator"`
}

// DiscordGuild represents a Discord guild/server from API
type DiscordGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions string `json:"permissions"`
	Features    []string `json:"features"`
}

// DiscordRole represents a Discord role from API
type DiscordRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int    `json:"position"`
	Permissions string `json:"permissions"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`
}

// DiscordMember represents a Discord guild member from API
type DiscordMember struct {
	User         *DiscordUser `json:"user"`
	Nick         string       `json:"nick"`
	Avatar       string       `json:"avatar"`
	Roles        []string     `json:"roles"`
	JoinedAt     string       `json:"joined_at"`
	PremiumSince string       `json:"premium_since"`
	Deaf         bool         `json:"deaf"`
	Mute         bool         `json:"mute"`
	Pending      bool         `json:"pending"`
	Permissions  string       `json:"permissions"`
}

// DiscordTokenResponse represents Discord OAuth2 token response
type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         User   `json:"user"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// APIError represents API error response
type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}