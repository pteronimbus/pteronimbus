package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/testutils"
	"gorm.io/gorm"
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

	authService := NewAuthService(nil, mockDiscord, mockJWT, mockRedis)

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

			authService := NewAuthService(nil, mockDiscord, mockJWT, mockRedis)

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

			authService := NewAuthService(nil, mockDiscord, mockJWT, mockRedis)

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

			authService := NewAuthService(nil, mockDiscord, mockJWT, mockRedis)

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

			authService := NewAuthService(nil, mockDiscord, mockJWT, mockRedis)

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

func TestAuthService_SuperAdminRoleAssignment(t *testing.T) {
	// Setup test database
	db, cleanup := setupTestDatabaseWithModels(t)
	defer cleanup()

	// Create test configuration with super admin Discord ID
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "test-secret-key",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour * 24 * 7,
			Issuer:          "pteronimbus-test",
		},
		RBAC: config.RBACConfig{
			SuperAdminDiscordID: "197918357025062922", // Your Discord ID
			RoleSyncTTL:         time.Minute * 5,
			GuildCacheTTL:       time.Minute * 5,
			GracePeriod:         time.Minute * 2,
		},
	}

	// Create services
	rbacService := NewRBACService(db, &cfg.RBAC)
	mockDiscord := new(MockDiscordService)
	mockRedis := new(MockRedisService)

	// Create JWT service with RBAC integration
	jwtService := NewJWTServiceWithRBAC(cfg, rbacService)

	// Create auth service with RBAC integration
	authService := NewAuthServiceWithRBAC(db, mockDiscord, jwtService, mockRedis, rbacService)

	tests := []struct {
		name                    string
		discordUserID           string
		expectSuperAdminRole    bool
		setupMocks              func(*MockDiscordService, *MockRedisService)
		checkSuperAdminStatus   func(*testing.T, string)
	}{
		{
			name:                 "new user with super admin Discord ID should get super admin role",
			discordUserID:        "197918357025062922",
			expectSuperAdminRole: true,
			setupMocks: func(mockDiscord *MockDiscordService, mockRedis *MockRedisService) {
				// Mock Discord token exchange
				mockDiscord.On("ExchangeCodeForToken", mock.Anything, "test_code").Return(&models.DiscordTokenResponse{
					AccessToken:  "test_access_token",
					RefreshToken: "test_refresh_token",
					ExpiresIn:    3600,
				}, nil)

				// Mock Discord user info
				mockDiscord.On("GetUserInfo", mock.Anything, "test_access_token").Return(&models.DiscordUser{
					ID:       "197918357025062922",
					Username: "testuser",
					Avatar:   "test_avatar",
					Email:    "test@example.com",
				}, nil)

				// Mock Redis session storage
				mockRedis.On("StoreSession", mock.Anything, mock.Anything).Return(nil)
			},
			checkSuperAdminStatus: func(t *testing.T, userID string) {
				// Check that the user has super admin role (system-wide, not tenant-scoped)
				isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), userID)
				assert.NoError(t, err)
				assert.True(t, isSuperAdmin, "User should have super admin role")

				// Super admin role is system-wide, not tenant-scoped, so no tenant role should exist
				var userTenant models.UserTenant
				err = db.Where("user_id = ?", userID).First(&userTenant).Error
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err, "Super admin should not have tenant-scoped roles")
			},
		},
		{
			name:                 "new user without super admin Discord ID should not get super admin role",
			discordUserID:        "123456789012345678",
			expectSuperAdminRole: false,
			setupMocks: func(mockDiscord *MockDiscordService, mockRedis *MockRedisService) {
				// Mock Discord token exchange
				mockDiscord.On("ExchangeCodeForToken", mock.Anything, "test_code").Return(&models.DiscordTokenResponse{
					AccessToken:  "test_access_token",
					RefreshToken: "test_refresh_token",
					ExpiresIn:    3600,
				}, nil)

				// Mock Discord user info
				mockDiscord.On("GetUserInfo", mock.Anything, "test_access_token").Return(&models.DiscordUser{
					ID:       "123456789012345678",
					Username: "regularuser",
					Avatar:   "test_avatar",
					Email:    "regular@example.com",
				}, nil)

				// Mock Redis session storage
				mockRedis.On("StoreSession", mock.Anything, mock.Anything).Return(nil)
			},
			checkSuperAdminStatus: func(t *testing.T, userID string) {
				// Get user details for debugging
				var user models.User
				err := db.Where("id = ?", userID).First(&user).Error
				if err != nil {
					t.Logf("Failed to get user: %v", err)
				} else {
					t.Logf("User details: ID=%s, DiscordID=%s, Username=%s", user.ID, user.DiscordUserID, user.Username)
				}

				// Check that the user does not have super admin role
				isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), userID)
				if err != nil {
					t.Logf("IsSuperAdmin error: %v", err)
				}
				t.Logf("User %s isSuperAdmin: %v", userID, isSuperAdmin)
				assert.NoError(t, err)
				assert.False(t, isSuperAdmin, "User should not have super admin role")

				// Regular users should not have any tenant roles (since they're not in any tenant)
				var userTenant models.UserTenant
				err = db.Where("user_id = ?", userID).First(&userTenant).Error
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err, "Regular user should not have any tenant roles")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks before each test case
			mockDiscord.ExpectedCalls = nil
			mockRedis.ExpectedCalls = nil

			// Setup mocks
			tt.setupMocks(mockDiscord, mockRedis)

			// Call HandleCallback
			authResponse, err := authService.HandleCallback(context.Background(), "test_code")

			// Assertions
			assert.NoError(t, err)
			assert.NotNil(t, authResponse)
			assert.NotEmpty(t, authResponse.AccessToken)
			assert.NotEmpty(t, authResponse.RefreshToken)
			assert.NotNil(t, authResponse.User)

			// Check super admin status
			tt.checkSuperAdminStatus(t, authResponse.User.ID)

			// Verify mocks
			mockDiscord.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}

