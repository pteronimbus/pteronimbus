package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// StringArray is a custom type for handling PostgreSQL text arrays
type StringArray []string

// Scan implements the sql.Scanner interface for reading from database
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = StringArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return errors.New("cannot scan into StringArray")
	}
}

// Value implements the driver.Valuer interface for writing to database
func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "{}", nil
	}
	
	// Convert to PostgreSQL array format
	result := "{"
	for i, str := range s {
		if i > 0 {
			result += ","
		}
		// Escape quotes and wrap in quotes
		escaped := `"` + str + `"`
		result += escaped
	}
	result += "}"
	
	return result, nil
}

// Tenant represents a Discord server with Pteronimbus installed
type Tenant struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DiscordServerID string         `json:"discord_server_id" gorm:"uniqueIndex;not null"`
	Name            string         `json:"name" gorm:"not null"`
	Icon            string         `json:"icon"`
	OwnerID         string         `json:"owner_id" gorm:"not null"`
	Config          TenantConfig   `json:"config" gorm:"type:jsonb"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Users        []UserTenant          `json:"users,omitempty" gorm:"foreignKey:TenantID"`
	DiscordRoles []TenantDiscordRole   `json:"discord_roles,omitempty" gorm:"foreignKey:TenantID"`
	DiscordUsers []TenantDiscordUser   `json:"discord_users,omitempty" gorm:"foreignKey:TenantID"`
	GameServers  []GameServer          `json:"game_servers,omitempty" gorm:"foreignKey:TenantID"`
}

// TenantConfig holds tenant-specific configuration
type TenantConfig struct {
	DefaultGameTemplate  string            `json:"default_game_template,omitempty"`
	ResourceLimits       ResourceLimits    `json:"resource_limits,omitempty"`
	NotificationChannels []string          `json:"notification_channels,omitempty"`
	Settings             map[string]string `json:"settings,omitempty"`
}

// Scan implements the sql.Scanner interface for reading from database
func (tc *TenantConfig) Scan(value interface{}) error {
	if value == nil {
		*tc = TenantConfig{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, tc)
	case string:
		return json.Unmarshal([]byte(v), tc)
	default:
		return errors.New("cannot scan into TenantConfig")
	}
}

// Value implements the driver.Valuer interface for writing to database
func (tc TenantConfig) Value() (driver.Value, error) {
	return json.Marshal(tc)
}

// ResourceLimits defines resource constraints for a tenant
type ResourceLimits struct {
	MaxGameServers int    `json:"max_game_servers"`
	MaxCPU         string `json:"max_cpu"`
	MaxMemory      string `json:"max_memory"`
	MaxStorage     string `json:"max_storage"`
}

// UserTenant represents the relationship between users and tenants
type UserTenant struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      string         `json:"user_id" gorm:"not null;index"`
	TenantID    string         `json:"tenant_id" gorm:"not null;index"`
	Roles       StringArray    `json:"roles" gorm:"type:text[]"`
	Permissions StringArray    `json:"permissions" gorm:"type:text[]"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// TenantDiscordRole represents a Discord role within a tenant
type TenantDiscordRole struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID        string         `json:"tenant_id" gorm:"not null;index"`
	DiscordRoleID   string         `json:"discord_role_id" gorm:"not null"`
	Name            string         `json:"name" gorm:"not null"`
	Color           int            `json:"color"`
	Position        int            `json:"position"`
	Permissions     StringArray    `json:"permissions" gorm:"type:text[]"`
	Mentionable     bool           `json:"mentionable"`
	Hoist           bool           `json:"hoist"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// TenantDiscordUser represents a Discord user within a tenant context
type TenantDiscordUser struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID        string         `json:"tenant_id" gorm:"not null;index"`
	DiscordUserID   string         `json:"discord_user_id" gorm:"not null"`
	Username        string         `json:"username" gorm:"not null"`
	DisplayName     string         `json:"display_name"`
	Avatar          string         `json:"avatar"`
	Roles           StringArray    `json:"roles" gorm:"type:text[]"`
	JoinedAt        *time.Time     `json:"joined_at"`
	LastSyncAt      time.Time      `json:"last_sync_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// GameServer represents a game server instance
type GameServer struct {
	ID         string            `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TenantID   string            `json:"tenant_id" gorm:"not null;index"`
	TemplateID string            `json:"template_id"`
	Name       string            `json:"name" gorm:"not null"`
	GameType   string            `json:"game_type" gorm:"not null"`
	Config     GameServerConfig  `json:"config" gorm:"type:jsonb"`
	Status     GameServerStatus  `json:"status" gorm:"type:jsonb"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `json:"-" gorm:"index"`

	// Relationships
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// GameServerConfig holds game server configuration
type GameServerConfig struct {
	Image           string            `json:"image"`
	Ports           []Port            `json:"ports"`
	Environment     map[string]string `json:"environment"`
	Resources       ResourceRequirements `json:"resources"`
	PersistentData  []VolumeMount     `json:"persistent_data"`
	StartupCommand  []string          `json:"startup_command"`
}

// Scan implements the sql.Scanner interface for reading from database
func (gsc *GameServerConfig) Scan(value interface{}) error {
	if value == nil {
		*gsc = GameServerConfig{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, gsc)
	case string:
		return json.Unmarshal([]byte(v), gsc)
	default:
		return errors.New("cannot scan into GameServerConfig")
	}
}

// Value implements the driver.Valuer interface for writing to database
func (gsc GameServerConfig) Value() (driver.Value, error) {
	return json.Marshal(gsc)
}

// GameServerStatus represents the current status of a game server
type GameServerStatus struct {
	Phase       string    `json:"phase"` // Pending, Running, Stopped, Failed
	Message     string    `json:"message"`
	LastUpdated time.Time `json:"last_updated"`
	PlayerCount int       `json:"player_count"`
	Uptime      string    `json:"uptime"`
}

// Scan implements the sql.Scanner interface for reading from database
func (gss *GameServerStatus) Scan(value interface{}) error {
	if value == nil {
		*gss = GameServerStatus{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, gss)
	case string:
		return json.Unmarshal([]byte(v), gss)
	default:
		return errors.New("cannot scan into GameServerStatus")
	}
}

// Value implements the driver.Valuer interface for writing to database
func (gss GameServerStatus) Value() (driver.Value, error) {
	return json.Marshal(gss)
}

// Port represents a network port configuration
type Port struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// ResourceRequirements defines resource requirements for a game server
type ResourceRequirements struct {
	Requests ResourceList `json:"requests"`
	Limits   ResourceList `json:"limits"`
}

// ResourceList defines CPU and memory resources
type ResourceList struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// VolumeMount represents a volume mount configuration
type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mount_path"`
	Size      string `json:"size"`
}

// TableName returns the table name for Tenant
func (Tenant) TableName() string {
	return "tenants"
}

// TableName returns the table name for UserTenant
func (UserTenant) TableName() string {
	return "user_tenants"
}

// TableName returns the table name for TenantDiscordRole
func (TenantDiscordRole) TableName() string {
	return "discord_roles"
}

// TableName returns the table name for TenantDiscordUser
func (TenantDiscordUser) TableName() string {
	return "discord_users"
}

// Activity represents an activity log entry
type Activity struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Timestamp string    `json:"timestamp"`
}

// DiscordStats represents Discord statistics for a tenant
type DiscordStats struct {
	MemberCount int    `json:"memberCount"`
	RoleCount   int    `json:"roleCount"`
	LastSync    string `json:"lastSync"`
}

// TableName returns the table name for GameServer
func (GameServer) TableName() string {
	return "game_servers"
}