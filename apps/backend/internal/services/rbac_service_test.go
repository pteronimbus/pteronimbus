package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupRBACTest(t *testing.T) (*RBACService, *gorm.DB, func()) {
	db, cleanup := testutils.SetupTestDatabaseWithModels(t,
		&models.User{},
		&models.Tenant{},
		&models.UserTenant{},
		&models.TenantDiscordRole{},
		&models.TenantDiscordUser{},
		&models.Permission{},
		&models.Role{},
		&models.PermissionAuditLog{},
		&models.GuildMembershipCache{},
	)
	
	rbacConfig := &config.RBACConfig{
		SuperAdminDiscordID: "superadmin123",
		RoleSyncTTL:         time.Minute * 5,
		GuildCacheTTL:       time.Minute * 5,
		GracePeriod:         time.Minute * 2,
	}
	
	rbacService := NewRBACService(db, rbacConfig)
	
	return rbacService, db, cleanup
}

func TestRBACService_HasPermission(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create Discord role
	discordRole := &models.TenantDiscordRole{
		TenantID:      tenant.ID,
		DiscordRoleID: "admin",
		Name:          "Admin",
		Permissions:   models.StringArray{models.PermissionServerRead, models.PermissionServerWrite},
	}
	err = db.Create(discordRole).Error
	require.NoError(t, err)

	// Create user-tenant relationship
	userTenant := &models.UserTenant{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		Roles:       models.StringArray{"admin"},
		Permissions: models.StringArray{models.PermissionServerCreate},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Test direct permissions
	hasPermission, err := rbacService.HasPermission(context.Background(), user.ID, tenant.ID, models.PermissionServerCreate)
	require.NoError(t, err)
	assert.True(t, hasPermission)

	// Test role-based permissions
	hasPermission, err = rbacService.HasPermission(context.Background(), user.ID, tenant.ID, models.PermissionServerRead)
	require.NoError(t, err)
	assert.True(t, hasPermission)

	// Test permission user doesn't have
	hasPermission, err = rbacService.HasPermission(context.Background(), user.ID, tenant.ID, models.PermissionServerDelete)
	require.NoError(t, err)
	assert.False(t, hasPermission)

	// Test wildcard permission
	userTenant.Permissions = models.StringArray{models.PermissionAdminAll}
	err = db.Save(userTenant).Error
	require.NoError(t, err)

	hasPermission, err = rbacService.HasPermission(context.Background(), user.ID, tenant.ID, "any:permission")
	require.NoError(t, err)
	assert.True(t, hasPermission)
}

func TestRBACService_IsSuperAdmin(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test tenant
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	// Test super admin by Discord ID (initial setup)
	superAdmin := &models.User{
		DiscordUserID: "superadmin123", // Matches config
		Username:      "superadmin",
	}
	err = db.Create(superAdmin).Error
	require.NoError(t, err)

	isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), superAdmin.ID)
	require.NoError(t, err)
	assert.True(t, isSuperAdmin)

	// Test super admin by role assignment
	roleAdmin := &models.User{
		DiscordUserID: "roleadmin123",
		Username:      "roleadmin",
	}
	err = db.Create(roleAdmin).Error
	require.NoError(t, err)

	// Assign super admin role
	userTenant := &models.UserTenant{
		UserID:   roleAdmin.ID,
		TenantID: tenant.ID,
		Roles:    models.StringArray{"superadmin"},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	isSuperAdmin, err = rbacService.IsSuperAdmin(context.Background(), roleAdmin.ID)
	require.NoError(t, err)
	assert.True(t, isSuperAdmin)

	// Test super admin by wildcard permission
	permAdmin := &models.User{
		DiscordUserID: "permadmin123",
		Username:      "permadmin",
	}
	err = db.Create(permAdmin).Error
	require.NoError(t, err)

	// Assign wildcard permission
	permUserTenant := &models.UserTenant{
		UserID:      permAdmin.ID,
		TenantID:    tenant.ID,
		Permissions: models.StringArray{models.PermissionAdminAll},
	}
	err = db.Create(permUserTenant).Error
	require.NoError(t, err)

	isSuperAdmin, err = rbacService.IsSuperAdmin(context.Background(), permAdmin.ID)
	require.NoError(t, err)
	assert.True(t, isSuperAdmin)

	// Test regular user check
	regularUser := &models.User{
		DiscordUserID: "regular123",
		Username:      "regularuser",
	}
	err = db.Create(regularUser).Error
	require.NoError(t, err)

	isSuperAdmin, err = rbacService.IsSuperAdmin(context.Background(), regularUser.ID)
	require.NoError(t, err)
	assert.False(t, isSuperAdmin)

	// Test non-existent user
	isSuperAdmin, err = rbacService.IsSuperAdmin(context.Background(), uuid.New().String())
	require.NoError(t, err)
	assert.False(t, isSuperAdmin)
}

