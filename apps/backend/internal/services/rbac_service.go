package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"gorm.io/gorm"
)

// RBACService handles role-based access control operations
type RBACService struct {
	db     *gorm.DB
	config *config.RBACConfig
}

// NewRBACService creates a new RBAC service
func NewRBACService(db *gorm.DB, config *config.RBACConfig) *RBACService {
	return &RBACService{
		db:     db,
		config: config,
	}
}

// HasPermission checks if a user has a specific permission in a tenant
func (rs *RBACService) HasPermission(ctx context.Context, userID, tenantID, permission string) (bool, error) {
	// First check if user is super admin
	isSuperAdmin, err := rs.IsSuperAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check super admin status: %w", err)
	}
	if isSuperAdmin {
		return true, nil
	}

	// Get user-tenant relationship
	var userTenant models.UserTenant
	err = rs.db.WithContext(ctx).Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&userTenant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user-tenant relationship: %w", err)
	}

	// Check direct permissions
	for _, perm := range userTenant.Permissions {
		if perm == permission || perm == models.PermissionAdminAll {
			return true, nil
		}
	}

	// Check role-based permissions
	if len(userTenant.Roles) > 0 {
		hasPermission, err := rs.checkRolePermissions(ctx, tenantID, userTenant.Roles, permission)
		if err != nil {
			return false, fmt.Errorf("failed to check role permissions: %w", err)
		}
		if hasPermission {
			return true, nil
		}
	}

	return false, nil
}

// checkRolePermissions checks if any of the user's roles have the required permission
func (rs *RBACService) checkRolePermissions(ctx context.Context, tenantID string, roleIDs []string, permission string) (bool, error) {
	// Check Discord roles first
	var discordRoles []models.TenantDiscordRole
	err := rs.db.WithContext(ctx).Where("tenant_id = ? AND discord_role_id IN ?", tenantID, roleIDs).Find(&discordRoles).Error
	if err != nil {
		return false, fmt.Errorf("failed to get Discord roles: %w", err)
	}

	for _, role := range discordRoles {
		for _, perm := range role.Permissions {
			if perm == permission || perm == models.PermissionAdminAll {
				return true, nil
			}
		}
	}

	// Check internal roles
	var roles []models.Role
	err = rs.db.WithContext(ctx).Where("tenant_id = ? AND name IN ?", tenantID, roleIDs).Find(&roles).Error
	if err != nil {
		return false, fmt.Errorf("failed to get internal roles: %w", err)
	}

	for _, role := range roles {
		for _, perm := range role.Permissions {
			if perm == permission || perm == models.PermissionAdminAll {
				return true, nil
			}
		}
	}

	return false, nil
}

// IsSuperAdmin checks if a user is a super admin
func (rs *RBACService) IsSuperAdmin(ctx context.Context, userID string) (bool, error) {
	// Get user details
	var user models.User
	err := rs.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user's Discord ID matches super admin Discord ID
	// This is the primary way to identify super admins (system-wide, not tenant-scoped)
	if rs.config.SuperAdminDiscordID != "" && strings.EqualFold(user.DiscordUserID, rs.config.SuperAdminDiscordID) {
		return true, nil
	}

	// Ensure superadmin system role exists before checking
	err = rs.ensureSuperAdminSystemRole(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to ensure super admin system role: %w", err)
	}

	// Check if user has system admin role
	hasSystemAdmin, err := rs.HasSystemPermission(ctx, userID, models.PermissionSystemAdmin)
	if err != nil {
		return false, fmt.Errorf("failed to check system admin permission: %w", err)
	}
	if hasSystemAdmin {
		return true, nil
	}

	return false, nil
}

