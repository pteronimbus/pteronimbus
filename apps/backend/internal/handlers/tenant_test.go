package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// MockTenantService is a mock implementation of TenantServiceInterface
type MockTenantService struct {
	mock.Mock
}

func (m *MockTenantService) CreateTenant(ctx context.Context, discordGuild *models.DiscordGuild, ownerID string) (*models.Tenant, error) {
	args := m.Called(ctx, discordGuild, ownerID)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantService) GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantService) GetTenantByDiscordServerID(ctx context.Context, discordServerID string) (*models.Tenant, error) {
	args := m.Called(ctx, discordServerID)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantService) GetUserTenants(ctx context.Context, userID string) ([]models.Tenant, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Tenant), args.Error(1)
}

func (m *MockTenantService) AddUserToTenant(ctx context.Context, userID, tenantID string, roles []string, permissions []string) error {
	args := m.Called(ctx, userID, tenantID, roles, permissions)
	return args.Error(0)
}

func (m *MockTenantService) RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error {
	args := m.Called(ctx, userID, tenantID)
	return args.Error(0)
}

func (m *MockTenantService) SyncDiscordRoles(ctx context.Context, tenantID, botToken string) error {
	args := m.Called(ctx, tenantID, botToken)
	return args.Error(0)
}

func (m *MockTenantService) SyncDiscordUsers(ctx context.Context, tenantID, botToken string) error {
	args := m.Called(ctx, tenantID, botToken)
	return args.Error(0)
}

func (m *MockTenantService) UpdateTenantConfig(ctx context.Context, tenantID string, config models.TenantConfig) error {
	args := m.Called(ctx, tenantID, config)
	return args.Error(0)
}

func (m *MockTenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	args := m.Called(ctx, tenantID)
	return args.Error(0)
}

func (m *MockTenantService) HasPermission(ctx context.Context, userID, tenantID, permission string) (bool, error) {
	args := m.Called(ctx, userID, tenantID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockTenantService) CheckManageServerPermission(discordGuild *models.DiscordGuild) bool {
	args := m.Called(discordGuild)
	return args.Bool(0)
}

// MockDiscordServiceForHandler is a mock for the Discord service used in handlers
type MockDiscordServiceForHandler struct {
	mock.Mock
}

func (m *MockDiscordServiceForHandler) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockDiscordServiceForHandler) ExchangeCodeForToken(ctx context.Context, code string) (*models.DiscordTokenResponse, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.DiscordTokenResponse), args.Error(1)
}

func (m *MockDiscordServiceForHandler) GetUserInfo(ctx context.Context, accessToken string) (*models.DiscordUser, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(*models.DiscordUser), args.Error(1)
}

func (m *MockDiscordServiceForHandler) RefreshToken(ctx context.Context, refreshToken string) (*models.DiscordTokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*models.DiscordTokenResponse), args.Error(1)
}

func (m *MockDiscordServiceForHandler) GetUserGuilds(ctx context.Context, accessToken string) ([]models.DiscordGuild, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).([]models.DiscordGuild), args.Error(1)
}

func (m *MockDiscordServiceForHandler) GetGuildRoles(ctx context.Context, botToken, guildID string) ([]models.DiscordRole, error) {
	args := m.Called(ctx, botToken, guildID)
	return args.Get(0).([]models.DiscordRole), args.Error(1)
}

func (m *MockDiscordServiceForHandler) GetGuildMembers(ctx context.Context, botToken, guildID string, limit int) ([]models.DiscordMember, error) {
	args := m.Called(ctx, botToken, guildID, limit)
	return args.Get(0).([]models.DiscordMember), args.Error(1)
}

func (m *MockDiscordServiceForHandler) GetGuildMember(ctx context.Context, botToken, guildID, userID string) (*models.DiscordMember, error) {
	args := m.Called(ctx, botToken, guildID, userID)
	return args.Get(0).(*models.DiscordMember), args.Error(1)
}

// MockAuthServiceForTenant is a mock for the Auth service used in tenant handlers
type MockAuthServiceForTenant struct {
	mock.Mock
}

func (m *MockAuthServiceForTenant) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *MockAuthServiceForTenant) HandleCallback(ctx context.Context, code string) (*models.AuthResponse, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthServiceForTenant) RefreshToken(ctx context.Context, refreshTokenString string) (*models.AuthResponse, error) {
	args := m.Called(ctx, refreshTokenString)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthServiceForTenant) ValidateAccessToken(ctx context.Context, accessToken string) (*models.User, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthServiceForTenant) ParseTokenClaims(accessToken string) (*models.JWTClaims, error) {
	args := m.Called(accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.JWTClaims), args.Error(1)
}

func (m *MockAuthServiceForTenant) Logout(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}

// MockRedisServiceForTenant is a mock for the Redis service used in tenant handlers
type MockRedisServiceForTenant struct {
	mock.Mock
}

