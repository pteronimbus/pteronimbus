package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// TenantService handles tenant-related operations
type TenantService struct {
	db             *gorm.DB
	discordService DiscordServiceInterface
}

// NewTenantService creates a new tenant service
func NewTenantService(db *gorm.DB, discordService DiscordServiceInterface) *TenantService {
	return &TenantService{
		db:             db,
		discordService: discordService,
	}
}

// CreateTenant creates a new tenant from a Discord guild
func (ts *TenantService) CreateTenant(ctx context.Context, discordGuild *models.DiscordGuild, ownerID string) (*models.Tenant, error) {
	// Check if tenant already exists
	var existingTenant models.Tenant
	err := ts.db.Where("discord_server_id = ?", discordGuild.ID).First(&existingTenant).Error
	if err == nil {
		return nil, fmt.Errorf("tenant already exists for Discord server %s", discordGuild.ID)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing tenant: %w", err)
	}

	// Create new tenant
	tenant := &models.Tenant{
		DiscordServerID: discordGuild.ID,
		Name:            discordGuild.Name,
		Icon:            discordGuild.Icon,
		OwnerID:         ownerID,
		Config: models.TenantConfig{
			ResourceLimits: models.ResourceLimits{
				MaxGameServers: 5,
				MaxCPU:         "2",
				MaxMemory:      "4Gi",
				MaxStorage:     "10Gi",
			},
			Settings: make(map[string]string),
		},
	}

	err = ts.db.Create(tenant).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

// GetTenant retrieves a tenant by ID
func (ts *TenantService) GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := ts.db.Preload("Users").Preload("DiscordRoles").First(&tenant, "id = ?", tenantID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return &tenant, nil
}

// GetTenantByDiscordServerID retrieves a tenant by Discord server ID
func (ts *TenantService) GetTenantByDiscordServerID(ctx context.Context, discordServerID string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := ts.db.Preload("Users").Preload("DiscordRoles").First(&tenant, "discord_server_id = ?", discordServerID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return &tenant, nil
}

// GetUserTenants retrieves all tenants a user has access to
func (ts *TenantService) GetUserTenants(ctx context.Context, userID string) ([]models.Tenant, error) {
	var tenants []models.Tenant
	err := ts.db.Joins("JOIN user_tenants ON tenants.id = user_tenants.tenant_id").
		Where("user_tenants.user_id = ?", userID).
		Find(&tenants).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user tenants: %w", err)
	}

	return tenants, nil
}

// AddUserToTenant adds a user to a tenant with specified roles
func (ts *TenantService) AddUserToTenant(ctx context.Context, userID, tenantID string, roles []string, permissions []string) error {
	// Check if user-tenant relationship already exists
	var existingUserTenant models.UserTenant
	err := ts.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&existingUserTenant).Error
	if err == nil {
		// Update existing relationship
		existingUserTenant.Roles = models.StringArray(roles)
		existingUserTenant.Permissions = models.StringArray(permissions)
		existingUserTenant.UpdatedAt = time.Now()
		return ts.db.Save(&existingUserTenant).Error
	}
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check existing user-tenant relationship: %w", err)
	}

	// Create new user-tenant relationship
	userTenant := &models.UserTenant{
		UserID:      userID,
		TenantID:    tenantID,
		Roles:       models.StringArray(roles),
		Permissions: models.StringArray(permissions),
	}

	err = ts.db.Create(userTenant).Error
	if err != nil {
		return fmt.Errorf("failed to add user to tenant: %w", err)
	}

	return nil
}

// RemoveUserFromTenant removes a user from a tenant
func (ts *TenantService) RemoveUserFromTenant(ctx context.Context, userID, tenantID string) error {
	err := ts.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).Delete(&models.UserTenant{}).Error
	if err != nil {
		return fmt.Errorf("failed to remove user from tenant: %w", err)
	}

	return nil
}

