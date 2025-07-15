package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// Mock services for testing
type MockDiscordService struct {
	mock.Mock
}

func (m *MockDiscordService) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockDiscordService) ExchangeCodeForToken(ctx context.Context, code string) (*models.DiscordTokenResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DiscordTokenResponse), args.Error(1)
}

func (m *MockDiscordService) GetUserInfo(ctx context.Context, accessToken string) (*models.DiscordUser, error) {
	args := m.Called(ctx, accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DiscordUser), args.Error(1)
}

func (m *MockDiscordService) RefreshToken(ctx context.Context, refreshToken string) (*models.DiscordTokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DiscordTokenResponse), args.Error(1)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(user *models.User, sessionID string) (string, time.Time, error) {
	args := m.Called(user, sessionID)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockJWTService) GenerateRefreshToken(user *models.User, sessionID string) (string, time.Time, error) {
	args := m.Called(user, sessionID)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JWTClaims), args.Error(1)
}

type MockRedisService struct {
	mock.Mock
}

func (m *MockRedisService) StoreSession(ctx context.Context, session *models.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRedisService) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockRedisService) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockRedisService) DeleteSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func TestAuthService_GetAuthURL(t *testing.T) {
	mockDiscord := new(MockDiscordService)
	mockJWT := new(MockJWTService)
	mockRedis := new(MockRedisService)

	authService := NewAuthService(mockDiscord, mockJWT, mockRedis)

	state := "test_state"
	expectedURL := "https://discord.com/oauth2/authorize?client_id=test&state=test_state"

	mockDiscord.On("GetAuthURL", state).Return(expectedURL)

	result := authService.GetAuthURL(state)

	assert.Equal(t, expectedURL, result)
	mockDiscord.AssertExpectations(t)
}

