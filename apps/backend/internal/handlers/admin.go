package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	adminService *services.AdminService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// CheckAccess checks if the current user has admin access
func (h *AdminHandler) CheckAccess(c *gin.Context) {
	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	// Check if user has superadmin access
	hasAccess, err := h.adminService.CheckSuperAdminAccess(c.Request.Context(), userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to check admin access",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"hasAdminAccess": hasAccess,
	})
}

// GetStats returns admin-level statistics
func (h *AdminHandler) GetStats(c *gin.Context) {
	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	// Check if user has superadmin access
	hasAccess, err := h.adminService.CheckSuperAdminAccess(c.Request.Context(), userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to check admin access",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Insufficient permissions to access admin statistics",
		})
		return
	}

	// Get admin statistics
	stats, err := h.adminService.GetAdminStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get admin statistics",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"stats":   stats,
	})
}

// CleanupInactiveControllers removes inactive controllers
func (h *AdminHandler) CleanupInactiveControllers(c *gin.Context) {
	// Get authenticated user
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not authenticated",
		})
		return
	}

	userModel := user.(*models.User)

	// Check if user has superadmin access
	hasAccess, err := h.adminService.CheckSuperAdminAccess(c.Request.Context(), userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to check admin access",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "FORBIDDEN",
			Message: "Insufficient permissions to perform cleanup",
		})
		return
	}

	// Perform cleanup
	err = h.adminService.CleanupInactiveControllers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to cleanup inactive controllers",
			Details: map[string]interface{}{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Inactive controllers cleaned up successfully",
	})
}
