package handlers

import (
	"net/http"

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
				"reason": "Session was created before Discord token integration. Please re-authenticate.",
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
				"reason": "Session was created before Discord token integration. Please re-authenticate.",
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
	botToken := "YOUR_BOT_TOKEN" // This should come from configuration

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

// Helper function to get session from Redis
func (th *TenantHandler) getSessionFromRedis(c *gin.Context, sessionID string) (*models.Session, error) {
	return th.redisService.GetSession(c.Request.Context(), sessionID)
}