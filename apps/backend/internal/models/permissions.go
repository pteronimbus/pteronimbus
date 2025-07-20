package models

// Permission constants for different resources and actions
const (
	// Server permissions
	PermissionServerCreate = "server:create"
	PermissionServerRead   = "server:read"
	PermissionServerWrite  = "server:write"
	PermissionServerDelete = "server:delete"
	PermissionServerStart  = "server:start"
	PermissionServerStop   = "server:stop"
	PermissionServerRestart = "server:restart"

	// Console permissions
	PermissionConsoleRead    = "console:read"
	PermissionConsoleExecute = "console:execute"

	// File permissions
	PermissionFileRead   = "file:read"
	PermissionFileWrite  = "file:write"
	PermissionFileDelete = "file:delete"

	// Backup permissions
	PermissionBackupCreate = "backup:create"
	PermissionBackupRead   = "backup:read"
	PermissionBackupDelete = "backup:delete"
	PermissionBackupRestore = "backup:restore"

	// Log permissions
	PermissionLogRead = "log:read"

	// Template permissions
	PermissionTemplateCreate = "template:create"
	PermissionTemplateRead   = "template:read"
	PermissionTemplateWrite  = "template:write"
	PermissionTemplateDelete = "template:delete"

	// User permissions
	PermissionUserCreate = "user:create"
	PermissionUserRead   = "user:read"
	PermissionUserWrite  = "user:write"
	PermissionUserDelete = "user:delete"

	// Role permissions
	PermissionRoleCreate = "role:create"
	PermissionRoleRead   = "role:read"
	PermissionRoleWrite  = "role:write"
	PermissionRoleDelete = "role:delete"

	// Admin permissions
	PermissionAdminAll = "*"
	PermissionSuperAdmin = "superadmin"
)

// PermissionScope defines the scope of a permission
type PermissionScope string

const (
	ScopeOwn    PermissionScope = "own"
	ScopeTenant PermissionScope = "tenant"
	ScopeGlobal PermissionScope = "global"
)

// PermissionDefinition represents a permission with resource, action, and scope
type PermissionDefinition struct {
	Resource string          `json:"resource"`
	Action   string          `json:"action"`
	Scope    PermissionScope `json:"scope"`
}

// NewPermissionDefinition creates a new permission definition
func NewPermissionDefinition(resource, action string, scope PermissionScope) PermissionDefinition {
	return PermissionDefinition{
		Resource: resource,
		Action:   action,
		Scope:    scope,
	}
}

// String returns the permission as a string in format "resource:action"
func (p PermissionDefinition) String() string {
	return p.Resource + ":" + p.Action
}

// Matches checks if this permission matches another permission
func (p PermissionDefinition) Matches(other PermissionDefinition) bool {
	// Check if either permission is a wildcard
	if p.Resource == "*" || other.Resource == "*" {
		return true
	}
	if p.Action == "*" || other.Action == "*" {
		return p.Resource == other.Resource
	}
	
	return p.Resource == other.Resource && p.Action == other.Action
}

// HasPermission checks if a user has a specific permission
func (u *User) HasPermission(permission string, tenantID string) bool {
	// This method would be implemented to check against the user's permissions
	// For now, we'll return false as this should be implemented in the service layer
	return false
}

// GetDefaultPermissions returns default permissions for new users
func GetDefaultPermissions() []string {
	return []string{
		PermissionServerRead,
		PermissionLogRead,
	}
}

// GetAdminPermissions returns all admin permissions
func GetAdminPermissions() []string {
	return []string{
		PermissionAdminAll,
	}
}

// GetServerManagerPermissions returns permissions for server managers
func GetServerManagerPermissions() []string {
	return []string{
		PermissionServerCreate,
		PermissionServerRead,
		PermissionServerWrite,
		PermissionServerDelete,
		PermissionServerStart,
		PermissionServerStop,
		PermissionServerRestart,
		PermissionConsoleRead,
		PermissionConsoleExecute,
		PermissionFileRead,
		PermissionFileWrite,
		PermissionBackupCreate,
		PermissionBackupRead,
		PermissionBackupRestore,
		PermissionLogRead,
	}
}

// GetModeratorPermissions returns permissions for moderators
func GetModeratorPermissions() []string {
	return []string{
		PermissionServerRead,
		PermissionConsoleRead,
		PermissionFileRead,
		PermissionLogRead,
		PermissionUserRead,
	}
}

// GetUserPermissions returns basic user permissions
func GetUserPermissions() []string {
	return []string{
		PermissionServerRead,
		PermissionLogRead,
	}
} 