// HasSystemPermission checks if a user has a specific system permission
func (rs *RBACService) HasSystemPermission(ctx context.Context, userID, permission string) (bool, error) {
	// Get user's system roles
	var userSystemRoles []models.UserSystemRole
	err := rs.db.WithContext(ctx).
		Preload("SystemRole").
		Where("user_id = ?", userID).
		Find(&userSystemRoles).Error
	if err != nil {
		return false, fmt.Errorf("failed to get user system roles: %w", err)
	}

	// Check permissions from system roles
	for _, userSystemRole := range userSystemRoles {
		for _, perm := range userSystemRole.SystemRole.Permissions {
			if perm == permission || perm == models.PermissionAdminAll {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetUserSystemRoles returns all system roles for a user
func (rs *RBACService) GetUserSystemRoles(ctx context.Context, userID string) ([]models.SystemRole, error) {
	var userSystemRoles []models.UserSystemRole
	err := rs.db.WithContext(ctx).
		Preload("SystemRole").
		Where("user_id = ?", userID).
		Find(&userSystemRoles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user system roles: %w", err)
	}

	systemRoles := make([]models.SystemRole, len(userSystemRoles))
	for i, userSystemRole := range userSystemRoles {
		systemRoles[i] = userSystemRole.SystemRole
	}

	return systemRoles, nil
}

// AssignSystemRoleToUser assigns a system role to a user
func (rs *RBACService) AssignSystemRoleToUser(ctx context.Context, userID, systemRoleName string) error {
	// Get system role
	var systemRole models.SystemRole
	err := rs.db.WithContext(ctx).Where("name = ?", systemRoleName).First(&systemRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("system role '%s' not found", systemRoleName)
		}
		return fmt.Errorf("failed to get system role: %w", err)
	}

	// Check if user already has this system role
	var existingUserSystemRole models.UserSystemRole
	err = rs.db.WithContext(ctx).
		Where("user_id = ? AND system_role_id = ?", userID, systemRole.ID).
		First(&existingUserSystemRole).Error
	
	if err == nil {
		// User already has this role
		return nil
	}
	
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing user system role: %w", err)
	}

	// Assign system role to user
	userSystemRole := &models.UserSystemRole{
		UserID:       userID,
		SystemRoleID: systemRole.ID,
	}

	err = rs.db.WithContext(ctx).Create(userSystemRole).Error
	if err != nil {
		return fmt.Errorf("failed to assign system role to user: %w", err)
	}

	return nil
}

// RemoveSystemRoleFromUser removes a system role from a user
func (rs *RBACService) RemoveSystemRoleFromUser(ctx context.Context, userID, systemRoleName string) error {
	// Get system role
	var systemRole models.SystemRole
	err := rs.db.WithContext(ctx).Where("name = ?", systemRoleName).First(&systemRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("system role '%s' not found", systemRoleName)
		}
		return fmt.Errorf("failed to get system role: %w", err)
	}

	// Remove system role from user
	err = rs.db.WithContext(ctx).
		Where("user_id = ? AND system_role_id = ?", userID, systemRole.ID).
		Delete(&models.UserSystemRole{}).Error
	if err != nil {
		return fmt.Errorf("failed to remove system role from user: %w", err)
	}

	return nil
}

// CreateSystemRole creates a new system role
func (rs *RBACService) CreateSystemRole(ctx context.Context, name, description string, permissions []string) (*models.SystemRole, error) {
	systemRole := &models.SystemRole{
		Name:        name,
		Description: description,
		Permissions: models.StringArray(permissions),
	}

	err := rs.db.WithContext(ctx).Create(systemRole).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create system role: %w", err)
	}

	return systemRole, nil
}

// GetSystemRoles returns all system roles
func (rs *RBACService) GetSystemRoles(ctx context.Context) ([]models.SystemRole, error) {
	var systemRoles []models.SystemRole
	err := rs.db.WithContext(ctx).Find(&systemRoles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get system roles: %w", err)
	}

	return systemRoles, nil
}

