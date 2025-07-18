package services

import (
	"context"
	"time"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"gorm.io/gorm"
)

// GameServerService implements GameServerServiceInterface
type GameServerService struct {
	db *gorm.DB
}

// NewGameServerService creates a new game server service
func NewGameServerService(db *gorm.DB) GameServerServiceInterface {
	return &GameServerService{
		db: db,
	}
}

// GetTenantServers retrieves all game servers for a tenant
func (gss *GameServerService) GetTenantServers(ctx context.Context, tenantID string) ([]models.GameServer, error) {
	var servers []models.GameServer

	// For now, return mock data until we implement the full database integration
	// In a real implementation, this would query the database:
	// err := gss.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&servers).Error

	servers = []models.GameServer{
		{
			ID:         "server-1",
			TenantID:   tenantID,
			TemplateID: "template-minecraft",
			Name:       "Survival World",
			GameType:   "minecraft",
			Config: models.GameServerConfig{
				Environment: map[string]string{
					"EULA":        "TRUE",
					"DIFFICULTY":  "normal",
					"MAX_PLAYERS": "20",
				},
			},
			Status: models.GameServerStatus{
				Phase:       "Running",
				Message:     "Server is running normally",
				LastUpdated: time.Now(),
				PlayerCount: 5,
				Uptime:      "2h 30m",
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:         "server-2",
			TenantID:   tenantID,
			TemplateID: "template-cs2",
			Name:       "Competitive Server",
			GameType:   "cs2",
			Config: models.GameServerConfig{
				Environment: map[string]string{
					"GAME_MODE": "competitive",
					"MAP":       "de_dust2",
				},
			},
			Status: models.GameServerStatus{
				Phase:       "Stopped",
				Message:     "Server is stopped",
				LastUpdated: time.Now().Add(-1 * time.Hour),
				PlayerCount: 0,
				Uptime:      "0m",
			},
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
	}

	return servers, nil
}

// GetTenantActivity retrieves recent activity for a tenant
func (gss *GameServerService) GetTenantActivity(ctx context.Context, tenantID string, limit int) ([]models.Activity, error) {
	// For now, return mock activity data
	// In a real implementation, this would query an audit log table

	activities := []models.Activity{
		{
			ID:        "activity-1",
			Type:      "server_started",
			Message:   "Server 'Survival World' was started",
			Timestamp: time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
		},
		{
			ID:        "activity-2",
			Type:      "user_joined",
			Message:   "Player 'Steve' joined the server",
			Timestamp: time.Now().Add(-35 * time.Minute).Format(time.RFC3339),
		},
		{
			ID:        "activity-3",
			Type:      "server_stopped",
			Message:   "Server 'Competitive Server' was stopped",
			Timestamp: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		},
		{
			ID:        "activity-4",
			Type:      "server_created",
			Message:   "New server 'Survival World' was created",
			Timestamp: time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
		},
		{
			ID:        "activity-5",
			Type:      "role_updated",
			Message:   "Discord roles were synchronized",
			Timestamp: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		},
	}

	// Apply limit if specified
	if limit > 0 && limit < len(activities) {
		activities = activities[:limit]
	}

	return activities, nil
}

// GetTenantDiscordStats retrieves Discord statistics for a tenant
func (gss *GameServerService) GetTenantDiscordStats(ctx context.Context, tenantID string) (*models.DiscordStats, error) {
	var roleCount int64
	var userCount int64
	var lastSync time.Time

	// Count roles
	if err := gss.db.Model(&models.TenantDiscordRole{}).Where("tenant_id = ?", tenantID).Count(&roleCount).Error; err != nil {
		return nil, err
	}

	// Count users
	if err := gss.db.Model(&models.TenantDiscordUser{}).Where("tenant_id = ?", tenantID).Count(&userCount).Error; err != nil {
		return nil, err
	}

	// Get last sync time (latest of users or roles)
	var userLastSync, roleLastSync time.Time
	gss.db.Model(&models.TenantDiscordUser{}).Where("tenant_id = ?", tenantID).Select("MAX(last_sync_at)").Scan(&userLastSync)
	gss.db.Model(&models.TenantDiscordRole{}).Where("tenant_id = ?", tenantID).Select("MAX(updated_at)").Scan(&roleLastSync)
	if userLastSync.After(roleLastSync) {
		lastSync = userLastSync
	} else {
		lastSync = roleLastSync
	}

	stats := &models.DiscordStats{
		MemberCount: int(userCount),
		RoleCount:   int(roleCount),
		LastSync:    lastSync.Format(time.RFC3339),
	}
	return stats, nil
}
