package services

import (
	"context"
	"testing"
	"time"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/testutils"
)

func setupGameServerTestDB(t *testing.T) (*gorm.DB, func()) {
	return testutils.SetupTestDatabaseWithModels(t,
		&models.GameServer{},
		&models.Tenant{},
		&models.User{},
		&models.UserTenant{},
		&models.TenantDiscordRole{},
		&models.TenantDiscordUser{},
		&models.Session{},
	)
}

func TestGameServerService_GetTenantServers(t *testing.T) {
	db, cleanup := setupGameServerTestDB(t)
	defer cleanup()
	service := NewGameServerService(db)
	ctx := context.Background()

	// Test getting servers for a tenant
	servers, err := service.GetTenantServers(ctx, "tenant-123")

	assert.NoError(t, err)
	assert.NotNil(t, servers)
	assert.Len(t, servers, 2) // Mock data returns 2 servers

	// Verify the mock data structure
	assert.Equal(t, "server-1", servers[0].ID)
	assert.Equal(t, "tenant-123", servers[0].TenantID)
	assert.Equal(t, "Survival World", servers[0].Name)
	assert.Equal(t, "minecraft", servers[0].GameType)
	assert.Equal(t, "Running", servers[0].Status.Phase)
	assert.Equal(t, 5, servers[0].Status.PlayerCount)

	assert.Equal(t, "server-2", servers[1].ID)
	assert.Equal(t, "tenant-123", servers[1].TenantID)
	assert.Equal(t, "Competitive Server", servers[1].Name)
	assert.Equal(t, "cs2", servers[1].GameType)
	assert.Equal(t, "Stopped", servers[1].Status.Phase)
	assert.Equal(t, 0, servers[1].Status.PlayerCount)
}

func TestGameServerService_GetTenantActivity(t *testing.T) {
	db, cleanup := setupGameServerTestDB(t)
	defer cleanup()
	service := NewGameServerService(db)
	ctx := context.Background()

	// Test getting activity without limit
	activities, err := service.GetTenantActivity(ctx, "tenant-123", 0)

	assert.NoError(t, err)
	assert.NotNil(t, activities)
	assert.Len(t, activities, 5) // Mock data returns 5 activities

	// Verify the mock data structure
	assert.Equal(t, "activity-1", activities[0].ID)
	assert.Equal(t, "server_started", activities[0].Type)
	assert.Contains(t, activities[0].Message, "Server 'Survival World' was started")

	// Test getting activity with limit
	limitedActivities, err := service.GetTenantActivity(ctx, "tenant-123", 3)

	assert.NoError(t, err)
	assert.NotNil(t, limitedActivities)
	assert.Len(t, limitedActivities, 3) // Should be limited to 3

	// Test getting activity with limit larger than available
	largeLimit, err := service.GetTenantActivity(ctx, "tenant-123", 10)

	assert.NoError(t, err)
	assert.NotNil(t, largeLimit)
	assert.Len(t, largeLimit, 5) // Should return all 5 available
}

func TestGameServerService_GetTenantDiscordStats(t *testing.T) {
	db, cleanup := setupGameServerTestDB(t)
	defer cleanup()
	service := NewGameServerService(db)
	ctx := context.Background()

	// Create a proper tenant with UUID for the test
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         "user-123",
	}
	err := db.Create(tenant).Error
	assert.NoError(t, err)

	// Create some Discord roles and users for the tenant
	discordRole := &models.TenantDiscordRole{
		TenantID:      tenant.ID,
		DiscordRoleID: "role-123",
		Name:          "Admin",
	}
	err = db.Create(discordRole).Error
	assert.NoError(t, err)

	discordUser := &models.TenantDiscordUser{
		TenantID:      tenant.ID,
		DiscordUserID: "user-123",
		Username:      "testuser",
		LastSyncAt:    time.Now(),
	}
	err = db.Create(discordUser).Error
	assert.NoError(t, err)

	// Test getting Discord stats
	stats, err := service.GetTenantDiscordStats(ctx, tenant.ID)

	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Verify the actual data from database
	assert.Equal(t, 1, stats.MemberCount) // 1 user created
	assert.Equal(t, 1, stats.RoleCount)   // 1 role created
	assert.NotEmpty(t, stats.LastSync)
}

func TestGameServerService_GetTenantServers_DifferentTenant(t *testing.T) {
	db, cleanup := setupGameServerTestDB(t)
	defer cleanup()
	service := NewGameServerService(db)
	ctx := context.Background()

	// Test that different tenant IDs still work (mock data uses the passed tenantID)
	servers, err := service.GetTenantServers(ctx, "different-tenant")

	assert.NoError(t, err)
	assert.NotNil(t, servers)
	assert.Len(t, servers, 2)

	// Verify that the tenant ID is correctly set in the mock data
	assert.Equal(t, "different-tenant", servers[0].TenantID)
	assert.Equal(t, "different-tenant", servers[1].TenantID)
}

func TestGameServerService_ActivityTypes(t *testing.T) {
	db, cleanup := setupGameServerTestDB(t)
	defer cleanup()
	service := NewGameServerService(db)
	ctx := context.Background()

	activities, err := service.GetTenantActivity(ctx, "tenant-123", 0)

	assert.NoError(t, err)
	assert.NotNil(t, activities)

	// Verify different activity types are present
	activityTypes := make(map[string]bool)
	for _, activity := range activities {
		activityTypes[activity.Type] = true
	}

	assert.True(t, activityTypes["server_started"])
	assert.True(t, activityTypes["user_joined"])
	assert.True(t, activityTypes["server_stopped"])
	assert.True(t, activityTypes["server_created"])
	assert.True(t, activityTypes["role_updated"])
}