// SyncDiscordRoles synchronizes Discord roles for a tenant
func (ts *TenantService) SyncDiscordRoles(ctx context.Context, tenantID, botToken string) error {
	// Get tenant
	tenant, err := ts.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Get Discord roles
	discordRoles, err := ts.discordService.GetGuildRoles(ctx, botToken, tenant.DiscordServerID)
	if err != nil {
		return fmt.Errorf("failed to get Discord roles: %w", err)
	}

	// Sync roles to database
	for _, discordRole := range discordRoles {
		var existingRole models.TenantDiscordRole
		err := ts.db.Where("tenant_id = ? AND discord_role_id = ?", tenantID, discordRole.ID).First(&existingRole).Error
		
		if err == gorm.ErrRecordNotFound {
			// Create new role
			newRole := &models.TenantDiscordRole{
				TenantID:      tenantID,
				DiscordRoleID: discordRole.ID,
				Name:          discordRole.Name,
				Color:         discordRole.Color,
				Position:      discordRole.Position,
				Permissions:   models.StringArray{}, // Will be mapped later
				Mentionable:   discordRole.Mentionable,
				Hoist:         discordRole.Hoist,
			}
			
			err = ts.db.Create(newRole).Error
			if err != nil {
				return fmt.Errorf("failed to create Discord role: %w", err)
			}
		} else if err == nil {
			// Update existing role
			existingRole.Name = discordRole.Name
			existingRole.Color = discordRole.Color
			existingRole.Position = discordRole.Position
			existingRole.Mentionable = discordRole.Mentionable
			existingRole.Hoist = discordRole.Hoist
			existingRole.UpdatedAt = time.Now()
			
			err = ts.db.Save(&existingRole).Error
			if err != nil {
				return fmt.Errorf("failed to update Discord role: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check existing Discord role: %w", err)
		}
	}

	return nil
}

// SyncDiscordUsers synchronizes Discord users for a tenant
func (ts *TenantService) SyncDiscordUsers(ctx context.Context, tenantID, botToken string) error {
	// Get tenant
	tenant, err := ts.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Get Discord members (limit to 1000 for now)
	discordMembers, err := ts.discordService.GetGuildMembers(ctx, botToken, tenant.DiscordServerID, 1000)
	if err != nil {
		return fmt.Errorf("failed to get Discord members: %w", err)
	}

	// Sync users to database
	for _, member := range discordMembers {
		if member.User == nil {
			continue
		}

		var existingUser models.TenantDiscordUser
		err := ts.db.Where("tenant_id = ? AND discord_user_id = ?", tenantID, member.User.ID).First(&existingUser).Error
		
		joinedAt, _ := time.Parse(time.RFC3339, member.JoinedAt)
		
		if err == gorm.ErrRecordNotFound {
			// Create new user
			newUser := &models.TenantDiscordUser{
				TenantID:      tenantID,
				DiscordUserID: member.User.ID,
				Username:      member.User.Username,
				DisplayName:   member.Nick,
				Avatar:        member.Avatar,
				Roles:         models.StringArray(member.Roles),
				JoinedAt:      &joinedAt,
				LastSyncAt:    time.Now(),
			}
			
			err = ts.db.Create(newUser).Error
			if err != nil {
				return fmt.Errorf("failed to create Discord user: %w", err)
			}
		} else if err == nil {
			// Update existing user
			existingUser.Username = member.User.Username
			existingUser.DisplayName = member.Nick
			existingUser.Avatar = member.Avatar
			existingUser.Roles = models.StringArray(member.Roles)
			existingUser.JoinedAt = &joinedAt
			existingUser.LastSyncAt = time.Now()
			existingUser.UpdatedAt = time.Now()
			
			err = ts.db.Save(&existingUser).Error
			if err != nil {
				return fmt.Errorf("failed to update Discord user: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check existing Discord user: %w", err)
		}
	}

	return nil
}

// UpdateTenantConfig updates tenant configuration
func (ts *TenantService) UpdateTenantConfig(ctx context.Context, tenantID string, config models.TenantConfig) error {
	err := ts.db.Model(&models.Tenant{}).Where("id = ?", tenantID).Update("config", config).Error
	if err != nil {
		return fmt.Errorf("failed to update tenant config: %w", err)
	}

	return nil
}

// DeleteTenant deletes a tenant and all associated data
func (ts *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	// Start transaction
	tx := ts.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete associated data
	if err := tx.Where("tenant_id = ?", tenantID).Delete(&models.UserTenant{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user-tenant relationships: %w", err)
	}

	if err := tx.Where("tenant_id = ?", tenantID).Delete(&models.TenantDiscordRole{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete Discord roles: %w", err)
	}

	if err := tx.Where("tenant_id = ?", tenantID).Delete(&models.TenantDiscordUser{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete Discord users: %w", err)
	}

	if err := tx.Where("tenant_id = ?", tenantID).Delete(&models.GameServer{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete game servers: %w", err)
	}

	// Delete tenant
	if err := tx.Delete(&models.Tenant{}, "id = ?", tenantID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	return tx.Commit().Error
}

// HasPermission checks if a user has a specific permission in a tenant
func (ts *TenantService) HasPermission(ctx context.Context, userID, tenantID, permission string) (bool, error) {
	var userTenant models.UserTenant
	err := ts.db.Where("user_id = ? AND tenant_id = ?", userID, tenantID).First(&userTenant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user-tenant relationship: %w", err)
	}

	// Check direct permissions
	for _, perm := range userTenant.Permissions {
		if perm == permission || perm == "*" {
			return true, nil
		}
	}

	// Check role-based permissions
	var discordRoles []models.TenantDiscordRole
	err = ts.db.Where("tenant_id = ? AND discord_role_id IN ?", tenantID, userTenant.Roles).Find(&discordRoles).Error
	if err != nil {
		return false, fmt.Errorf("failed to get Discord roles: %w", err)
	}

	for _, role := range discordRoles {
		for _, perm := range role.Permissions {
			if perm == permission || perm == "*" {
				return true, nil
			}
		}
	}

	return false, nil
}

// CheckManageServerPermission checks if a user has manage server permissions in Discord
func (ts *TenantService) CheckManageServerPermission(discordGuild *models.DiscordGuild) bool {
	if discordGuild.Owner {
		return true
	}

	// Parse permissions string to int64
	permissions, err := strconv.ParseInt(discordGuild.Permissions, 10, 64)
	if err != nil {
		return false
	}

	// Check for MANAGE_GUILD permission (0x00000020)
	const MANAGE_GUILD = 0x00000020
	return (permissions & MANAGE_GUILD) == MANAGE_GUILD
}