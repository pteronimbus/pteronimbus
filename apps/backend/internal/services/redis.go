package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// RedisService handles Redis operations
type RedisService struct {
	client *redis.Client
}

// NewRedisService creates a new Redis service
func NewRedisService(cfg *config.Config) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &RedisService{
		client: rdb,
	}
}

// StoreSession stores a session in Redis
func (r *RedisService) StoreSession(ctx context.Context, session *models.Session) error {
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Store session with expiration
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	key := fmt.Sprintf("session:%s", session.ID)
	err = r.client.Set(ctx, key, sessionData, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store session: %w", err)
	}

	// Also store refresh token mapping
	refreshKey := fmt.Sprintf("refresh_token:%s", session.RefreshToken)
	err = r.client.Set(ctx, refreshKey, session.ID, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store refresh token mapping: %w", err)
	}

	return nil
}

// GetSession retrieves a session from Redis
func (r *RedisService) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	sessionData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session models.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// GetSessionByRefreshToken retrieves a session by refresh token
func (r *RedisService) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	sessionID, err := r.client.Get(ctx, refreshKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get session ID: %w", err)
	}

	return r.GetSession(ctx, sessionID)
}

// DeleteSession deletes a session from Redis
func (r *RedisService) DeleteSession(ctx context.Context, sessionID string) error {
	// Get session first to get refresh token
	session, err := r.GetSession(ctx, sessionID)
	if err != nil {
		// If session doesn't exist, consider it already deleted
		return nil
	}

	// Delete session
	sessionKey := fmt.Sprintf("session:%s", sessionID)
	err = r.client.Del(ctx, sessionKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Delete refresh token mapping
	refreshKey := fmt.Sprintf("refresh_token:%s", session.RefreshToken)
	err = r.client.Del(ctx, refreshKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete refresh token mapping: %w", err)
	}

	return nil
}

// UpdateSessionExpiry updates session expiry time
func (r *RedisService) UpdateSessionExpiry(ctx context.Context, sessionID string, expiresAt time.Time) error {
	session, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = expiresAt

	// Re-store with new expiry
	return r.StoreSession(ctx, session)
}

// Ping checks Redis connectivity
func (r *RedisService) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisService) Close() error {
	return r.client.Close()
}