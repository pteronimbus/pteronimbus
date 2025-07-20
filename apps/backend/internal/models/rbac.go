package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission represents a permission definition
type Permission struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Resource    string         `json:"resource" gorm:"not null"`
	Action      string         `json:"action" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Role represents a role definition within a tenant
type Role struct {
	ID           string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID     string         `json:"tenant_id" gorm:"not null;index"`
	Name         string         `json:"name" gorm:"not null"`
	Permissions  StringArray    `json:"permissions" gorm:"type:text[]"`
	IsSystemRole bool           `json:"is_system_role" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// SystemRole represents a system-wide role (not tenant-scoped)
type SystemRole struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Permissions StringArray    `json:"permissions" gorm:"type:text[]"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserSystemRole represents the relationship between users and system roles
type UserSystemRole struct {
	ID         string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string         `json:"user_id" gorm:"not null;index"`
	SystemRoleID string       `json:"system_role_id" gorm:"not null;index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User       User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SystemRole SystemRole `json:"system_role,omitempty" gorm:"foreignKey:SystemRoleID"`
}

// PermissionAuditLog represents audit logging for permission changes
type PermissionAuditLog struct {
	ID           string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID       string         `json:"user_id" gorm:"index"`
	TenantID     *string        `json:"tenant_id" gorm:"index"` // Nullable for system-level operations
	Action       string         `json:"action" gorm:"not null"`
	ResourceType string         `json:"resource_type" gorm:"not null"`
	ResourceID   string         `json:"resource_id"`
	OldValue     string         `json:"old_value"`
	NewValue     string         `json:"new_value"`
	Reason       string         `json:"reason"`
	PerformedBy  string         `json:"performed_by" gorm:"index"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User       User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Tenant     Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	Performer  User   `json:"performer,omitempty" gorm:"foreignKey:PerformedBy"`
}

// GuildMembershipCache represents cached Discord guild membership data
type GuildMembershipCache struct {
	ID         string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string         `json:"user_id" gorm:"not null;index"`
	GuildID    string         `json:"guild_id" gorm:"not null"`
	Roles      StringArray    `json:"roles" gorm:"type:text[]"`
	Permissions int64         `json:"permissions" gorm:"not null;default:0"`
	LastSync   time.Time      `json:"last_sync"`
	ExpiresAt  time.Time      `json:"expires_at" gorm:"not null;index"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}

// TableName returns the table name for SystemRole
func (SystemRole) TableName() string {
	return "system_roles"
}

// TableName returns the table name for UserSystemRole
func (UserSystemRole) TableName() string {
	return "user_system_roles"
}

// TableName returns the table name for PermissionAuditLog
func (PermissionAuditLog) TableName() string {
	return "permission_audit_log"
}

// TableName returns the table name for GuildMembershipCache
func (GuildMembershipCache) TableName() string {
	return "guild_membership_cache"
} 