func TestRBACService_CreateRole(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create tenant
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	// Create role
	permissions := []string{models.PermissionServerRead, models.PermissionServerWrite}
	role, err := rbacService.CreateRole(context.Background(), tenant.ID, "TestRole", permissions, false)
	require.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, tenant.ID, role.TenantID)
	assert.Equal(t, "TestRole", role.Name)
	assert.Equal(t, models.StringArray(permissions), role.Permissions)
	assert.False(t, role.IsSystemRole)

	// Verify role was created in database
	var dbRole models.Role
	err = db.Where("id = ?", role.ID).First(&dbRole).Error
	require.NoError(t, err)
	assert.Equal(t, "TestRole", dbRole.Name)
}

func TestRBACService_UpdateRole(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create tenant and role
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	role, err := rbacService.CreateRole(context.Background(), tenant.ID, "TestRole", []string{models.PermissionServerRead}, false)
	require.NoError(t, err)

	// Update role
	newPermissions := []string{models.PermissionServerRead, models.PermissionServerWrite}
	updatedRole, err := rbacService.UpdateRole(context.Background(), role.ID, "UpdatedRole", newPermissions)
	require.NoError(t, err)
	assert.Equal(t, "UpdatedRole", updatedRole.Name)
	assert.Equal(t, models.StringArray(newPermissions), updatedRole.Permissions)
}