// CreateRole creates a new role in a tenant
func (rs *RBACService) CreateRole(ctx context.Context, tenantID, name string, permissions []string, isSystemRole bool) (*models.Role, error) {
	role := &models.Role{
		TenantID:     tenantID,
		Name:         name,
		Permissions:  models.StringArray(permissions),
		IsSystemRole: isSystemRole,
	}

	err := rs.db.WithContext(ctx).Create(role).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// UpdateRole updates an existing role
func (rs *RBACService) UpdateRole(ctx context.Context, roleID string, name string, permissions []string) (*models.Role, error) {
	var role models.Role
	err := rs.db.WithContext(ctx).Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	role.Name = name
	role.Permissions = models.StringArray(permissions)
	role.UpdatedAt = time.Now()

	err = rs.db.WithContext(ctx).Save(&role).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &role, nil
}

// DeleteRole deletes a role
func (rs *RBACService) DeleteRole(ctx context.Context, roleID string) error {
	var role models.Role
	err := rs.db.WithContext(ctx).Where("id = ?", roleID).First(&role).Error
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	if role.IsSystemRole {
		return fmt.Errorf("cannot delete system role")
	}

	err = rs.db.WithContext(ctx).Delete(&role).Error
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// GetRoles returns all roles for a tenant
func (rs *RBACService) GetRoles(ctx context.Context, tenantID string) ([]models.Role, error) {
	var roles []models.Role
	err := rs.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return roles, nil
}

// AssignRoleToUser assigns a role to a user in a tenant
func (rs *RBACService) AssignRoleToUser(ctx context.Context, userID, tenantID, roleName string) error {
	// Get or create user-tenant relationship
	var userTenant models.UserTenant
	err := rs.db.WithContext(ctx).Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&userTenant).Error
	if err == gorm.ErrRecordNotFound {
		userTenant = models.UserTenant{
			UserID:      userID,
			TenantID:    tenantID,
			Roles:       models.StringArray{},
			Permissions: models.StringArray{},
		}
	} else if err != nil {
		return fmt.Errorf("failed to get user-tenant relationship: %w", err)
	}

	// Add role if not already present
	roleExists := false
	for _, role := range userTenant.Roles {
		if role == roleName {
			roleExists = true
			break
		}
	}

	if !roleExists {
		userTenant.Roles = append(userTenant.Roles, roleName)
		userTenant.UpdatedAt = time.Now()

		if userTenant.ID == "" {
			err = rs.db.WithContext(ctx).Create(&userTenant).Error
		} else {
			err = rs.db.WithContext(ctx).Save(&userTenant).Error
		}

		if err != nil {
			return fmt.Errorf("failed to assign role to user: %w", err)
		}
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user in a tenant
func (rs *RBACService) RemoveRoleFromUser(ctx context.Context, userID, tenantID, roleName string) error {
	var userTenant models.UserTenant
	err := rs.db.WithContext(ctx).Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&userTenant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // User doesn't have any roles in this tenant
		}
		return fmt.Errorf("failed to get user-tenant relationship: %w", err)
	}

	// Remove role if present
	var newRoles []string
	for _, role := range userTenant.Roles {
		if role != roleName {
			newRoles = append(newRoles, role)
		}
	}

	userTenant.Roles = models.StringArray(newRoles)
	userTenant.UpdatedAt = time.Now()

	err = rs.db.WithContext(ctx).Save(&userTenant).Error
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	return nil
}

// LogPermissionChange logs a permission change for audit purposes
func (rs *RBACService) LogPermissionChange(ctx context.Context, userID, tenantID, action, resourceType, resourceID, oldValue, newValue, reason, performedBy string) error {
	auditLog := &models.PermissionAuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		OldValue:     oldValue,
		NewValue:     newValue,
		Reason:       reason,
		PerformedBy:  performedBy,
	}

	// Only set TenantID if it's not empty (for system-level operations)
	if tenantID != "" {
		auditLog.TenantID = &tenantID
	}

	err := rs.db.WithContext(ctx).Create(auditLog).Error
	if err != nil {
		return fmt.Errorf("failed to log permission change: %w", err)
	}

	return nil
}

// GetUserPermissions returns all permissions for a user in a tenant
func (rs *RBACService) GetUserPermissions(ctx context.Context, userID, tenantID string) ([]string, error) {
	// Check if user is super admin
	isSuperAdmin, err := rs.IsSuperAdmin(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check super admin status: %w", err)
	}
	if isSuperAdmin {
		return []string{models.PermissionAdminAll}, nil
	}

	var userTenant models.UserTenant
	err = rs.db.WithContext(ctx).Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&userTenant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get user-tenant relationship: %w", err)
	}

	permissions := make([]string, 0)
	permissionSet := make(map[string]bool)

	// Add direct permissions
	for _, perm := range userTenant.Permissions {
		if !permissionSet[perm] {
			permissions = append(permissions, perm)
			permissionSet[perm] = true
		}
	}

	// Add role-based permissions
	if len(userTenant.Roles) > 0 {
		rolePermissions, err := rs.getRolePermissions(ctx, tenantID, userTenant.Roles)
		if err != nil {
			return nil, fmt.Errorf("failed to get role permissions: %w", err)
		}

		for _, perm := range rolePermissions {
			if !permissionSet[perm] {
				permissions = append(permissions, perm)
				permissionSet[perm] = true
			}
		}
	}

	return permissions, nil
}

// getRolePermissions gets all permissions from the specified roles
func (rs *RBACService) getRolePermissions(ctx context.Context, tenantID string, roleIDs []string) ([]string, error) {
	permissions := make([]string, 0)

	// Get Discord role permissions
	var discordRoles []models.TenantDiscordRole
	err := rs.db.WithContext(ctx).Where("tenant_id = ? AND discord_role_id IN ?", tenantID, roleIDs).Find(&discordRoles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get Discord roles: %w", err)
	}

	for _, role := range discordRoles {
		permissions = append(permissions, []string(role.Permissions)...)
	}

	// Get internal role permissions
	var roles []models.Role
	err = rs.db.WithContext(ctx).Where("tenant_id = ? AND name IN ?", tenantID, roleIDs).Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get internal roles: %w", err)
	}

	for _, role := range roles {
		permissions = append(permissions, []string(role.Permissions)...)
	}

	return permissions, nil
} 

