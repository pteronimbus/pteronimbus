package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// TenantHandler handles tenant-related HTTP requests
type TenantHandler struct {
	tenantService  services.TenantServiceInterface
	discordService services.DiscordServiceInterface
	authService    services.AuthServiceInterface
	redisService   services.RedisServiceInterface
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantService services.TenantServiceInterface, discordService services.DiscordServiceInterface, authService services.AuthServiceInterface, redisService services.RedisServiceInterface) *TenantHandler {
	return &TenantHandler{
		tenantService:  tenantService,
		discordService: discordService,
		authService:    authService,
		redisService:   redisService,
	}
}

// GetUserTenants retrieves all tenants a user has access to
func (th *TenantHandler) GetUserTenants(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	tenants, err := th.tenantService.GetUserTenants(c.Request.Context(), userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get user tenants",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenants": tenants,
	})
}

// GetAvailableGuilds retrieves Discord guilds where user can install Pteronimbus
func (th *TenantHandler) GetAvailableGuilds(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	_ = user.(*models.User)

	// Get user's Discord access token from session
	sessionID, exists := c.Get("session_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "Session not found",
		})
		return
	}

	session, err := th.getSessionFromRedis(c, sessionID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "Failed to get session",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Check if Discord access token exists
	if session.DiscordAccessToken == "" {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "DISCORD_TOKEN_MISSING",
			Message: "Discord access token not found in session. Please log in again.",
			Details: map[string]interface{}{
				"session_id": sessionID,
				"reason":     "Session was created before Discord token integration. Please re-authenticate.",
			},
		})
		return
	}

	// Get user's Discord guilds
	guilds, err := th.discordService.GetUserGuilds(c.Request.Context(), session.DiscordAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "DISCORD_API_ERROR",
			Message: "Failed to get Discord guilds",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Filter guilds where user has manage server permissions
	var availableGuilds []models.DiscordGuild
	for _, guild := range guilds {
		if th.tenantService.CheckManageServerPermission(&guild) {
			availableGuilds = append(availableGuilds, guild)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"guilds": availableGuilds,
	})
}

