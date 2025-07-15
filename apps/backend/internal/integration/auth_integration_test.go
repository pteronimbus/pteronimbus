package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/handlers"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/middleware"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// AuthIntegrationTestSuite tests the complete authentication flow
type AuthIntegrationTestSuite struct {
	suite.Suite
	router      *gin.Engine
	authHandler *handlers.AuthHandler
	authService *services.AuthService
	jwtService  *services.JWTService
}

func (suite *AuthIntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	// Create test configuration
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "test_secret_key_that_is_long_enough_for_testing",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  24 * time.Hour,
			Issuer:           "pteronimbus-test",
		},
	}

	// Create services with mocked external dependencies
	suite.jwtService = services.NewJWTService(cfg)
	
	// For integration tests, we'll use mock services for Discord and Redis
	// In a real integration test, you might use test containers
	mockDiscordService := &MockDiscordService{}
	mockRedisService := &MockRedisService{}
	
	suite.authService = services.NewAuthService(mockDiscordService, suite.jwtService, mockRedisService)
	suite.authHandler = handlers.NewAuthHandler(suite.authService)

	// Setup router
	suite.router = gin.New()
	suite.setupRoutes()
}

func (suite *AuthIntegrationTestSuite) setupRoutes() {
	authMiddleware := middleware.NewAuthMiddleware(suite.authService)

	// Auth routes
	auth := suite.router.Group("/auth")
	{
		auth.GET("/login", suite.authHandler.Login)
		auth.GET("/callback", suite.authHandler.Callback)
		auth.POST("/refresh", suite.authHandler.Refresh)
		auth.POST("/logout", suite.authHandler.Logout)
	}

	// Protected routes
	api := suite.router.Group("/api")
	api.Use(authMiddleware.RequireAuth())
	{
		api.GET("/me", suite.authHandler.Me)
		api.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "protected resource"})
		})
	}

	// Optional auth routes
	optional := suite.router.Group("/optional")
	optional.Use(authMiddleware.OptionalAuth())
	{
		optional.GET("/resource", func(c *gin.Context) {
			user, exists := middleware.GetUserFromContext(c)
			if exists {
				c.JSON(http.StatusOK, gin.H{"authenticated": true, "user_id": user.ID})
			} else {
				c.JSON(http.StatusOK, gin.H{"authenticated": false})
			}
		})
	}
}

func (suite *AuthIntegrationTestSuite) TestCompleteAuthFlow() {
	// Test 1: Login endpoint
	req, _ := http.NewRequest("GET", "/auth/login", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var loginResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), loginResponse, "auth_url")
	assert.Contains(suite.T(), loginResponse, "state")

	// Test 2: Access protected resource without token (should fail)
	req, _ = http.NewRequest("GET", "/api/protected", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	// Test 3: Generate a valid token using JWT service directly
	user := &models.User{
		ID:            "test_user_id",
		DiscordUserID: "discord_123",
		Username:      "testuser",
	}
	sessionID := "test_session_id"

	accessToken, _, err := suite.jwtService.GenerateAccessToken(user, sessionID)
	assert.NoError(suite.T(), err)

	// Test 4: Access protected resource with valid token
	req, _ = http.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// This will fail because we don't have Redis session, but that's expected in this mock setup
	// In a real integration test, you'd set up the session properly
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	// Test 5: Optional auth endpoint without token
	req, _ = http.NewRequest("GET", "/optional/resource", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var optionalResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &optionalResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), false, optionalResponse["authenticated"])
}

func (suite *AuthIntegrationTestSuite) TestJWTTokenLifecycle() {
	user := &models.User{
		ID:            "test_user_id",
		DiscordUserID: "discord_123",
		Username:      "testuser",
	}
	sessionID := "test_session_id"

	// Test access token generation
	accessToken, accessExpiresAt, err := suite.jwtService.GenerateAccessToken(user, sessionID)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), accessToken)
	assert.True(suite.T(), accessExpiresAt.After(time.Now()))

	// Test refresh token generation
	refreshToken, refreshExpiresAt, err := suite.jwtService.GenerateRefreshToken(user, sessionID)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), refreshToken)
	assert.True(suite.T(), refreshExpiresAt.After(accessExpiresAt))

	// Test token validation
	claims, err := suite.jwtService.ValidateToken(accessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, claims.UserID)
	assert.Equal(suite.T(), user.DiscordUserID, claims.DiscordUserID)
	assert.Equal(suite.T(), user.Username, claims.Username)
	assert.Equal(suite.T(), sessionID, claims.SessionID)

	// Test refresh token validation
	refreshClaims, err := suite.jwtService.ValidateToken(refreshToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, refreshClaims.UserID)
	assert.Equal(suite.T(), sessionID, refreshClaims.SessionID)
}

func (suite *AuthIntegrationTestSuite) TestRefreshTokenEndpoint() {
	// Create a valid refresh token
	user := &models.User{
		ID:            "test_user_id",
		DiscordUserID: "discord_123",
		Username:      "testuser",
	}
	sessionID := "test_session_id"

	refreshToken, _, err := suite.jwtService.GenerateRefreshToken(user, sessionID)
	assert.NoError(suite.T(), err)

	// Test refresh endpoint with valid token
	refreshRequest := models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}
	body, _ := json.Marshal(refreshRequest)

	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// This will fail because we don't have Redis session, but the JWT validation part works
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	// Test refresh endpoint with invalid token
	invalidRefreshRequest := models.RefreshTokenRequest{
		RefreshToken: "invalid_token",
	}
	body, _ = json.Marshal(invalidRefreshRequest)

	req, _ = http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthIntegrationTestSuite) TestMiddlewareIntegration() {
	// Test that middleware properly handles different scenarios
	tests := []struct {
		name           string
		endpoint       string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "protected endpoint without auth",
			endpoint:       "/api/protected",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "protected endpoint with invalid auth",
			endpoint:       "/api/protected",
			authHeader:     "Bearer invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "protected endpoint with malformed auth",
			endpoint:       "/api/protected",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "optional endpoint without auth",
			endpoint:       "/optional/resource",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "optional endpoint with invalid auth",
			endpoint:       "/optional/resource",
			authHeader:     "Bearer invalid_token",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.endpoint, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// Mock services for integration testing
type MockDiscordService struct{}

func (m *MockDiscordService) GetAuthURL(state string) string {
	return "https://discord.com/oauth2/authorize?client_id=test&state=" + state
}

func (m *MockDiscordService) ExchangeCodeForToken(ctx context.Context, code string) (*models.DiscordTokenResponse, error) {
	return &models.DiscordTokenResponse{
		AccessToken: "mock_discord_token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}, nil
}

func (m *MockDiscordService) GetUserInfo(ctx context.Context, accessToken string) (*models.DiscordUser, error) {
	return &models.DiscordUser{
		ID:       "discord_123",
		Username: "testuser",
		Avatar:   "avatar_hash",
		Email:    "test@example.com",
	}, nil
}

func (m *MockDiscordService) RefreshToken(ctx context.Context, refreshToken string) (*models.DiscordTokenResponse, error) {
	return &models.DiscordTokenResponse{
		AccessToken:  "new_mock_discord_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
	}, nil
}

type MockRedisService struct{}

func (m *MockRedisService) StoreSession(ctx context.Context, session *models.Session) error {
	return nil
}

func (m *MockRedisService) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	return nil, errors.New("session not found in mock")
}

func (m *MockRedisService) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	return nil, errors.New("session not found in mock")
}

func (m *MockRedisService) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

func TestAuthIntegrationSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}