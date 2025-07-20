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
	"github.com/stretchr/testify/mock"
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
		&models.SystemRole{},
		&models.UserSystemRole{},
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

	// Test super admin by system role assignment
	roleAdmin := &models.User{
		DiscordUserID: "roleadmin123",
		Username:      "roleadmin",
	}
	err = db.Create(roleAdmin).Error
	require.NoError(t, err)

	// Create super admin system role first
	_, err = rbacService.CreateSystemRole(context.Background(), "superadmin", "Super administrator", []string{models.PermissionSystemAdmin})
	require.NoError(t, err)

	// Assign super admin system role
	err = rbacService.AssignSystemRoleToUser(context.Background(), roleAdmin.ID, "superadmin")
	require.NoError(t, err)

	isSuperAdmin, err = rbacService.IsSuperAdmin(context.Background(), roleAdmin.ID)
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

// Removed old failing tests - replaced with TestRBACService_SuperAdminAutoAssignmentDuringSignup and TestRBACService_ManualSuperAdminRoleAssignment 

func TestRBACService_SystemRoles(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test users
	superAdmin := &models.User{
		DiscordUserID: "superadmin123",
		Username:      "superadmin",
	}
	err := db.Create(superAdmin).Error
	require.NoError(t, err)

	regularUser := &models.User{
		DiscordUserID: "regular123",
		Username:      "regularuser",
	}
	err = db.Create(regularUser).Error
	require.NoError(t, err)

	// Test creating system role
	t.Run("CreateSystemRole", func(t *testing.T) {
		systemRole, err := rbacService.CreateSystemRole(context.Background(), "test_admin", "Test admin role", []string{models.PermissionSystemAdmin})
		require.NoError(t, err)
		assert.Equal(t, "test_admin", systemRole.Name)
		assert.Equal(t, "Test admin role", systemRole.Description)
		assert.Contains(t, systemRole.Permissions, models.PermissionSystemAdmin)
	})

	// Test assigning system role to user
	t.Run("AssignSystemRoleToUser", func(t *testing.T) {
		// Create system role first
		_, err := rbacService.CreateSystemRole(context.Background(), "test_role", "Test role", []string{"test:permission"})
		require.NoError(t, err)

		// Assign role to user
		err = rbacService.AssignSystemRoleToUser(context.Background(), regularUser.ID, "test_role")
		require.NoError(t, err)

		// Verify role was assigned
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), regularUser.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "test_role", systemRoles[0].Name)
	})

	// Test removing system role from user
	t.Run("RemoveSystemRoleFromUser", func(t *testing.T) {
		// Create and assign role
		_, err := rbacService.CreateSystemRole(context.Background(), "remove_test", "Remove test", []string{"test:permission"})
		require.NoError(t, err)
		err = rbacService.AssignSystemRoleToUser(context.Background(), regularUser.ID, "remove_test")
		require.NoError(t, err)

		// Remove role
		err = rbacService.RemoveSystemRoleFromUser(context.Background(), regularUser.ID, "remove_test")
		require.NoError(t, err)

		// Verify role was removed
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), regularUser.ID)
		require.NoError(t, err)
		// Should only have the test_role from previous test
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "test_role", systemRoles[0].Name)
	})

	// Test system permission checking
	t.Run("HasSystemPermission", func(t *testing.T) {
		// Create system role with specific permission
		_, err := rbacService.CreateSystemRole(context.Background(), "permission_test", "Permission test", []string{"custom:permission"})
		require.NoError(t, err)

		// Assign role to user
		err = rbacService.AssignSystemRoleToUser(context.Background(), regularUser.ID, "permission_test")
		require.NoError(t, err)

		// Check permission
		hasPermission, err := rbacService.HasSystemPermission(context.Background(), regularUser.ID, "custom:permission")
		require.NoError(t, err)
		assert.True(t, hasPermission)

		// Check non-existent permission
		hasPermission, err = rbacService.HasSystemPermission(context.Background(), regularUser.ID, "nonexistent:permission")
		require.NoError(t, err)
		assert.False(t, hasPermission)
	})

	// Test super admin through system role
	t.Run("SuperAdminThroughSystemRole", func(t *testing.T) {
		// Create super admin system role
		_, err := rbacService.CreateSystemRole(context.Background(), "superadmin", "Super administrator", []string{models.PermissionSystemAdmin})
		require.NoError(t, err)

		// Assign super admin role to user
		err = rbacService.AssignSystemRoleToUser(context.Background(), regularUser.ID, "superadmin")
		require.NoError(t, err)

		// Check if user is super admin
		isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), regularUser.ID)
		require.NoError(t, err)
		assert.True(t, isSuperAdmin)

		// Check system admin permission
		hasPermission, err := rbacService.HasSystemPermission(context.Background(), regularUser.ID, models.PermissionSystemAdmin)
		require.NoError(t, err)
		assert.True(t, hasPermission)
	})

	// Test getting all system roles
	t.Run("GetSystemRoles", func(t *testing.T) {
		systemRoles, err := rbacService.GetSystemRoles(context.Background())
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(systemRoles), 3) // Should have at least the roles we created
		
		// Check for specific roles
		roleNames := make([]string, len(systemRoles))
		for i, role := range systemRoles {
			roleNames[i] = role.Name
		}
		assert.Contains(t, roleNames, "superadmin")
		assert.Contains(t, roleNames, "test_role")
		assert.Contains(t, roleNames, "permission_test")
	})
}

