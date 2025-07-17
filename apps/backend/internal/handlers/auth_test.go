package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// MockAuthService for testing
type MockAuthService struct {
	mock.Mock
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

func (m *MockAuthService) Logout(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful login",
			setupMock: func(m *MockAuthService) {
				m.On("GetAuthURL", mock.AnythingOfType("string")).Return("https://discord.com/oauth2/authorize?client_id=test")
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "auth_url")
				assert.Contains(t, response, "state")
				assert.Equal(t, "https://discord.com/oauth2/authorize?client_id=test", response["auth_url"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			tt.setupMock(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := setupTestRouter()
			router.GET("/auth/login", handler.Login)

			req, _ := http.NewRequest("GET", "/auth/login", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Callback(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful callback",
			setupRequest: func(req *http.Request) {
				q := req.URL.Query()
				q.Add("code", "test_code")
				q.Add("state", "test_state")
				req.URL.RawQuery = q.Encode()
			},
			setupMock: func(m *MockAuthService) {
				authResponse := &models.AuthResponse{
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					ExpiresIn:    3600,
					User: models.User{
						ID:       "user_id",
						Username: "testuser",
					},
				}
				m.On("HandleCallback", mock.Anything, "test_code").Return(authResponse, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "access_token", response.AccessToken)
				assert.Equal(t, "refresh_token", response.RefreshToken)
			},
		},
		{
			name: "missing code parameter",
			setupRequest: func(req *http.Request) {
				q := req.URL.Query()
				q.Add("state", "test_state")
				req.URL.RawQuery = q.Encode()
			},
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "VALIDATION_ERROR", response.Code)
				assert.Contains(t, response.Message, "Authorization code is required")
			},
		},
		{
			name: "missing state parameter",
			setupRequest: func(req *http.Request) {
				q := req.URL.Query()
				q.Add("code", "test_code")
				req.URL.RawQuery = q.Encode()
			},
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "VALIDATION_ERROR", response.Code)
				assert.Contains(t, response.Message, "State parameter is required")
			},
		},
		{
			name: "discord api error",
			setupRequest: func(req *http.Request) {
				q := req.URL.Query()
				q.Add("code", "test_code")
				q.Add("state", "test_state")
				req.URL.RawQuery = q.Encode()
			},
			setupMock: func(m *MockAuthService) {
				m.On("HandleCallback", mock.Anything, "test_code").Return(nil, errors.New("discord api error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "DISCORD_API_ERROR", response.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			tt.setupMock(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := setupTestRouter()
			router.GET("/auth/callback", handler.Callback)

			req, _ := http.NewRequest("GET", "/auth/callback", nil)
			tt.setupRequest(req)
			
			// Set the state cookie for tests that need it
			if req.URL.Query().Get("state") == "test_state" {
				req.AddCookie(&http.Cookie{
					Name:  "oauth_state",
					Value: "test_state",
				})
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Refresh(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful refresh",
			requestBody: models.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			setupMock: func(m *MockAuthService) {
				authResponse := &models.AuthResponse{
					AccessToken:  "new_access_token",
					RefreshToken: "valid_refresh_token",
					ExpiresIn:    3600,
					User: models.User{
						ID:       "user_id",
						Username: "testuser",
					},
				}
				m.On("RefreshToken", mock.Anything, "valid_refresh_token").Return(authResponse, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "new_access_token", response.AccessToken)
			},
		},
		{
			name: "missing refresh token",
			requestBody: models.RefreshTokenRequest{
				RefreshToken: "",
			},
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "VALIDATION_ERROR", response.Code)
			},
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid json",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "VALIDATION_ERROR", response.Code)
			},
		},
		{
			name: "invalid refresh token",
			requestBody: models.RefreshTokenRequest{
				RefreshToken: "invalid_refresh_token",
			},
			setupMock: func(m *MockAuthService) {
				m.On("RefreshToken", mock.Anything, "invalid_refresh_token").Return(nil, errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "UNAUTHORIZED", response.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			tt.setupMock(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := setupTestRouter()
			router.POST("/auth/refresh", handler.Refresh)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful me request",
			setupContext: func(c *gin.Context) {
				user := &models.User{
					ID:       "user_id",
					Username: "testuser",
				}
				c.Set("user", user)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "user")
			},
		},
		{
			name:           "user not in context",
			setupContext:   func(c *gin.Context) {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "UNAUTHORIZED", response.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			handler := NewAuthHandler(mockAuthService)
			router := setupTestRouter()
			
			router.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})
			
			router.GET("/auth/me", handler.Me)

			req, _ := http.NewRequest("GET", "/auth/me", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		setupMock      func(*MockAuthService)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:       "successful logout",
			authHeader: "Bearer valid_token",
			setupMock: func(m *MockAuthService) {
				m.On("Logout", mock.Anything, "valid_token").Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
			},
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "VALIDATION_ERROR", response.Code)
			},
		},
		{
			name:           "invalid authorization header format",
			authHeader:     "InvalidFormat",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "VALIDATION_ERROR", response.Code)
			},
		},
		{
			name:       "logout service error",
			authHeader: "Bearer valid_token",
			setupMock: func(m *MockAuthService) {
				m.On("Logout", mock.Anything, "valid_token").Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.APIError
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "INTERNAL_ERROR", response.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := new(MockAuthService)
			tt.setupMock(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := setupTestRouter()
			router.POST("/auth/logout", handler.Logout)

			req, _ := http.NewRequest("POST", "/auth/logout", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.checkResponse(t, w)
			mockAuthService.AssertExpectations(t)
		})
	}
}