func TestAuthService_SuperAdminJWTInclusion(t *testing.T) {
	// Setup test database
	db, cleanup := setupTestDatabaseWithModels(t)
	defer cleanup()

	// Create test configuration with super admin Discord ID
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "test-secret-key",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour * 24 * 7,
			Issuer:          "pteronimbus-test",
		},
		RBAC: config.RBACConfig{
			SuperAdminDiscordID: "197918357025062922",
			RoleSyncTTL:         time.Minute * 5,
			GuildCacheTTL:       time.Minute * 5,
			GracePeriod:         time.Minute * 2,
		},
	}

	// Create services
	rbacService := NewRBACService(db, &cfg.RBAC)
	mockDiscord := new(MockDiscordService)
	mockRedis := new(MockRedisService)

	// Create JWT service with RBAC integration
	jwtService := NewJWTServiceWithRBAC(cfg, rbacService)

	// Create auth service with RBAC integration
	authService := NewAuthServiceWithRBAC(db, mockDiscord, jwtService, mockRedis, rbacService)

	tests := []struct {
		name                 string
		discordUserID        string
		expectSuperAdminJWT  bool
		setupMocks           func(*MockDiscordService, *MockRedisService)
	}{
		{
			name:                "super admin user should have IsSuperAdmin in JWT",
			discordUserID:       "197918357025062922",
			expectSuperAdminJWT: true,
			setupMocks: func(mockDiscord *MockDiscordService, mockRedis *MockRedisService) {
				// Mock Discord token exchange
				mockDiscord.On("ExchangeCodeForToken", mock.Anything, "test_code").Return(&models.DiscordTokenResponse{
					AccessToken:  "test_access_token",
					RefreshToken: "test_refresh_token",
					ExpiresIn:    3600,
				}, nil)

				// Mock Discord user info
				mockDiscord.On("GetUserInfo", mock.Anything, "test_access_token").Return(&models.DiscordUser{
					ID:       "197918357025062922",
					Username: "testuser",
					Avatar:   "test_avatar",
					Email:    "test@example.com",
				}, nil)

				// Mock Redis session storage
				mockRedis.On("StoreSession", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:                "regular user should not have IsSuperAdmin in JWT",
			discordUserID:       "123456789012345678",
			expectSuperAdminJWT: false,
			setupMocks: func(mockDiscord *MockDiscordService, mockRedis *MockRedisService) {
				// Mock Discord token exchange
				mockDiscord.On("ExchangeCodeForToken", mock.Anything, "test_code").Return(&models.DiscordTokenResponse{
					AccessToken:  "test_access_token",
					RefreshToken: "test_refresh_token",
					ExpiresIn:    3600,
				}, nil)

				// Mock Discord user info
				mockDiscord.On("GetUserInfo", mock.Anything, "test_access_token").Return(&models.DiscordUser{
					ID:       "123456789012345678",
					Username: "regularuser",
					Avatar:   "test_avatar",
					Email:    "regular@example.com",
				}, nil)

				// Mock Redis session storage
				mockRedis.On("StoreSession", mock.Anything, mock.Anything).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks before each test case
			mockDiscord.ExpectedCalls = nil
			mockRedis.ExpectedCalls = nil

			// Setup mocks
			tt.setupMocks(mockDiscord, mockRedis)

			// Call HandleCallback
			authResponse, err := authService.HandleCallback(context.Background(), "test_code")

			// Assertions
			assert.NoError(t, err)
			assert.NotNil(t, authResponse)
			assert.NotEmpty(t, authResponse.AccessToken)

			// Decode JWT to check super admin status
			claims, err := jwtService.ValidateToken(authResponse.AccessToken)
			if err != nil {
				t.Logf("JWT validation error: %v", err)
				t.Logf("Access token: %s", authResponse.AccessToken)
				t.Fatalf("JWT validation failed: %v", err)
			}
			assert.NoError(t, err)
			assert.NotNil(t, claims)

			// Check IsSuperAdmin field in JWT claims
			assert.Equal(t, tt.expectSuperAdminJWT, claims.IsSuperAdmin, 
				"JWT IsSuperAdmin should be %v for user %s", tt.expectSuperAdminJWT, tt.discordUserID)

			// Verify mocks
			mockDiscord.AssertExpectations(t)
			mockRedis.AssertExpectations(t)
		})
	}
}

func setupTestDatabaseWithModels(t *testing.T) (*gorm.DB, func()) {
	return testutils.SetupTestDatabaseWithModels(t,
		&models.User{},
		&models.Session{},
		&models.UserTenant{},
		&models.Tenant{},
		&models.Permission{},
		&models.Role{},
		&models.SystemRole{},
		&models.UserSystemRole{},
		&models.PermissionAuditLog{},
		&models.GuildMembershipCache{},
	)
}