func TestAuthService_HandleCallback(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		setupMocks    func(*MockDiscordService, *MockJWTService, *MockRedisService)
		expectedError bool
		checkResult   func(*testing.T, *models.AuthResponse)
	}{
		{
			name: "successful callback",
			code: "valid_code",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				discordToken := &models.DiscordTokenResponse{
					AccessToken: "discord_access_token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}
				discordUser := &models.DiscordUser{
					ID:       "discord_user_id",
					Username: "testuser",
					Avatar:   "avatar_hash",
					Email:    "test@example.com",
				}

				discord.On("ExchangeCodeForToken", mock.Anything, "valid_code").Return(discordToken, nil)
				discord.On("GetUserInfo", mock.Anything, "discord_access_token").Return(discordUser, nil)

				expiresAt := time.Now().Add(time.Hour)
				jwt.On("GenerateAccessToken", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return("access_token", expiresAt, nil)
				jwt.On("GenerateRefreshToken", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return("refresh_token", expiresAt, nil)

				redis.On("StoreSession", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil)
			},
			expectedError: false,
			checkResult: func(t *testing.T, result *models.AuthResponse) {
				assert.Equal(t, "access_token", result.AccessToken)
				assert.Equal(t, "refresh_token", result.RefreshToken)
				assert.Equal(t, "testuser", result.User.Username)
				assert.Equal(t, "discord_user_id", result.User.DiscordUserID)
			},
		},
		{
			name: "discord token exchange error",
			code: "invalid_code",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				discord.On("ExchangeCodeForToken", mock.Anything, "invalid_code").Return(nil, errors.New("invalid code"))
			},
			expectedError: true,
		},
		{
			name: "discord user info error",
			code: "valid_code",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				discordToken := &models.DiscordTokenResponse{
					AccessToken: "discord_access_token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}

				discord.On("ExchangeCodeForToken", mock.Anything, "valid_code").Return(discordToken, nil)
				discord.On("GetUserInfo", mock.Anything, "discord_access_token").Return(nil, errors.New("user info error"))
			},
			expectedError: true,
		},
		{
			name: "jwt generation error",
			code: "valid_code",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				discordToken := &models.DiscordTokenResponse{
					AccessToken: "discord_access_token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}
				discordUser := &models.DiscordUser{
					ID:       "discord_user_id",
					Username: "testuser",
					Avatar:   "avatar_hash",
					Email:    "test@example.com",
				}

				discord.On("ExchangeCodeForToken", mock.Anything, "valid_code").Return(discordToken, nil)
				discord.On("GetUserInfo", mock.Anything, "discord_access_token").Return(discordUser, nil)

				jwt.On("GenerateAccessToken", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return("", time.Time{}, errors.New("jwt error"))
			},
			expectedError: true,
		},
		{
			name: "redis storage error",
			code: "valid_code",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				discordToken := &models.DiscordTokenResponse{
					AccessToken: "discord_access_token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}
				discordUser := &models.DiscordUser{
					ID:       "discord_user_id",
					Username: "testuser",
					Avatar:   "avatar_hash",
					Email:    "test@example.com",
				}

				discord.On("ExchangeCodeForToken", mock.Anything, "valid_code").Return(discordToken, nil)
				discord.On("GetUserInfo", mock.Anything, "discord_access_token").Return(discordUser, nil)

				expiresAt := time.Now().Add(time.Hour)
				jwt.On("GenerateAccessToken", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return("access_token", expiresAt, nil)
				jwt.On("GenerateRefreshToken", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return("refresh_token", expiresAt, nil)

				redis.On("StoreSession", mock.Anything, mock.AnythingOfType("*models.Session")).Return(errors.New("redis error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDiscord := new(MockDiscordService)
			mockJWT := new(MockJWTService)
			mockRedis := new(MockRedisService)

			tt.setupMocks(mockDiscord, mockJWT, mockRedis)

			authService := NewAuthService(mockDiscord, mockJWT, mockRedis)

			result, err := authService.HandleCallback(context.Background(), tt.code)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				tt.checkResult(t, result)
			}

			mockDiscord.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		refreshToken  string
		setupMocks    func(*MockDiscordService, *MockJWTService, *MockRedisService)
		expectedError bool
		checkResult   func(*testing.T, *models.AuthResponse)
	}{
		{
			name:         "successful token refresh",
			refreshToken: "valid_refresh_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}
				session := &models.Session{
					ID:           "session_id",
					UserID:       "user_id",
					AccessToken:  "old_access_token",
					RefreshToken: "valid_refresh_token",
					ExpiresAt:    time.Now().Add(time.Hour),
					CreatedAt:    time.Now(),
				}

				jwt.On("ValidateToken", "valid_refresh_token").Return(claims, nil)
				redis.On("GetSessionByRefreshToken", mock.Anything, "valid_refresh_token").Return(session, nil)

				expiresAt := time.Now().Add(time.Hour)
				jwt.On("GenerateAccessToken", mock.AnythingOfType("*models.User"), "session_id").Return("new_access_token", expiresAt, nil)
				redis.On("StoreSession", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil)
			},
			expectedError: false,
			checkResult: func(t *testing.T, result *models.AuthResponse) {
				assert.Equal(t, "new_access_token", result.AccessToken)
				assert.Equal(t, "valid_refresh_token", result.RefreshToken)
			},
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid_refresh_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				jwt.On("ValidateToken", "invalid_refresh_token").Return(nil, errors.New("invalid token"))
			},
			expectedError: true,
		},
		{
			name:         "session not found",
			refreshToken: "valid_refresh_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}

				jwt.On("ValidateToken", "valid_refresh_token").Return(claims, nil)
				redis.On("GetSessionByRefreshToken", mock.Anything, "valid_refresh_token").Return(nil, errors.New("session not found"))
			},
			expectedError: true,
		},
		{
			name:         "expired session",
			refreshToken: "valid_refresh_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}
				session := &models.Session{
					ID:           "session_id",
					UserID:       "user_id",
					AccessToken:  "old_access_token",
					RefreshToken: "valid_refresh_token",
					ExpiresAt:    time.Now().Add(-time.Hour), // Expired
					CreatedAt:    time.Now(),
				}

				jwt.On("ValidateToken", "valid_refresh_token").Return(claims, nil)
				redis.On("GetSessionByRefreshToken", mock.Anything, "valid_refresh_token").Return(session, nil)
				redis.On("DeleteSession", mock.Anything, "session_id").Return(nil)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDiscord := new(MockDiscordService)
			mockJWT := new(MockJWTService)
			mockRedis := new(MockRedisService)

			tt.setupMocks(mockDiscord, mockJWT, mockRedis)

			authService := NewAuthService(mockDiscord, mockJWT, mockRedis)

			result, err := authService.RefreshToken(context.Background(), tt.refreshToken)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				tt.checkResult(t, result)
			}

			mockDiscord.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestAuthService_ValidateAccessToken(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		setupMocks    func(*MockDiscordService, *MockJWTService, *MockRedisService)
		expectedError bool
		checkResult   func(*testing.T, *models.User)
	}{
		{
			name:        "valid access token",
			accessToken: "valid_access_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}
				session := &models.Session{
					ID:           "session_id",
					UserID:       "user_id",
					AccessToken:  "valid_access_token",
					RefreshToken: "refresh_token",
					ExpiresAt:    time.Now().Add(time.Hour),
					CreatedAt:    time.Now(),
				}

				jwt.On("ValidateToken", "valid_access_token").Return(claims, nil)
				redis.On("GetSession", mock.Anything, "session_id").Return(session, nil)
			},
			expectedError: false,
			checkResult: func(t *testing.T, user *models.User) {
				assert.Equal(t, "user_id", user.ID)
				assert.Equal(t, "discord_user_id", user.DiscordUserID)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name:        "invalid access token",
			accessToken: "invalid_access_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				jwt.On("ValidateToken", "invalid_access_token").Return(nil, errors.New("invalid token"))
			},
			expectedError: true,
		},
		{
			name:        "session not found",
			accessToken: "valid_access_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}

				jwt.On("ValidateToken", "valid_access_token").Return(claims, nil)
				redis.On("GetSession", mock.Anything, "session_id").Return(nil, errors.New("session not found"))
			},
			expectedError: true,
		},
		{
			name:        "expired session",
			accessToken: "valid_access_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}
				session := &models.Session{
					ID:           "session_id",
					UserID:       "user_id",
					AccessToken:  "valid_access_token",
					RefreshToken: "refresh_token",
					ExpiresAt:    time.Now().Add(-time.Hour), // Expired
					CreatedAt:    time.Now(),
				}

				jwt.On("ValidateToken", "valid_access_token").Return(claims, nil)
				redis.On("GetSession", mock.Anything, "session_id").Return(session, nil)
				redis.On("DeleteSession", mock.Anything, "session_id").Return(nil)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDiscord := new(MockDiscordService)
			mockJWT := new(MockJWTService)
			mockRedis := new(MockRedisService)

			tt.setupMocks(mockDiscord, mockJWT, mockRedis)

			authService := NewAuthService(mockDiscord, mockJWT, mockRedis)

			result, err := authService.ValidateAccessToken(context.Background(), tt.accessToken)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				tt.checkResult(t, result)
			}

			mockDiscord.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		setupMocks    func(*MockDiscordService, *MockJWTService, *MockRedisService)
		expectedError bool
	}{
		{
			name:        "successful logout",
			accessToken: "valid_access_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}

				jwt.On("ValidateToken", "valid_access_token").Return(claims, nil)
				redis.On("DeleteSession", mock.Anything, "session_id").Return(nil)
			},
			expectedError: false,
		},
		{
			name:        "invalid access token",
			accessToken: "invalid_access_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				jwt.On("ValidateToken", "invalid_access_token").Return(nil, errors.New("invalid token"))
			},
			expectedError: true,
		},
		{
			name:        "redis delete error",
			accessToken: "valid_access_token",
			setupMocks: func(discord *MockDiscordService, jwt *MockJWTService, redis *MockRedisService) {
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
				}

				jwt.On("ValidateToken", "valid_access_token").Return(claims, nil)
				redis.On("DeleteSession", mock.Anything, "session_id").Return(errors.New("redis error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDiscord := new(MockDiscordService)
			mockJWT := new(MockJWTService)
			mockRedis := new(MockRedisService)

			tt.setupMocks(mockDiscord, mockJWT, mockRedis)

			authService := NewAuthService(mockDiscord, mockJWT, mockRedis)

			err := authService.Logout(context.Background(), tt.accessToken)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDiscord.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}