func TestRBACService_DeleteRole(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create tenant and role
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	role, err := rbacService.CreateRole(context.Background(), tenant.ID, "TestRole", []string{models.PermissionServerRead}, false)
	require.NoError(t, err)

	// Delete role
	err = rbacService.DeleteRole(context.Background(), role.ID)
	require.NoError(t, err)

	// Verify role was deleted
	var dbRole models.Role
	err = db.Where("id = ?", role.ID).First(&dbRole).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestRBACService_DeleteSystemRole(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create tenant and system role
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	role, err := rbacService.CreateRole(context.Background(), tenant.ID, "SystemRole", []string{models.PermissionServerRead}, true)
	require.NoError(t, err)

	// Try to delete system role (should fail)
	err = rbacService.DeleteRole(context.Background(), role.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot delete system role")
}

func TestRBACService_AssignRoleToUser(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Assign role to user
	err = rbacService.AssignRoleToUser(context.Background(), user.ID, tenant.ID, "admin")
	require.NoError(t, err)

	// Verify role was assigned
	var userTenant models.UserTenant
	err = db.Where("user_id = ? AND tenant_id = ?", user.ID, tenant.ID).First(&userTenant).Error
	require.NoError(t, err)
	assert.Contains(t, userTenant.Roles, "admin")

	// Assign same role again (should not duplicate)
	err = rbacService.AssignRoleToUser(context.Background(), user.ID, tenant.ID, "admin")
	require.NoError(t, err)

	// Verify no duplication
	err = db.Where("user_id = ? AND tenant_id = ?", user.ID, tenant.ID).First(&userTenant).Error
	require.NoError(t, err)
	assert.Len(t, userTenant.Roles, 1)
	assert.Contains(t, userTenant.Roles, "admin")
}

func TestRBACService_RemoveRoleFromUser(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create user-tenant relationship with roles
	userTenant := &models.UserTenant{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		Roles:       models.StringArray{"admin", "moderator"},
		Permissions: models.StringArray{},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Remove role from user
	err = rbacService.RemoveRoleFromUser(context.Background(), user.ID, tenant.ID, "admin")
	require.NoError(t, err)

	// Verify role was removed
	err = db.Where("user_id = ? AND tenant_id = ?", user.ID, tenant.ID).First(&userTenant).Error
	require.NoError(t, err)
	assert.NotContains(t, userTenant.Roles, "admin")
	assert.Contains(t, userTenant.Roles, "moderator")
}

func TestRBACService_GetUserPermissions(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create Discord role
	discordRole := &models.TenantDiscordRole{
		TenantID:      tenant.ID,
		DiscordRoleID: "admin",
		Name:          "Admin",
		Permissions:   models.StringArray{models.PermissionServerRead, models.PermissionServerWrite},
	}
	err = db.Create(discordRole).Error
	require.NoError(t, err)

	// Create user-tenant relationship
	userTenant := &models.UserTenant{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		Roles:       models.StringArray{"admin"},
		Permissions: models.StringArray{models.PermissionServerCreate},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Get user permissions
	permissions, err := rbacService.GetUserPermissions(context.Background(), user.ID, tenant.ID)
	require.NoError(t, err)
	assert.Len(t, permissions, 3)
	assert.Contains(t, permissions, models.PermissionServerCreate)
	assert.Contains(t, permissions, models.PermissionServerRead)
	assert.Contains(t, permissions, models.PermissionServerWrite)
}

func TestRBACService_GetUserPermissions_SuperAdmin(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create super admin user
	superAdmin := &models.User{
		DiscordUserID: "superadmin123", // Matches config
		Username:      "superadmin",
	}
	err := db.Create(superAdmin).Error
	require.NoError(t, err)

	// Create tenant
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err = db.Create(tenant).Error
	require.NoError(t, err)

	// Get super admin permissions
	permissions, err := rbacService.GetUserPermissions(context.Background(), superAdmin.ID, tenant.ID)
	require.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Contains(t, permissions, models.PermissionAdminAll)
}

func TestRBACService_LogPermissionChange(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	performer := &models.User{
		DiscordUserID: "performer-123",
		Username:      "performer",
	}
	err = db.Create(performer).Error
	require.NoError(t, err)

	// Log permission change
	err = rbacService.LogPermissionChange(
		context.Background(),
		user.ID,
		tenant.ID,
		"role_assigned",
		"user",
		user.ID,
		"",
		"admin",
		"User promoted to admin",
		performer.ID,
	)
	require.NoError(t, err)

	// Verify audit log was created
	var auditLog models.PermissionAuditLog
	err = db.Where("user_id = ? AND tenant_id = ?", user.ID, tenant.ID).First(&auditLog).Error
	require.NoError(t, err)
	assert.Equal(t, "role_assigned", auditLog.Action)
	assert.Equal(t, "user", auditLog.ResourceType)
	assert.Equal(t, user.ID, auditLog.ResourceID)
	assert.Equal(t, "admin", auditLog.NewValue)
	assert.Equal(t, "User promoted to admin", auditLog.Reason)
	assert.Equal(t, performer.ID, auditLog.PerformedBy)
}

func TestRBACService_GetRoles(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create tenant
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	// Create roles
	_, err = rbacService.CreateRole(context.Background(), tenant.ID, "Role1", []string{models.PermissionServerRead}, false)
	require.NoError(t, err)

	_, err = rbacService.CreateRole(context.Background(), tenant.ID, "Role2", []string{models.PermissionServerWrite}, false)
	require.NoError(t, err)

	// Get roles
	roles, err := rbacService.GetRoles(context.Background(), tenant.ID)
	require.NoError(t, err)
	assert.Len(t, roles, 2)
	assert.Equal(t, "Role1", roles[0].Name)
	assert.Equal(t, "Role2", roles[1].Name)
}

func TestRBACService_AssignSuperAdminRole(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test tenant
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	// Create super admin user (by Discord ID)
	superAdmin := &models.User{
		DiscordUserID: "superadmin123",
		Username:      "superadmin",
	}
	err = db.Create(superAdmin).Error
	require.NoError(t, err)

	// Create target user
	targetUser := &models.User{
		DiscordUserID: "target123",
		Username:      "targetuser",
	}
	err = db.Create(targetUser).Error
	require.NoError(t, err)

	// Create context with super admin user
	ctx := context.WithValue(context.Background(), "user_id", superAdmin.ID)

	// Assign super admin role
	err = rbacService.AssignSuperAdminRole(ctx, targetUser.ID, tenant.ID)
	require.NoError(t, err)

	// Verify role was assigned
	var userTenant models.UserTenant
	err = db.Where("user_id = ? AND tenant_id = ?", targetUser.ID, tenant.ID).First(&userTenant).Error
	require.NoError(t, err)
	assert.Contains(t, userTenant.Roles, "superadmin")

	// Verify target user is now super admin
	isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), targetUser.ID)
	require.NoError(t, err)
	assert.True(t, isSuperAdmin)

	// Test assigning to non-super admin (should fail)
	regularUser := &models.User{
		DiscordUserID: "regular123",
		Username:      "regularuser",
	}
	err = db.Create(regularUser).Error
	require.NoError(t, err)

	regularCtx := context.WithValue(context.Background(), "user_id", regularUser.ID)
	err = rbacService.AssignSuperAdminRole(regularCtx, targetUser.ID, tenant.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only super admins can assign super admin role")
}

func TestRBACService_RemoveSuperAdminRole(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test tenant
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	// Create super admin user (by Discord ID)
	superAdmin := &models.User{
		DiscordUserID: "superadmin123",
		Username:      "superadmin",
	}
	err = db.Create(superAdmin).Error
	require.NoError(t, err)

	// Create target user with super admin role
	targetUser := &models.User{
		DiscordUserID: "target123",
		Username:      "targetuser",
	}
	err = db.Create(targetUser).Error
	require.NoError(t, err)

	userTenant := &models.UserTenant{
		UserID:   targetUser.ID,
		TenantID: tenant.ID,
		Roles:    models.StringArray{"superadmin", "admin"},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Create context with super admin user
	ctx := context.WithValue(context.Background(), "user_id", superAdmin.ID)

	// Remove super admin role
	err = rbacService.RemoveSuperAdminRole(ctx, targetUser.ID, tenant.ID)
	require.NoError(t, err)

	// Verify role was removed
	err = db.Where("user_id = ? AND tenant_id = ?", targetUser.ID, tenant.ID).First(&userTenant).Error
	require.NoError(t, err)
	assert.NotContains(t, userTenant.Roles, "superadmin")
	assert.Contains(t, userTenant.Roles, "admin") // Other roles should remain

	// Verify target user is no longer super admin
	isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), targetUser.ID)
	require.NoError(t, err)
	assert.False(t, isSuperAdmin)

	// Test removing from initial super admin (should fail)
	err = rbacService.RemoveSuperAdminRole(ctx, superAdmin.ID, tenant.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot remove super admin role from the initial super admin")

	// Test removing by non-super admin (should fail)
	regularUser := &models.User{
		DiscordUserID: "regular123",
		Username:      "regularuser",
	}
	err = db.Create(regularUser).Error
	require.NoError(t, err)

	regularCtx := context.WithValue(context.Background(), "user_id", regularUser.ID)
	err = rbacService.RemoveSuperAdminRole(regularCtx, targetUser.ID, tenant.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only super admins can remove super admin role")
} 