// CreateTenant creates a new tenant from a Discord guild
func (th *TenantHandler) CreateTenant(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	var req struct {
		GuildID string `json:"guild_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Get user's Discord access token
	sessionID, exists := c.Get("session_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "Session not found",
		})
		return
	}

	session, err := th.getSessionFromRedis(c, sessionID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "Failed to get session",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Check if Discord access token exists
	if session.DiscordAccessToken == "" {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "DISCORD_TOKEN_MISSING",
			Message: "Discord access token not found in session. Please log in again.",
			Details: map[string]interface{}{
				"session_id": sessionID,
				"reason":     "Session was created before Discord token integration. Please re-authenticate.",
			},
		})
		return
	}

	// Get user's Discord guilds to verify they have access
	guilds, err := th.discordService.GetUserGuilds(c.Request.Context(), session.DiscordAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "DISCORD_API_ERROR",
			Message: "Failed to get Discord guilds",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Find the requested guild and verify permissions
	var targetGuild *models.DiscordGuild
	for _, guild := range guilds {
		if guild.ID == req.GuildID {
			targetGuild = &guild
			break
		}
	}

	if targetGuild == nil {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Guild not found or access denied",
		})
		return
	}

	if !th.tenantService.CheckManageServerPermission(targetGuild) {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Insufficient permissions to install Pteronimbus to this server",
		})
		return
	}

	// Create tenant
	tenant, err := th.tenantService.CreateTenant(c.Request.Context(), targetGuild, userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create tenant",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Add the user to the tenant as owner
	err = th.tenantService.AddUserToTenant(c.Request.Context(), userModel.ID, tenant.ID, []string{"owner"}, []string{"*"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to add user to tenant",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tenant": tenant,
	})
}

// GetTenant retrieves a specific tenant
func (th *TenantHandler) GetTenant(c *gin.Context) {
	tenantID := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Tenant ID is required",
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	// Check if user has access to this tenant
	hasAccess, err := th.tenantService.HasPermission(c.Request.Context(), userModel.ID, tenantID, "read")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to check tenant access",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Access denied to this tenant",
		})
		return
	}

	tenant, err := th.tenantService.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIError{
			Code:    "NOT_FOUND",
			Message: "Tenant not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant": tenant,
	})
}

// UpdateTenantConfig updates tenant configuration
func (th *TenantHandler) UpdateTenantConfig(c *gin.Context) {
	tenantID := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Tenant ID is required",
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	// Check if user has manage permissions for this tenant
	hasAccess, err := th.tenantService.HasPermission(c.Request.Context(), userModel.ID, tenantID, "manage")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to check tenant access",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Insufficient permissions to manage this tenant",
		})
		return
	}

	var config models.TenantConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	err = th.tenantService.UpdateTenantConfig(c.Request.Context(), tenantID, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update tenant config",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant configuration updated successfully",
	})
}

// SyncTenantData synchronizes Discord roles and users for a tenant
func (th *TenantHandler) SyncTenantData(c *gin.Context) {
	tenantID := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Tenant ID is required",
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	// Check if user has manage permissions for this tenant
	hasAccess, err := th.tenantService.HasPermission(c.Request.Context(), userModel.ID, tenantID, "manage")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to check tenant access",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Insufficient permissions to manage this tenant",
		})
		return
	}

	// Get bot token from config or environment
	// For now, we'll use a placeholder - this should be configured properly
	botToken := th.discordService.(*services.DiscordService).BotToken()
	if botToken == "" {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "DISCORD_SYNC_ERROR",
			Message: "Bot token not configured",
		})
		return
	}

	// Sync Discord roles
	err = th.tenantService.SyncDiscordRoles(c.Request.Context(), tenantID, botToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "DISCORD_SYNC_ERROR",
			Message: "Failed to sync Discord roles",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Sync Discord users
	err = th.tenantService.SyncDiscordUsers(c.Request.Context(), tenantID, botToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "DISCORD_SYNC_ERROR",
			Message: "Failed to sync Discord users",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant data synchronized successfully",
	})
}

// DeleteTenant deletes a tenant
func (th *TenantHandler) DeleteTenant(c *gin.Context) {
	tenantID := c.Param("id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Tenant ID is required",
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	// Check if user is the owner of this tenant
	tenant, err := th.tenantService.GetTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIError{
			Code:    "NOT_FOUND",
			Message: "Tenant not found",
		})
		return
	}

	if tenant.OwnerID != userModel.ID {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Only the tenant owner can delete the tenant",
		})
		return
	}

	err = th.tenantService.DeleteTenant(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to delete tenant",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tenant deleted successfully",
	})
}

// GetBotStatus checks if the bot is present in the tenant's guild and has required permissions
func (th *TenantHandler) GetBotStatus(c *gin.Context) {
	ctx := c.Request.Context()
	tenantId := c.Param("id")

	// Get tenant
	tenant, err := th.tenantService.GetTenant(ctx, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tenant not found"})
		return
	}
	guildID := tenant.DiscordServerID

	// Get bot user ID from config
	botToken := th.discordService.(*services.DiscordService).BotToken()
	if botToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bot token not configured"})
		return
	}

	// Get bot user info
	botUser, err := th.discordService.GetUserInfo(ctx, botToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get bot user info"})
		return
	}

	// Check if bot is a member of the guild
	member, err := th.discordService.GetGuildMember(ctx, botToken, guildID, botUser.ID)
	present := err == nil && member != nil

	missing := []string{}
	if present {
		// Discord permissions bitfield for required permissions
		const requiredPerms int64 = 277293902870

		// Get the bot's member object in the guild (already fetched as 'member')
		// The member.Roles and member.Permissions fields may be relevant

		// Calculate the bot's effective permissions in the guild
		// This requires fetching all roles and summing the permissions for the bot's roles
		roles, err := th.discordService.GetGuildRoles(ctx, botToken, guildID)
		var effectivePerms int64 = 0
		if err == nil && roles != nil {
			for _, role := range roles {
				for _, botRoleID := range member.Roles {
					if role.ID == botRoleID {
						// Permissions is a string, so parse it to int64
						permInt, perr := strconv.ParseInt(role.Permissions, 10, 64)
						if perr == nil {
							effectivePerms |= permInt
						}
					}
				}
			}
		}
		// If the bot has the ADMINISTRATOR permission, it has all permissions
		const permAdministrator int64 = 0x00000008
		if (effectivePerms & permAdministrator) == permAdministrator {
			// Bot has all permissions
		} else {
			// Check for each required permission bit
			permMap := map[int64]string{
				0x00000008:   "ADMINISTRATOR",
				0x00000020:   "MANAGE_CHANNELS",
				0x00000040:   "MANAGE_GUILD",
				0x00000080:   "ADD_REACTIONS",
				0x00000400:   "VIEW_AUDIT_LOG",
				0x00000800:   "PRIORITY_SPEAKER",
				0x00001000:   "STREAM",
				0x00002000:   "VIEW_CHANNEL",
				0x00004000:   "SEND_MESSAGES",
				0x00008000:   "SEND_TTS_MESSAGES",
				0x00010000:   "MANAGE_MESSAGES",
				0x00020000:   "EMBED_LINKS",
				0x00040000:   "ATTACH_FILES",
				0x00080000:   "READ_MESSAGE_HISTORY",
				0x00100000:   "MENTION_EVERYONE",
				0x00200000:   "USE_EXTERNAL_EMOJIS",
				0x00400000:   "VIEW_GUILD_INSIGHTS",
				0x01000000:   "MANAGE_ROLES",
				0x02000000:   "MANAGE_WEBHOOKS",
				0x04000000:   "MANAGE_EMOJIS_AND_STICKERS",
				0x08000000:   "USE_APPLICATION_COMMANDS",
				0x10000000:   "REQUEST_TO_SPEAK",
				0x20000000:   "MANAGE_EVENTS",
				0x40000000:   "MANAGE_THREADS",
				0x80000000:   "CREATE_PUBLIC_THREADS",
				0x100000000:  "CREATE_PRIVATE_THREADS",
				0x200000000:  "USE_EXTERNAL_STICKERS",
				0x400000000:  "SEND_MESSAGES_IN_THREADS",
				0x800000000:  "START_EMBEDDED_ACTIVITIES",
				0x1000000000: "MODERATE_MEMBERS",
			}
			for bit, name := range permMap {
				if (requiredPerms&bit) == bit && (effectivePerms&bit) != bit {
					missing = append(missing, name)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"present":            present,
		"missingPermissions": missing,
	})
}

// Helper function to get session from Redis
func (th *TenantHandler) getSessionFromRedis(c *gin.Context, sessionID string) (*models.Session, error) {
	return th.redisService.GetSession(c.Request.Context(), sessionID)
}