func (m *MockRedisServiceForTenant) StoreSession(ctx context.Context, session *models.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRedisServiceForTenant) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockRedisServiceForTenant) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockRedisServiceForTenant) DeleteSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func setupTenantHandler() (*TenantHandler, *MockTenantService, *MockDiscordServiceForHandler, *MockAuthServiceForTenant, *MockRedisServiceForTenant) {
	mockTenantService := new(MockTenantService)
	mockDiscordService := new(MockDiscordServiceForHandler)
	mockAuthService := new(MockAuthServiceForTenant)
	mockRedisService := new(MockRedisServiceForTenant)
	
	handler := NewTenantHandler(mockTenantService, mockDiscordService, mockAuthService, mockRedisService)
	
	return handler, mockTenantService, mockDiscordService, mockAuthService, mockRedisService
}

func setupGinContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}
	
	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	return c, w
}

func TestTenantHandler_GetUserTenants(t *testing.T) {
	handler, mockTenantService, _, _, _ := setupTenantHandler()
	
	// Setup test data
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	tenants := []models.Tenant{
		{
			ID:              "tenant-1",
			DiscordServerID: "guild-1",
			Name:            "Test Guild 1",
			OwnerID:         user.ID,
		},
		{
			ID:              "tenant-2",
			DiscordServerID: "guild-2",
			Name:            "Test Guild 2",
			OwnerID:         "other-user",
		},
	}
	
	c, w := setupGinContext("GET", "/api/tenants", nil)
	c.Set("user", user)
	
	mockTenantService.On("GetUserTenants", mock.Anything, user.ID).Return(tenants, nil)
	
	// Execute
	handler.GetUserTenants(c)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	tenantsData := response["tenants"].([]interface{})
	assert.Len(t, tenantsData, 2)
	
	mockTenantService.AssertExpectations(t)
}

func TestTenantHandler_GetUserTenants_Unauthorized(t *testing.T) {
	handler, _, _, _, _ := setupTenantHandler()
	
	c, w := setupGinContext("GET", "/api/tenants", nil)
	// Don't set user in context
	
	// Execute
	handler.GetUserTenants(c)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response models.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "UNAUTHORIZED", response.Code)
}

func TestTenantHandler_CreateTenant(t *testing.T) {
	handler, mockTenantService, mockDiscordService, _, mockRedisService := setupTenantHandler()
	
	// Setup test data
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	guildID := "guild-123"
	requestBody := map[string]string{
		"guild_id": guildID,
	}
	
	discordGuilds := []models.DiscordGuild{
		{
			ID:          guildID,
			Name:        "Test Guild",
			Icon:        "icon-hash",
			Owner:       true,
			Permissions: "2147483647", // All permissions
		},
	}
	
	expectedTenant := &models.Tenant{
		ID:              "tenant-123",
		DiscordServerID: guildID,
		Name:            "Test Guild",
		Icon:            "icon-hash",
		OwnerID:         user.ID,
	}
	
	c, w := setupGinContext("POST", "/api/tenants", requestBody)
	c.Set("user", user)
	c.Set("session_id", "session-123")
	
	// Mock Redis session retrieval
	mockSession := &models.Session{
		ID:                 "session-123",
		AccessToken:        "jwt_access_token",
		DiscordAccessToken: "discord_access_token",
	}
	mockRedisService.On("GetSession", mock.Anything, "session-123").Return(mockSession, nil)
	mockDiscordService.On("GetUserGuilds", mock.Anything, "discord_access_token").Return(discordGuilds, nil)
	mockTenantService.On("CheckManageServerPermission", &discordGuilds[0]).Return(true)
	mockTenantService.On("CreateTenant", mock.Anything, &discordGuilds[0], user.ID).Return(expectedTenant, nil)
	mockTenantService.On("AddUserToTenant", mock.Anything, user.ID, expectedTenant.ID, []string{"owner"}, []string{"*"}).Return(nil)
	
	// Execute
	handler.CreateTenant(c)
	
	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	tenantData := response["tenant"].(map[string]interface{})
	assert.Equal(t, expectedTenant.ID, tenantData["id"])
	assert.Equal(t, expectedTenant.Name, tenantData["name"])
	
	mockDiscordService.AssertExpectations(t)
	mockTenantService.AssertExpectations(t)
}

func TestTenantHandler_CreateTenant_InvalidGuild(t *testing.T) {
	handler, _, mockDiscordService, _, mockRedisService := setupTenantHandler()
	
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	requestBody := map[string]string{
		"guild_id": "invalid-guild",
	}
	
	discordGuilds := []models.DiscordGuild{
		{
			ID:   "other-guild",
			Name: "Other Guild",
		},
	}
	
	c, w := setupGinContext("POST", "/api/tenants", requestBody)
	c.Set("user", user)
	c.Set("session_id", "session-123")
	
	// Mock Redis session retrieval
	mockSession := &models.Session{
		ID:                 "session-123",
		AccessToken:        "jwt_access_token",
		DiscordAccessToken: "discord_access_token",
	}
	mockRedisService.On("GetSession", mock.Anything, "session-123").Return(mockSession, nil)
	mockDiscordService.On("GetUserGuilds", mock.Anything, "discord_access_token").Return(discordGuilds, nil)
	
	// Execute
	handler.CreateTenant(c)
	
	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response models.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "FORBIDDEN", response.Code)
	
	mockDiscordService.AssertExpectations(t)
}

