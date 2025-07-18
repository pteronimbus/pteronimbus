package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// GameServerHandler handles game server-related HTTP requests
type GameServerHandler struct {
	gameServerService services.GameServerServiceInterface
	tenantService     services.TenantServiceInterface
}

// NewGameServerHandler creates a new game server handler
func NewGameServerHandler(gameServerService services.GameServerServiceInterface, tenantService services.TenantServiceInterface) *GameServerHandler {
	return &GameServerHandler{
		gameServerService: gameServerService,
		tenantService:     tenantService,
	}
}

// GetTenantServers retrieves all game servers for a tenant
func (gsh *GameServerHandler) GetTenantServers(c *gin.Context) {
	tenant, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "TENANT_REQUIRED",
			Message: "Tenant context is required",
		})
		return
	}

	tenantModel := tenant.(*models.Tenant)

	servers, err := gsh.gameServerService.GetTenantServers(c.Request.Context(), tenantModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get tenant servers",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"servers": servers,
	})
}

// GetTenantActivity retrieves recent activity for a tenant
func (gsh *GameServerHandler) GetTenantActivity(c *gin.Context) {
	tenant, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "TENANT_REQUIRED",
			Message: "Tenant context is required",
		})
		return
	}

	tenantModel := tenant.(*models.Tenant)

	// Get limit from query parameter, default to 0 (no limit)
	limit := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	activities, err := gsh.gameServerService.GetTenantActivity(c.Request.Context(), tenantModel.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get tenant activity",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
	})
}

// GetTenantDiscordStats retrieves Discord statistics for a tenant
func (gsh *GameServerHandler) GetTenantDiscordStats(c *gin.Context) {
	tenant, exists := c.Get("tenant")
	if !exists {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "TENANT_REQUIRED",
			Message: "Tenant context is required",
		})
		return
	}

	tenantModel := tenant.(*models.Tenant)

	stats, err := gsh.gameServerService.GetTenantDiscordStats(c.Request.Context(), tenantModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get tenant Discord stats",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}