package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/handlers"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/middleware"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/testutils"
)

// TenantMockDiscordService for tenant integration tests
type TenantMockDiscordService struct {
	mock.Mock
}

func (m *TenantMockDiscordService) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *TenantMockDiscordService) ExchangeCodeForToken(ctx context.Context, code string) (*models.DiscordTokenResponse, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.DiscordTokenResponse), args.Error(1)
}

func (m *TenantMockDiscordService) GetUserInfo(ctx context.Context, accessToken string) (*models.DiscordUser, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(*models.DiscordUser), args.Error(1)
}

func (m *TenantMockDiscordService) RefreshToken(ctx context.Context, refreshToken string) (*models.DiscordTokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*models.DiscordTokenResponse), args.Error(1)
}

func (m *TenantMockDiscordService) GetUserGuilds(ctx context.Context, accessToken string) ([]models.DiscordGuild, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).([]models.DiscordGuild), args.Error(1)
}

func (m *TenantMockDiscordService) GetGuildRoles(ctx context.Context, botToken, guildID string) ([]models.DiscordRole, error) {
	args := m.Called(ctx, botToken, guildID)
	return args.Get(0).([]models.DiscordRole), args.Error(1)
}

func (m *TenantMockDiscordService) GetGuildMembers(ctx context.Context, botToken, guildID string, limit int) ([]models.DiscordMember, error) {
	args := m.Called(ctx, botToken, guildID, limit)
	return args.Get(0).([]models.DiscordMember), args.Error(1)
}

func (m *TenantMockDiscordService) GetGuildMember(ctx context.Context, botToken, guildID, userID string) (*models.DiscordMember, error) {
	args := m.Called(ctx, botToken, guildID, userID)
	return args.Get(0).(*models.DiscordMember), args.Error(1)
}

// TenantMockAuthService for tenant integration tests
type TenantMockAuthService struct {
	mock.Mock
}

func (m *TenantMockAuthService) GetAuthURL(state string) string {
	args := m.Called(state)
	return args.String(0)
}

func (m *TenantMockAuthService) HandleCallback(ctx context.Context, code string) (*models.AuthResponse, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *TenantMockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *TenantMockAuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*models.User, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *TenantMockAuthService) ParseTokenClaims(accessToken string) (*models.JWTClaims, error) {
	args := m.Called(accessToken)
	return args.Get(0).(*models.JWTClaims), args.Error(1)
}

func (m *TenantMockAuthService) Logout(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}

// setupIntegrationTest sets up a complete test environment
func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB, *TenantMockDiscordService, func()) {
	// Setup PostgreSQL test database with all required models
	db, cleanup := testutils.SetupTestDatabaseWithModels(t,
		&models.User{},
		&models.Tenant{},
		&models.UserTenant{},
		&models.TenantDiscordRole{},
		&models.TenantDiscordUser{},
		&models.GameServer{},
		&models.Session{},
	)

	// Setup mock services
	mockDiscordService := new(TenantMockDiscordService)
	
	// Create real services with mocked dependencies
	tenantService := services.NewTenantService(db, mockDiscordService)

	// Create mock auth and redis services (use existing MockRedisService)
	mockAuthService := &TenantMockAuthService{}
	mockRedisService := &MockRedisService{}

	// Setup handlers
	tenantHandler := handlers.NewTenantHandler(tenantService, mockDiscordService, mockAuthService, mockRedisService)

	// Setup middleware
	tenantMiddleware := middleware.NewTenantMiddleware(tenantService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add routes
	apiRoutes := router.Group("/api")
	{
		tenantRoutes := apiRoutes.Group("/tenants")
		{
			tenantRoutes.GET("", tenantHandler.GetUserTenants)
			tenantRoutes.POST("", tenantHandler.CreateTenant)
			tenantRoutes.GET("/:id", tenantHandler.GetTenant)
			tenantRoutes.PUT("/:id/config", tenantHandler.UpdateTenantConfig)
			tenantRoutes.DELETE("/:id", tenantHandler.DeleteTenant)
		}

		// Tenant-scoped routes
		tenantScopedRoutes := apiRoutes.Group("/tenant")
		tenantScopedRoutes.Use(tenantMiddleware.RequireTenant())
		{
			tenantScopedRoutes.GET("/info", func(c *gin.Context) {
				tenant, _ := c.Get("tenant")
				c.JSON(http.StatusOK, gin.H{
					"tenant": tenant,
				})
			})
		}
	}

	return router, db, mockDiscordService, cleanup
}

func TestTenantIntegration_ServiceInitialization(t *testing.T) {
	// Test that the integration setup works correctly
	router, db, mockDiscordService, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Verify that all components are properly initialized
	assert.NotNil(t, router)
	assert.NotNil(t, db)
	assert.NotNil(t, mockDiscordService)

	// Test that the database connection works
	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Ping()
	require.NoError(t, err)
}

func TestTenantIntegration_PermissionLogic(t *testing.T) {
	// Test the core permission checking logic without database dependencies
	tenantService := &services.TenantService{}

	// Test CheckManageServerPermission method
	testCases := []struct {
		name        string
		guild       *models.DiscordGuild
		expected    bool
		description string
	}{
		{
			name: "owner_has_permission",
			guild: &models.DiscordGuild{
				ID:          "guild-123",
				Name:        "Test Guild",
				Owner:       true,
				Permissions: "0",
			},
			expected:    true,
			description: "Guild owner should always have permission",
		},
		{
			name: "manage_guild_permission",
			guild: &models.DiscordGuild{
				ID:          "guild-123",
				Name:        "Test Guild",
				Owner:       false,
				Permissions: "32", // MANAGE_GUILD = 0x20 = 32
			},
			expected:    true,
			description: "User with MANAGE_GUILD permission should have access",
		},
		{
			name: "insufficient_permissions",
			guild: &models.DiscordGuild{
				ID:          "guild-123",
				Name:        "Test Guild",
				Owner:       false,
				Permissions: "1", // Basic permissions only
			},
			expected:    false,
			description: "User without MANAGE_GUILD permission should be denied",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tenantService.CheckManageServerPermission(tc.guild)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}