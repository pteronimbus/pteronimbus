package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// MockAuthService for middleware testing
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*models.User, error) {
	args := m.Called(ctx, accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) ParseTokenClaims(accessToken string) (*models.JWTClaims, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JWTClaims), args.Error(1)
}

func (m *MockAuthService) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockAuthService) HandleCallback(ctx context.Context, code string) (*models.AuthResponse, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Logout(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}

func setupTestMiddleware() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware_RequireAuth(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(*MockAuthService)
		expectedStatus int
		expectAbort    bool
		checkContext   func(*testing.T, *gin.Context)
	}{
		{
			name:       "valid bearer token",
			authHeader: "Bearer valid_token",
			setupMock: func(m *MockAuthService) {
				user := &models.User{
					ID:            "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
				}
				claims := &models.JWTClaims{
					UserID:        "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
					SessionID:     "session_id",
					SystemRoles:   []string{},
				}
				m.On("ValidateAccessToken", mock.Anything, "valid_token").Return(user, nil)
				m.On("ParseTokenClaims", "valid_token").Return(claims, nil)
			},
			expectedStatus: http.StatusOK,
			expectAbort:    false,
			checkContext: func(t *testing.T, c *gin.Context) {
				user, exists := GetUserFromContext(c)
				assert.True(t, exists)
				assert.Equal(t, "user_id", user.ID)
				assert.Equal(t, "testuser", user.Username)

				userID, exists := c.Get("user_id")
				assert.True(t, exists)
				assert.Equal(t, "user_id", userID)

				discordUserID, exists := c.Get("discord_user_id")
				assert.True(t, exists)
				assert.Equal(t, "discord_user_id", discordUserID)
			},
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "invalid authorization header format - no bearer",
			authHeader:     "InvalidFormat token",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:           "invalid authorization header format - no token",
			authHeader:     "Bearer",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:       "invalid authorization header format - only bearer",
			authHeader: "Bearer ",
			setupMock: func(m *MockAuthService) {
				// The middleware will try to validate an empty token
				m.On("ValidateAccessToken", mock.Anything, "").Return(nil, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalid_token",
			setupMock: func(m *MockAuthService) {
				m.On("ValidateAccessToken", mock.Anything, "invalid_token").Return(nil, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
		{
			name:       "expired token",
			authHeader: "Bearer expired_token",
			setupMock: func(m *MockAuthService) {
				m.On("ValidateAccessToken", mock.Anything, "expired_token").Return(nil, errors.New("token expired"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			tt.setupMock(mockAuthService)

			middleware := NewAuthMiddleware(mockAuthService)
			router := setupTestMiddleware()

			var contextToCheck *gin.Context
			var wasAborted bool

			router.Use(middleware.RequireAuth())
			router.GET("/protected", func(c *gin.Context) {
				contextToCheck = c
				wasAborted = c.IsAborted()
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectAbort {
				assert.True(t, wasAborted || w.Code != http.StatusOK)
			} else {
				assert.False(t, wasAborted)
				assert.NotNil(t, contextToCheck)
				tt.checkContext(t, contextToCheck)
			}

			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthMiddleware_OptionalAuth(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkContext   func(*testing.T, *gin.Context)
	}{
		{
			name:       "valid bearer token",
			authHeader: "Bearer valid_token",
			setupMock: func(m *MockAuthService) {
				user := &models.User{
					ID:            "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
				}
				m.On("ValidateAccessToken", mock.Anything, "valid_token").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				user, exists := GetUserFromContext(c)
				assert.True(t, exists)
				assert.Equal(t, "user_id", user.ID)
			},
		},
		{
			name:           "no authorization header",
			authHeader:     "",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				user, exists := GetUserFromContext(c)
				assert.False(t, exists)
				assert.Nil(t, user)
			},
		},
		{
			name:           "invalid authorization header format",
			authHeader:     "InvalidFormat token",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				user, exists := GetUserFromContext(c)
				assert.False(t, exists)
				assert.Nil(t, user)
			},
		},
		{
			name:       "invalid token - continues without user",
			authHeader: "Bearer invalid_token",
			setupMock: func(m *MockAuthService) {
				m.On("ValidateAccessToken", mock.Anything, "invalid_token").Return(nil, errors.New("invalid token"))
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				user, exists := GetUserFromContext(c)
				assert.False(t, exists)
				assert.Nil(t, user)
			},
		},
		{
			name:       "expired token - continues without user",
			authHeader: "Bearer expired_token",
			setupMock: func(m *MockAuthService) {
				m.On("ValidateAccessToken", mock.Anything, "expired_token").Return(nil, errors.New("token expired"))
			},
			expectedStatus: http.StatusOK,
			checkContext: func(t *testing.T, c *gin.Context) {
				user, exists := GetUserFromContext(c)
				assert.False(t, exists)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			tt.setupMock(mockAuthService)

			middleware := NewAuthMiddleware(mockAuthService)
			router := setupTestMiddleware()

			var contextToCheck *gin.Context

			router.Use(middleware.OptionalAuth())
			router.GET("/optional", func(c *gin.Context) {
				contextToCheck = c
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/optional", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.NotNil(t, contextToCheck)
			tt.checkContext(t, contextToCheck)

			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	tests := []struct {
		name        string
		setupContext func(*gin.Context)
		expectUser   bool
		checkUser    func(*testing.T, *models.User)
	}{
		{
			name: "user exists in context",
			setupContext: func(c *gin.Context) {
				user := &models.User{
					ID:            "user_id",
					DiscordUserID: "discord_user_id",
					Username:      "testuser",
				}
				c.Set("user", user)
			},
			expectUser: true,
			checkUser: func(t *testing.T, user *models.User) {
				assert.Equal(t, "user_id", user.ID)
				assert.Equal(t, "discord_user_id", user.DiscordUserID)
				assert.Equal(t, "testuser", user.Username)
			},
		},
		{
			name:         "user does not exist in context",
			setupContext: func(c *gin.Context) {},
			expectUser:   false,
		},
		{
			name: "wrong type in context",
			setupContext: func(c *gin.Context) {
				c.Set("user", "not_a_user_struct")
			},
			expectUser: false,
		},
		{
			name: "nil user in context",
			setupContext: func(c *gin.Context) {
				c.Set("user", nil)
			},
			expectUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestMiddleware()
			
			var user *models.User
			var exists bool

			router.GET("/test", func(c *gin.Context) {
				tt.setupContext(c)
				user, exists = GetUserFromContext(c)
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectUser, exists)
			if tt.expectUser {
				assert.NotNil(t, user)
				tt.checkUser(t, user)
			} else {
				assert.Nil(t, user)
			}
		})
	}
}

func TestAuthMiddleware_Integration(t *testing.T) {
	// Test that middleware properly integrates with multiple routes
	mockAuthService := new(MockAuthService)
	
	validUser := &models.User{
		ID:            "user_id",
		DiscordUserID: "discord_user_id",
		Username:      "testuser",
	}
	
	validClaims := &models.JWTClaims{
		UserID:        "user_id",
		DiscordUserID: "discord_user_id", 
		Username:      "testuser",
		SessionID:     "session_id",
	}
	
	mockAuthService.On("ValidateAccessToken", mock.Anything, "valid_token").Return(validUser, nil)
	mockAuthService.On("ParseTokenClaims", "valid_token").Return(validClaims, nil)
	mockAuthService.On("ValidateAccessToken", mock.Anything, "invalid_token").Return(nil, errors.New("invalid token"))

	middleware := NewAuthMiddleware(mockAuthService)
	router := setupTestMiddleware()

	// Protected route
	protected := router.Group("/api/protected")
	protected.Use(middleware.RequireAuth())
	protected.GET("/resource", func(c *gin.Context) {
		user, exists := GetUserFromContext(c)
		assert.True(t, exists)
		c.JSON(http.StatusOK, gin.H{"user_id": user.ID})
	})

	// Optional auth route
	optional := router.Group("/api/optional")
	optional.Use(middleware.OptionalAuth())
	optional.GET("/resource", func(c *gin.Context) {
		user, exists := GetUserFromContext(c)
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": user.ID, "authenticated": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"authenticated": false})
		}
	})

	// Public route (no middleware)
	router.GET("/api/public/resource", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "public"})
	})

	tests := []struct {
		name           string
		path           string
		authHeader     string
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:           "protected route with valid token",
			path:           "/api/protected/resource",
			authHeader:     "Bearer valid_token",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "user_id")
			},
		},
		{
			name:           "protected route without token",
			path:           "/api/protected/resource",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "protected route with invalid token",
			path:           "/api/protected/resource",
			authHeader:     "Bearer invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "optional route with valid token",
			path:           "/api/optional/resource",
			authHeader:     "Bearer valid_token",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "authenticated\":true")
			},
		},
		{
			name:           "optional route without token",
			path:           "/api/optional/resource",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "authenticated\":false")
			},
		},
		{
			name:           "public route",
			path:           "/api/public/resource",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Contains(t, w.Body.String(), "public")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}

	mockAuthService.AssertExpectations(t)
}