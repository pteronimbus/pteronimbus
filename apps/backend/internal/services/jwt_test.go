package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

func TestJWTService_GenerateAccessToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test_secret_key_that_is_long_enough",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "pteronimbus-test",
		},
	}

	jwtService := NewJWTService(cfg)

	user := &models.User{
		ID:            "user_id",
		DiscordUserID: "discord_user_id",
		Username:      "testuser",
	}
	sessionID := "session_id"

	token, expiresAt, err := jwtService.GenerateAccessToken(user, sessionID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))
	assert.True(t, expiresAt.Before(time.Now().Add(time.Hour+time.Minute))) // Should expire within an hour

	// Validate the token can be parsed
	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.DiscordUserID, claims.DiscordUserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, sessionID, claims.SessionID)
	assert.Equal(t, cfg.JWT.Issuer, claims.Issuer)
	assert.Empty(t, claims.SystemRoles) // Default should be empty
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test_secret_key_that_is_long_enough",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "pteronimbus-test",
		},
	}

	jwtService := NewJWTService(cfg)

	user := &models.User{
		ID:            "user_id",
		DiscordUserID: "discord_user_id",
		Username:      "testuser",
	}
	sessionID := "session_id"

	token, expiresAt, err := jwtService.GenerateRefreshToken(user, sessionID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))
	assert.True(t, expiresAt.Before(time.Now().Add(24*time.Hour+time.Minute))) // Should expire within 24 hours

	// Validate the token can be parsed
	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.DiscordUserID, claims.DiscordUserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, sessionID, claims.SessionID)
	assert.Equal(t, cfg.JWT.Issuer, claims.Issuer)
	assert.Empty(t, claims.SystemRoles) // Default should be empty
}

func TestJWTService_ValidateToken(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test_secret_key_that_is_long_enough",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "pteronimbus-test",
		},
	}

	jwtService := NewJWTService(cfg)

	tests := []struct {
		name          string
		setupToken    func() string
		expectedError bool
		checkClaims   func(*testing.T, *models.JWTClaims)
	}{
		{
			name: "valid token",
			setupToken: func() string {
				user := &models.User{
					ID:            "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
				}
				token, _, _ := jwtService.GenerateAccessToken(user, "session_id")
				return token
			},
			expectedError: false,
			checkClaims: func(t *testing.T, claims *models.JWTClaims) {
				assert.Equal(t, "user_id", claims.UserID)
				assert.Equal(t, "discord_user_id", claims.DiscordUserID)
				assert.Equal(t, "testuser", claims.Username)
				assert.Equal(t, "session_id", claims.SessionID)
				assert.Empty(t, claims.SystemRoles) // Default should be empty
			},
		},
		{
			name: "invalid token format",
			setupToken: func() string {
				return "invalid.token.format"
			},
			expectedError: true,
		},
		{
			name: "token with wrong signature",
			setupToken: func() string {
				// Create token with different secret
				wrongCfg := &config.Config{
					JWT: config.JWTConfig{
						Secret:           "wrong_secret_key_that_is_long_enough",
						AccessTokenTTL:   time.Hour,
						RefreshTokenTTL:  24 * time.Hour,
						Issuer:           "pteronimbus-test",
					},
				}
				wrongJWTService := NewJWTService(wrongCfg)
				user := &models.User{
					ID:            "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
				}
				token, _, _ := wrongJWTService.GenerateAccessToken(user, "session_id")
				return token
			},
			expectedError: true,
		},
		{
			name: "expired token",
			setupToken: func() string {
				// Create token with very short TTL
				shortCfg := &config.Config{
					JWT: config.JWTConfig{
						Secret:           "test_secret_key_that_is_long_enough",
						AccessTokenTTL:   -time.Hour, // Negative TTL to create expired token
						RefreshTokenTTL:  24 * time.Hour,
						Issuer:           "pteronimbus-test",
					},
				}
				shortJWTService := NewJWTService(shortCfg)
				user := &models.User{
					ID:            "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
				}
				token, _, _ := shortJWTService.GenerateAccessToken(user, "session_id")
				return token
			},
			expectedError: true,
		},
		{
			name: "malformed token",
			setupToken: func() string {
				return "not.a.jwt"
			},
			expectedError: true,
		},
		{
			name: "empty token",
			setupToken: func() string {
				return ""
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()

			claims, err := jwtService.ValidateToken(token)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				tt.checkClaims(t, claims)
			}
		})
	}
}

func TestJWTService_GetAccessTokenTTL(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test_secret_key_that_is_long_enough",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "pteronimbus-test",
		},
	}

	jwtService := NewJWTService(cfg)

	ttl := jwtService.GetAccessTokenTTL()

	assert.Equal(t, int64(3600), ttl) // 1 hour in seconds
}

func TestJWTService_TokenSigningMethod(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test_secret_key_that_is_long_enough",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "pteronimbus-test",
		},
	}

	jwtService := NewJWTService(cfg)

	user := &models.User{
		ID:            "user_id",
		DiscordUserID: "discord_user_id",
		Username:      "testuser",
	}

	tokenString, _, err := jwtService.GenerateAccessToken(user, "session_id")
	assert.NoError(t, err)

	// Parse token to check signing method
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, "HS256", token.Header["alg"])
}

func TestJWTService_ClaimsValidation(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test_secret_key_that_is_long_enough",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "pteronimbus-test",
		},
	}

	jwtService := NewJWTService(cfg)

	user := &models.User{
		ID:            "user_id",
		DiscordUserID: "discord_user_id",
		Username:      "testuser",
	}
	sessionID := "session_id"

	tokenString, expiresAt, err := jwtService.GenerateAccessToken(user, sessionID)
	assert.NoError(t, err)

	claims, err := jwtService.ValidateToken(tokenString)
	assert.NoError(t, err)

	// Check all claims are properly set
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.DiscordUserID, claims.DiscordUserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, sessionID, claims.SessionID)
	assert.Equal(t, cfg.JWT.Issuer, claims.Issuer)
	assert.Equal(t, user.ID, claims.Subject)
	assert.Empty(t, claims.SystemRoles) // Default should be empty

	// Check time claims
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)

	// Verify expiration time matches
	assert.WithinDuration(t, expiresAt, claims.ExpiresAt.Time, time.Second)

	// Verify issued at and not before are recent
	now := time.Now()
	assert.WithinDuration(t, now, claims.IssuedAt.Time, time.Minute)
	assert.WithinDuration(t, now, claims.NotBefore.Time, time.Minute)
}