func TestRBACService_AssignSuperAdminRole_SystemRoles(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test users
	superAdmin := &models.User{
		DiscordUserID: "superadmin123",
		Username:      "superadmin",
	}
	err := db.Create(superAdmin).Error
	require.NoError(t, err)

	targetUser := &models.User{
		DiscordUserID: "target123",
		Username:      "targetuser",
	}
	err = db.Create(targetUser).Error
	require.NoError(t, err)

	// Create context with super admin user
	ctx := context.WithValue(context.Background(), "user_id", superAdmin.ID)

	// Test assigning super admin role
	t.Run("AssignSuperAdminRole", func(t *testing.T) {
		err = rbacService.AssignSuperAdminRole(ctx, targetUser.ID)
		require.NoError(t, err)

		// Verify role was assigned
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "superadmin", systemRoles[0].Name)

		// Verify user is now super admin
		isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.True(t, isSuperAdmin)
	})

	// Test removing super admin role
	t.Run("RemoveSuperAdminRole", func(t *testing.T) {
		err = rbacService.RemoveSuperAdminRole(ctx, targetUser.ID)
		require.NoError(t, err)

		// Verify role was removed
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 0)

		// Verify user is no longer super admin
		isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.False(t, isSuperAdmin)
	})

	// Test permission checks
	t.Run("PermissionChecks", func(t *testing.T) {
		// Non-super admin should not be able to assign super admin role
		regularUser := &models.User{
			DiscordUserID: "regular123",
			Username:      "regularuser",
		}
		err = db.Create(regularUser).Error
		require.NoError(t, err)

		regularCtx := context.WithValue(context.Background(), "user_id", regularUser.ID)
		err = rbacService.AssignSuperAdminRole(regularCtx, targetUser.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only super admins can assign super admin role")
	})
}

func TestRBACService_AssignInitialSuperAdminRole(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create test user
	user := &models.User{
		DiscordUserID: "test123",
		Username:      "testuser",
	}
	err := db.Create(user).Error
	require.NoError(t, err)

	t.Run("AssignInitialSuperAdminRole", func(t *testing.T) {
		err = rbacService.AssignInitialSuperAdminRole(context.Background(), user.ID)
		require.NoError(t, err)

		// Verify super admin system role was created
		var systemRole models.SystemRole
		err = db.Where("name = ?", "superadmin").First(&systemRole).Error
		require.NoError(t, err)
		assert.Equal(t, "superadmin", systemRole.Name)
		assert.Equal(t, "Super administrator with full system access", systemRole.Description)
		assert.Contains(t, systemRole.Permissions, models.PermissionSystemAdmin)

		// Verify role was assigned to user
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "superadmin", systemRoles[0].Name)

		// Verify user is super admin
		isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), user.ID)
		require.NoError(t, err)
		assert.True(t, isSuperAdmin)
	})

	t.Run("AssignInitialSuperAdminRole_SecondTime", func(t *testing.T) {
		// Create a new user for this test
		user := &models.User{
			DiscordUserID: "superadmin123",
			Username:      "superadmin",
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Should not fail when called again (idempotent)
		err = rbacService.AssignInitialSuperAdminRole(context.Background(), user.ID)
		require.NoError(t, err)

		// Verify user still has super admin role
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "superadmin", systemRoles[0].Name)
	})
} 

