package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"gorm.io/gorm"
)

// SyncService handles syncing data from Discord.
type SyncService struct {
	db      *gorm.DB
	discord *discordgo.Session
}

// NewSyncService creates a new SyncService.
func NewSyncService(db *gorm.DB, discord *discordgo.Session) *SyncService {
	return &SyncService{db: db, discord: discord}
}

// SyncRoles syncs roles from a Discord server to a tenant.
func (s *SyncService) SyncRoles(tenantID string, guildID string) error {
	roles, err := s.discord.GuildRoles(guildID)
	if err != nil {
		return err
	}

	for _, role := range roles {
		dbRole := models.TenantDiscordRole{
			DiscordRoleID: role.ID,
			TenantID:      tenantID,
			Name:          role.Name,
			Color:         role.Color,
			Position:      role.Position,
			Mentionable:   role.Mentionable,
			Hoist:         role.Hoist,
		}

		if err := s.db.Where(models.TenantDiscordRole{DiscordRoleID: role.ID, TenantID: tenantID}).Assign(dbRole).FirstOrCreate(&dbRole).Error; err != nil {
			return err
		}
	}

	return nil
}

// SyncUsers syncs users from a Discord server to a tenant.
func (s *SyncService) SyncUsers(tenantID string, guildID string) error {
	members, err := s.discord.GuildMembers(guildID, "", 1000)
	if err != nil {
		return err
	}

	for _, member := range members {
		joinedAt := member.JoinedAt
		dbUser := models.TenantDiscordUser{
			DiscordUserID: member.User.ID,
			TenantID:      tenantID,
			Username:      member.User.Username,
			DisplayName:   member.Nick,
			Avatar:        member.User.AvatarURL("128"),
			JoinedAt:      &joinedAt,
		}

		if err := s.db.Where(models.TenantDiscordUser{DiscordUserID: member.User.ID, TenantID: tenantID}).Assign(dbUser).FirstOrCreate(&dbUser).Error; err != nil {
			return err
		}
	}

	return nil
}