func TestTenantHandler_GetTenant(t *testing.T) {
	handler, mockTenantService, _, _, _ := setupTenantHandler()
	
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	tenant := &models.Tenant{
		ID:              "tenant-123",
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         user.ID,
	}
	
	c, w := setupGinContext("GET", "/api/tenants/tenant-123", nil)
	c.Set("user", user)
	c.Params = gin.Params{
		{Key: "id", Value: "tenant-123"},
	}
	
	mockTenantService.On("HasPermission", mock.Anything, user.ID, "tenant-123", "read").Return(true, nil)
	mockTenantService.On("GetTenant", mock.Anything, "tenant-123").Return(tenant, nil)
	
	// Execute
	handler.GetTenant(c)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	tenantData := response["tenant"].(map[string]interface{})
	assert.Equal(t, tenant.ID, tenantData["id"])
	assert.Equal(t, tenant.Name, tenantData["name"])
	
	mockTenantService.AssertExpectations(t)
}

func TestTenantHandler_GetTenant_Forbidden(t *testing.T) {
	handler, mockTenantService, _, _, _ := setupTenantHandler()
	
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	c, w := setupGinContext("GET", "/api/tenants/tenant-123", nil)
	c.Set("user", user)
	c.Params = gin.Params{
		{Key: "id", Value: "tenant-123"},
	}
	
	mockTenantService.On("HasPermission", mock.Anything, user.ID, "tenant-123", "read").Return(false, nil)
	
	// Execute
	handler.GetTenant(c)
	
	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response models.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "FORBIDDEN", response.Code)
	
	mockTenantService.AssertExpectations(t)
}

func TestTenantHandler_UpdateTenantConfig(t *testing.T) {
	handler, mockTenantService, _, _, _ := setupTenantHandler()
	
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	config := models.TenantConfig{
		ResourceLimits: models.ResourceLimits{
			MaxGameServers: 10,
			MaxCPU:         "4",
			MaxMemory:      "8Gi",
			MaxStorage:     "20Gi",
		},
		Settings: map[string]string{
			"notification_channel": "general",
		},
	}
	
	c, w := setupGinContext("PUT", "/api/tenants/tenant-123/config", config)
	c.Set("user", user)
	c.Params = gin.Params{
		{Key: "id", Value: "tenant-123"},
	}
	
	mockTenantService.On("HasPermission", mock.Anything, user.ID, "tenant-123", "manage").Return(true, nil)
	mockTenantService.On("UpdateTenantConfig", mock.Anything, "tenant-123", config).Return(nil)
	
	// Execute
	handler.UpdateTenantConfig(c)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "updated successfully")
	
	mockTenantService.AssertExpectations(t)
}

func TestTenantHandler_DeleteTenant(t *testing.T) {
	handler, mockTenantService, _, _, _ := setupTenantHandler()
	
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	tenant := &models.Tenant{
		ID:              "tenant-123",
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         user.ID,
	}
	
	c, w := setupGinContext("DELETE", "/api/tenants/tenant-123", nil)
	c.Set("user", user)
	c.Params = gin.Params{
		{Key: "id", Value: "tenant-123"},
	}
	
	mockTenantService.On("GetTenant", mock.Anything, "tenant-123").Return(tenant, nil)
	mockTenantService.On("DeleteTenant", mock.Anything, "tenant-123").Return(nil)
	
	// Execute
	handler.DeleteTenant(c)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["message"], "deleted successfully")
	
	mockTenantService.AssertExpectations(t)
}

func TestTenantHandler_DeleteTenant_NotOwner(t *testing.T) {
	handler, mockTenantService, _, _, _ := setupTenantHandler()
	
	user := &models.User{
		ID:            "user-123",
		DiscordUserID: "discord-123",
		Username:      "testuser",
	}
	
	tenant := &models.Tenant{
		ID:              "tenant-123",
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         "other-user", // Different owner
	}
	
	c, w := setupGinContext("DELETE", "/api/tenants/tenant-123", nil)
	c.Set("user", user)
	c.Params = gin.Params{
		{Key: "id", Value: "tenant-123"},
	}
	
	mockTenantService.On("GetTenant", mock.Anything, "tenant-123").Return(tenant, nil)
	
	// Execute
	handler.DeleteTenant(c)
	
	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response models.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "FORBIDDEN", response.Code)
	
	mockTenantService.AssertExpectations(t)
}