func TestRBACService_SuperAdminAutoAssignmentDuringSignup(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	t.Run("SuperAdminSignupCreatesRoleAutomatically", func(t *testing.T) {
		// Create mock services for auth flow
		mockDiscord := new(MockDiscordService)
		mockJWT := new(MockJWTService)
		mockRedis := new(MockRedisService)

		// Setup Discord service mocks for super admin user
		discordToken := &models.DiscordTokenResponse{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			ExpiresIn:    3600,
		}
		discordUser := &models.DiscordUser{
			ID:       "superadmin123", // Matches config
			Username: "superadmin",
			Avatar:   "avatar123",
			Email:    "superadmin@example.com",
		}

		mockDiscord.On("ExchangeCodeForToken", mock.Anything, "fake_code").Return(discordToken, nil)
		mockDiscord.On("GetUserInfo", mock.Anything, "access_token").Return(discordUser, nil)

		// Setup JWT service mocks
		expiresAt := time.Now().Add(time.Hour)
		mockJWT.On("GenerateAccessToken", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return("access_token", expiresAt, nil)
		mockJWT.On("GenerateRefreshToken", mock.AnythingOfType("*models.User"), mock.AnythingOfType("string")).Return("refresh_token", expiresAt, nil)

		// Setup Redis service mocks
		mockRedis.On("StoreSession", mock.Anything, mock.AnythingOfType("*models.Session")).Return(nil)

		// Create auth service with RBAC integration
		authService := NewAuthServiceWithRBAC(db, mockDiscord, mockJWT, mockRedis, rbacService)

		// Test the callback flow
		authResponse, err := authService.HandleCallback(context.Background(), "fake_code")
		require.NoError(t, err)
		require.NotNil(t, authResponse)

		// Verify user was created
		var user models.User
		err = db.Where("discord_user_id = ?", "superadmin123").First(&user).Error
		require.NoError(t, err)

		// Verify super admin gets BOTH superadmin and systemuser roles
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 2, "Super admin should have both superadmin and systemuser roles")
		
		// Check for both roles
		roleNames := make([]string, len(systemRoles))
		for i, role := range systemRoles {
			roleNames[i] = role.Name
		}
		assert.Contains(t, roleNames, "superadmin", "Super admin should have superadmin role")
		assert.Contains(t, roleNames, "systemuser", "Super admin should also have systemuser role")
		
		// Verify superadmin role has correct permissions
		for _, role := range systemRoles {
			if role.Name == "superadmin" {
				assert.Contains(t, role.Permissions, models.PermissionSystemAdmin)
			}
		}
	})
} 

