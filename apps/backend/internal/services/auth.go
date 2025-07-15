package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// AuthService handles authentication operations
type AuthService struct {
	discordService DiscordServiceInterface
	jwtService     JWTServiceInterface
	redisService   RedisServiceInterface
}

// NewAuthService creates a new authentication service
func NewAuthService(discordService DiscordServiceInterface, jwtService JWTServiceInterface, redisService RedisServiceInterface) *AuthService {
	return &AuthService{
		discordService: discordService,
		jwtService:     jwtService,
		redisService:   redisService,
	}
}

// GetAuthURL generates Discord OAuth2 authorization URL
func (a *AuthService) GetAuthURL(state string) string {
	return a.discordService.GetAuthURL(state)
}

// HandleCallback processes Discord OAuth2 callback
func (a *AuthService) HandleCallback(ctx context.Context, code string) (*models.AuthResponse, error) {
	// Exchange code for Discord token
	discordToken, err := a.discordService.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Discord
	discordUser, err := a.discordService.GetUserInfo(ctx, discordToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Create or update user
	user := &models.User{
		ID:            uuid.New().String(),
		DiscordUserID: discordUser.ID,
		Username:      discordUser.Username,
		Avatar:        discordUser.Avatar,
		Email:         discordUser.Email,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Create session
	sessionID := uuid.New().String()
	
	// Generate JWT tokens
	accessToken, accessExpiresAt, err := a.jwtService.GenerateAccessToken(user, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshExpiresAt, err := a.jwtService.GenerateRefreshToken(user, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store session in Redis
	session := &models.Session{
		ID:           sessionID,
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshExpiresAt, // Use refresh token expiry for session
		CreatedAt:    time.Now(),
	}

	err = a.redisService.StoreSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	// Return auth response
	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessExpiresAt.Sub(time.Now()).Seconds()),
		User:         *user,
	}, nil
}

// RefreshToken refreshes an access token using refresh token
func (a *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := a.jwtService.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get session from Redis
	session, err := a.redisService.GetSessionByRefreshToken(ctx, refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Clean up expired session
		a.redisService.DeleteSession(ctx, session.ID)
		return nil, fmt.Errorf("session expired")
	}

	// Create user from claims
	user := &models.User{
		ID:            claims.UserID,
		DiscordUserID: claims.DiscordUserID,
		Username:      claims.Username,
	}

	// Generate new access token
	newAccessToken, accessExpiresAt, err := a.jwtService.GenerateAccessToken(user, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Update session with new access token
	session.AccessToken = newAccessToken
	err = a.redisService.StoreSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: refreshTokenString, // Keep the same refresh token
		ExpiresIn:    int64(accessExpiresAt.Sub(time.Now()).Seconds()),
		User:         *user,
	}, nil
}

// ValidateAccessToken validates an access token
func (a *AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*models.User, error) {
	// Validate JWT token
	claims, err := a.jwtService.ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	// Check if session exists in Redis
	session, err := a.redisService.GetSession(ctx, claims.SessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Clean up expired session
		a.redisService.DeleteSession(ctx, session.ID)
		return nil, fmt.Errorf("session expired")
	}

	// Return user from claims
	return &models.User{
		ID:            claims.UserID,
		DiscordUserID: claims.DiscordUserID,
		Username:      claims.Username,
	}, nil
}

// Logout invalidates a session
func (a *AuthService) Logout(ctx context.Context, accessToken string) error {
	// Validate access token to get session ID
	claims, err := a.jwtService.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("invalid access token: %w", err)
	}

	// Delete session from Redis
	err = a.redisService.DeleteSession(ctx, claims.SessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}