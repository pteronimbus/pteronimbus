package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGameServerService is a mock implementation of GameServerServiceInterface
type MockGameServerService struct {
	mock.Mock
}

func (m *MockGameServerService) GetTenantServers(ctx context.Context, tenantID string) ([]models.GameServer, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]models.GameServer), args.Error(1)
}

func (m *MockGameServerService) GetTenantActivity(ctx context.Context, tenantID string, limit int) ([]models.Activity, error) {
	args := m.Called(ctx, tenantID, limit)
	return args.Get(0).([]models.Activity), args.Error(1)
}

func (m *MockGameServerService) GetTenantDiscordStats(ctx context.Context, tenantID string) (*models.DiscordStats, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(*models.DiscordStats), args.Error(1)
}

// MockTenantService is a mock implementation of TenantServiceInterface for game server tests
type MockTenantServiceForGameServer struct {
	mock.Mock
}

func (m *MockTenantServiceForGameServer) CreateTenant(ctx context.Context, discordGuild *models.DiscordGuild, ownerID string) (*models.Tenant, error) {
	args := m.Called(ctx, discordGuild, ownerID)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantServiceForGameServer) GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantServiceForGameServer) GetTenantByDiscordServerID(ctx context.Context, discordServerID string) (*models.Tenant, error) {
	args := m.Called(ctx, discordServerID)
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockTenantServiceForGameServer) GetUserTenants(ctx context.Context, userID string) ([]models.Tenant, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Tenant), args.Error(1)
}

func (m *MockTenantServiceForGameServer) AddUserToTenant(ctx context.Context, userID, tenantID string, roles []string, permissions []string) error {
	args := m.Called(ctx, userID, tenantID, roles, permissions)
	return args.Error(0)
}

func (m *MockTenantServiceForGameServer) RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error {
	args := m.Called(ctx, userID, tenantID)
	return args.Error(0)
}

func (m *MockTenantServiceForGameServer) SyncDiscordRoles(ctx context.Context, tenantID, botToken string) error {
	args := m.Called(ctx, tenantID, botToken)
	return args.Error(0)
}

func (m *MockTenantServiceForGameServer) SyncDiscordUsers(ctx context.Context, tenantID, botToken string) error {
	args := m.Called(ctx, tenantID, botToken)
	return args.Error(0)
}

func (m *MockTenantServiceForGameServer) UpdateTenantConfig(ctx context.Context, tenantID string, config models.TenantConfig) error {
	args := m.Called(ctx, tenantID, config)
	return args.Error(0)
}

func (m *MockTenantServiceForGameServer) DeleteTenant(ctx context.Context, tenantID string) error {
	args := m.Called(ctx, tenantID)
	return args.Error(0)
}

func (m *MockTenantServiceForGameServer) HasPermission(ctx context.Context, userID, tenantID, permission string) (bool, error) {
	args := m.Called(ctx, userID, tenantID, permission)
	return args.Bool(0), args.Error(1)
}

func (m *MockTenantServiceForGameServer) CheckManageServerPermission(discordGuild *models.DiscordGuild) bool {
	args := m.Called(discordGuild)
	return args.Bool(0)
}

func setupGameServerHandler() (*GameServerHandler, *MockGameServerService, *MockTenantServiceForGameServer) {
	mockGameServerService := &MockGameServerService{}
	mockTenantService := &MockTenantServiceForGameServer{}
	handler := NewGameServerHandler(mockGameServerService, mockTenantService)
	return handler, mockGameServerService, mockTenantService
}

func setupGinContextForGameServer(method, url string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}

func TestGetTenantServers_Success(t *testing.T) {
	handler, mockGameServerService, _ := setupGameServerHandler()

	tenant := &models.Tenant{
		ID:              "tenant-123",
		DiscordServerID: "discord-123",
		Name:            "Test Server",
	}

	expectedServers := []models.GameServer{
		{
			ID:         "server-1",
			TenantID:   "tenant-123",
			TemplateID: "template-minecraft",
			Name:       "Survival World",
			GameType:   "minecraft",
			Status: models.GameServerStatus{
				Phase:       "Running",
				Message:     "Server is running normally",
				PlayerCount: 5,
			},
		},
	}

	c, w := setupGinContextForGameServer("GET", "/api/tenant/servers", nil)
	c.Set("tenant", tenant)

	mockGameServerService.On("GetTenantServers", mock.Anything, "tenant-123").Return(expectedServers, nil)

	handler.GetTenantServers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	servers, exists := response["servers"]
	assert.True(t, exists)
	assert.NotNil(t, servers)

	mockGameServerService.AssertExpectations(t)
}

func TestGetTenantServers_NoTenantContext(t *testing.T) {
	handler, _, _ := setupGameServerHandler()

	c, w := setupGinContextForGameServer("GET", "/api/tenant/servers", nil)
	// Don't set tenant in context

	handler.GetTenantServers(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "TENANT_REQUIRED", response.Code)
	assert.Equal(t, "Tenant context is required", response.Message)
}

func TestGetTenantActivity_Success(t *testing.T) {
	handler, mockGameServerService, _ := setupGameServerHandler()

	tenant := &models.Tenant{
		ID:              "tenant-123",
		DiscordServerID: "discord-123",
		Name:            "Test Server",
	}

	expectedActivity := []models.Activity{
		{
			ID:        "activity-1",
			Type:      "server_started",
			Message:   "Server 'Survival World' was started",
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	c, w := setupGinContextForGameServer("GET", "/api/tenant/activity?limit=10", nil)
	c.Set("tenant", tenant)

	mockGameServerService.On("GetTenantActivity", mock.Anything, "tenant-123", 10).Return(expectedActivity, nil)

	handler.GetTenantActivity(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	activities, exists := response["activities"]
	assert.True(t, exists)
	assert.NotNil(t, activities)

	mockGameServerService.AssertExpectations(t)
}

func TestGetTenantActivity_NoTenantContext(t *testing.T) {
	handler, _, _ := setupGameServerHandler()

	c, w := setupGinContextForGameServer("GET", "/api/tenant/activity", nil)
	// Don't set tenant in context

	handler.GetTenantActivity(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "TENANT_REQUIRED", response.Code)
	assert.Equal(t, "Tenant context is required", response.Message)
}

func TestGetTenantDiscordStats_Success(t *testing.T) {
	handler, mockGameServerService, _ := setupGameServerHandler()

	tenant := &models.Tenant{
		ID:              "tenant-123",
		DiscordServerID: "discord-123",
		Name:            "Test Server",
	}

	expectedStats := &models.DiscordStats{
		MemberCount: 42,
		RoleCount:   8,
		LastSync:    time.Now().Format(time.RFC3339),
	}

	c, w := setupGinContextForGameServer("GET", "/api/tenant/discord/stats", nil)
	c.Set("tenant", tenant)

	mockGameServerService.On("GetTenantDiscordStats", mock.Anything, "tenant-123").Return(expectedStats, nil)

	handler.GetTenantDiscordStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	stats, exists := response["stats"]
	assert.True(t, exists)
	assert.NotNil(t, stats)

	mockGameServerService.AssertExpectations(t)
}

func TestGetTenantDiscordStats_NoTenantContext(t *testing.T) {
	handler, _, _ := setupGameServerHandler()

	c, w := setupGinContextForGameServer("GET", "/api/tenant/discord/stats", nil)
	// Don't set tenant in context

	handler.GetTenantDiscordStats(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "TENANT_REQUIRED", response.Code)
	assert.Equal(t, "Tenant context is required", response.Message)
}