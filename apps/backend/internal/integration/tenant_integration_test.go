package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/handlers"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/middleware"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
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

// setupIntegrationTest sets up a complete test environment
func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB, *TenantMockDiscordService) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Create simplified models for SQLite testing
	type TestUser struct {
		ID            string `gorm:"primaryKey"`
		DiscordUserID string `gorm:"uniqueIndex;not null"`
		Username      string `gorm:"not null"`
		Avatar        string
		Email         string
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	type TestTenant struct {
		ID              string `gorm:"primaryKey"`
		DiscordServerID string `gorm:"uniqueIndex;not null"`
		Name            string `gorm:"not null"`
		Icon            string
		OwnerID         string `gorm:"not null"`
		Config          string `gorm:"type:text"` // Store as JSON string
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}

	type TestUserTenant struct {
		ID          string `gorm:"primaryKey"`
		UserID      string `gorm:"not null;index"`
		TenantID    string `gorm:"not null;index"`
		Roles       string `gorm:"type:text"` // Store as JSON string
		Permissions string `gorm:"type:text"` // Store as JSON string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	type TestTenantDiscordRole struct {
		ID            string `gorm:"primaryKey"`
		TenantID      string `gorm:"not null;index"`
		DiscordRoleID string `gorm:"not null"`
		Name          string `gorm:"not null"`
		Color         int
		Position      int
		Permissions   string `gorm:"type:text"` // Store as JSON string
		Mentionable   bool
		Hoist         bool
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	// Auto-migrate simplified models
	err = db.AutoMigrate(
		&TestUser{},
		&TestTenant{},
		&TestUserTenant{},
		&TestTenantDiscordRole{},
	)
	assert.NoError(t, err)

	// Setup mock services
	mockDiscordService := new(TenantMockDiscordService)
	
	// Create real services with mocked dependencies
	tenantService := services.NewTenantService(db, mockDiscordService)

	// Setup handlers
	tenantHandler := handlers.NewTenantHandler(tenantService, mockDiscordService, nil)

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

	return router, db, mockDiscordService
}

func TestTenantIntegration_ServiceInitialization(t *testing.T) {
	// Test that the integration setup works correctly
	router, db, mockDiscordService := setupIntegrationTest(t)

	// Verify that all components are properly initialized
	assert.NotNil(t, router)
	assert.NotNil(t, db)
	assert.NotNil(t, mockDiscordService)

	// Test that the database connection works
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	err = sqlDB.Ping()
	assert.NoError(t, err)
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