func TestRBACService_SystemUserAutoAssignmentDuringSignup(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	t.Run("RegularUserSignupCreatesSystemUserRoleAutomatically", func(t *testing.T) {
		// Create a regular user (not super admin)
		user := &models.User{
			DiscordUserID: "regularuser123", // Not the super admin Discord ID
			Username:      "regularuser",
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Assign systemuser role to the new user
		err = rbacService.AssignDefaultSystemUserRole(context.Background(), user.ID)
		require.NoError(t, err)

		// Verify systemuser role was assigned
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "systemuser", systemRoles[0].Name)
		assert.Contains(t, systemRoles[0].Permissions, models.PermissionTemplateRead)
		// Verify tenant-scoped permissions are NOT included
		assert.NotContains(t, systemRoles[0].Permissions, models.PermissionServerRead)
		assert.NotContains(t, systemRoles[0].Permissions, models.PermissionLogRead)
	})

	t.Run("AssignDefaultSystemUserRole_SecondTime", func(t *testing.T) {
		// Create another regular user
		user := &models.User{
			DiscordUserID: "anotheruser123",
			Username:      "anotheruser",
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		// Should not fail when called again (idempotent)
		err = rbacService.AssignDefaultSystemUserRole(context.Background(), user.ID)
		require.NoError(t, err)

		// Verify user still has systemuser role
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), user.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "systemuser", systemRoles[0].Name)
	})

	t.Run("SystemUserRoleCreation", func(t *testing.T) {
		// Test that the systemuser role is created with correct permissions
		var systemUserRole models.SystemRole
		err := db.Where("name = ?", "systemuser").First(&systemUserRole).Error
		require.NoError(t, err)
		assert.Equal(t, "systemuser", systemUserRole.Name)
		assert.Equal(t, "Default system user with basic access", systemUserRole.Description)
		assert.Contains(t, systemUserRole.Permissions, models.PermissionTemplateRead)
		// Verify tenant-scoped permissions are NOT included
		assert.NotContains(t, systemUserRole.Permissions, models.PermissionServerRead)
		assert.NotContains(t, systemUserRole.Permissions, models.PermissionLogRead)
	})
}

func TestRBACService_ManualSuperAdminRoleAssignment(t *testing.T) {
	rbacService, db, cleanup := setupRBACTest(t)
	defer cleanup()

	// Create super admin user (by Discord ID) - this will be the one assigning roles
	superAdmin := &models.User{
		DiscordUserID: "superadmin123", // Matches config
		Username:      "superadmin",
	}
	err := db.Create(superAdmin).Error
	require.NoError(t, err)

	// Create target user to assign role to
	targetUser := &models.User{
		DiscordUserID: "target123",
		Username:      "targetuser",
	}
	err = db.Create(targetUser).Error
	require.NoError(t, err)

	// Create context with super admin user
	ctx := context.WithValue(context.Background(), "user_id", superAdmin.ID)

	t.Run("SuperAdminCanAssignSuperAdminRole", func(t *testing.T) {
		// Assign super admin role to target user
		err = rbacService.AssignSuperAdminRole(ctx, targetUser.ID)
		require.NoError(t, err)

		// Verify role was assigned
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 1)
		assert.Equal(t, "superadmin", systemRoles[0].Name)

		// Verify target user is now super admin
		isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.True(t, isSuperAdmin)
	})

	t.Run("SuperAdminCanRemoveSuperAdminRole", func(t *testing.T) {
		// Remove super admin role from target user
		err = rbacService.RemoveSuperAdminRole(ctx, targetUser.ID)
		require.NoError(t, err)

		// Verify role was removed
		systemRoles, err := rbacService.GetUserSystemRoles(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.Len(t, systemRoles, 0)

		// Verify target user is no longer super admin
		isSuperAdmin, err := rbacService.IsSuperAdmin(context.Background(), targetUser.ID)
		require.NoError(t, err)
		assert.False(t, isSuperAdmin)
	})

	t.Run("NonSuperAdminCannotAssignSuperAdminRole", func(t *testing.T) {
		// Create regular user
		regularUser := &models.User{
			DiscordUserID: "regular123",
			Username:      "regularuser",
		}
		err = db.Create(regularUser).Error
		require.NoError(t, err)

		// Try to assign super admin role as regular user
		regularCtx := context.WithValue(context.Background(), "user_id", regularUser.ID)
		err = rbacService.AssignSuperAdminRole(regularCtx, targetUser.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only super admins can assign super admin role")
	})

	t.Run("CannotRemoveSuperAdminRoleFromInitialSuperAdmin", func(t *testing.T) {
		// Try to remove super admin role from the initial super admin (by Discord ID)
		err = rbacService.RemoveSuperAdminRole(ctx, superAdmin.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot remove super admin role from the initial super admin")
	})
}