// AssignSuperAdminRole assigns the super admin role to a user
func (rs *RBACService) AssignSuperAdminRole(ctx context.Context, userID string) error {
	// Check if the performing user is a super admin
	performingUserID, ok := ctx.Value("user_id").(string)
	if !ok {
		return fmt.Errorf("user context not found")
	}

	isSuperAdmin, err := rs.IsSuperAdmin(ctx, performingUserID)
	if err != nil {
		return fmt.Errorf("failed to check super admin status: %w", err)
	}
	if !isSuperAdmin {
		return fmt.Errorf("only super admins can assign super admin role")
	}

	// Prevent assigning super admin role to the initial super admin (by Discord ID)
	var user models.User
	err = rs.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if rs.config.SuperAdminDiscordID != "" && strings.EqualFold(user.DiscordUserID, rs.config.SuperAdminDiscordID) {
		return fmt.Errorf("cannot assign super admin role to the initial super admin")
	}

	// Ensure superadmin system role exists
	err = rs.ensureSuperAdminSystemRole(ctx)
	if err != nil {
		return fmt.Errorf("failed to ensure super admin system role: %w", err)
	}

	// Assign the super admin system role
	err = rs.AssignSystemRoleToUser(ctx, userID, "superadmin")
	if err != nil {
		return fmt.Errorf("failed to assign super admin system role: %w", err)
	}

	// Log the permission change
	err = rs.LogPermissionChange(ctx, userID, "", "assign", "system_role", "superadmin", "", "superadmin", "Super admin role assigned", performingUserID)
	if err != nil {
		return fmt.Errorf("failed to log permission change: %w", err)
	}

	return nil
}

// AssignInitialSuperAdminRole assigns the super admin role to a user during initial setup
// This bypasses the normal permission checks since it's used during user creation
func (rs *RBACService) AssignInitialSuperAdminRole(ctx context.Context, userID string) error {
	// Create the super admin system role if it doesn't exist
	var superAdminRole models.SystemRole
	err := rs.db.WithContext(ctx).Where("name = ?", "superadmin").First(&superAdminRole).Error
	if err == gorm.ErrRecordNotFound {
		// Create the super admin system role
		superAdminRole = models.SystemRole{
			Name:        "superadmin",
			Description: "Super administrator with full system access",
			Permissions: models.StringArray{models.PermissionSystemAdmin},
		}
		err = rs.db.WithContext(ctx).Create(&superAdminRole).Error
		if err != nil {
			return fmt.Errorf("failed to create super admin system role: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check super admin system role: %w", err)
	}

	// Assign the super admin system role to the user
	err = rs.AssignSystemRoleToUser(ctx, userID, "superadmin")
	if err != nil {
		return fmt.Errorf("failed to assign super admin system role: %w", err)
	}

	return nil
}

// RemoveSuperAdminRole removes the super admin role from a user
func (rs *RBACService) RemoveSuperAdminRole(ctx context.Context, userID string) error {
	// Check if the performing user is a super admin
	performingUserID, ok := ctx.Value("user_id").(string)
	if !ok {
		return fmt.Errorf("user context not found")
	}

	isSuperAdmin, err := rs.IsSuperAdmin(ctx, performingUserID)
	if err != nil {
		return fmt.Errorf("failed to check super admin status: %w", err)
	}
	if !isSuperAdmin {
		return fmt.Errorf("only super admins can remove super admin role")
	}

	// Prevent removing super admin role from the initial super admin
	var user models.User
	err = rs.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if rs.config.SuperAdminDiscordID != "" && strings.EqualFold(user.DiscordUserID, rs.config.SuperAdminDiscordID) {
		return fmt.Errorf("cannot remove super admin role from the initial super admin")
	}

	// Ensure superadmin system role exists
	err = rs.ensureSuperAdminSystemRole(ctx)
	if err != nil {
		return fmt.Errorf("failed to ensure super admin system role: %w", err)
	}

	// Remove the super admin system role
	err = rs.RemoveSystemRoleFromUser(ctx, userID, "superadmin")
	if err != nil {
		return fmt.Errorf("failed to remove super admin system role: %w", err)
	}

	// Log the permission change
	err = rs.LogPermissionChange(ctx, userID, "", "remove", "system_role", "superadmin", "superadmin", "", "Super admin role removed", performingUserID)
	if err != nil {
		return fmt.Errorf("failed to log permission change: %w", err)
	}

	return nil
} 

// ensureSuperAdminSystemRole ensures the superadmin system role exists
func (rs *RBACService) ensureSuperAdminSystemRole(ctx context.Context) error {
	var superAdminRole models.SystemRole
	err := rs.db.WithContext(ctx).Where("name = ?", "superadmin").First(&superAdminRole).Error
	if err == gorm.ErrRecordNotFound {
		// Create the super admin system role
		superAdminRole = models.SystemRole{
			Name:        "superadmin",
			Description: "Super administrator with full system access",
			Permissions: models.StringArray{models.PermissionSystemAdmin},
		}
		err = rs.db.WithContext(ctx).Create(&superAdminRole).Error
		if err != nil {
			return fmt.Errorf("failed to create super admin system role: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check super admin system role: %w", err)
	}

	return nil
} 