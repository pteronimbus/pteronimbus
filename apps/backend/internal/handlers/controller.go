package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// ControllerHandler handles controller-related HTTP requests
type ControllerHandler struct {
	controllerService *services.ControllerService
}

// NewControllerHandler creates a new controller handler
func NewControllerHandler(controllerService *services.ControllerService) *ControllerHandler {
	return &ControllerHandler{
		controllerService: controllerService,
	}
}

// Handshake handles controller registration and authentication
func (h *ControllerHandler) Handshake(c *gin.Context) {
	var req models.HandshakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.controllerService.Handshake(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error: " + err.Error(),
		})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusUnauthorized, response)
	}
}

// Heartbeat handles controller heartbeat updates
func (h *ControllerHandler) Heartbeat(c *gin.Context) {
	// Extract controller token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Missing authorization header",
		})
		return
	}

	// Extract token from "Bearer <token>" format
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid authorization header format",
		})
		return
	}

	token := tokenParts[1]
	controllerID, err := h.controllerService.ValidateControllerToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid or expired token",
		})
		return
	}

	var req models.HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format: " + err.Error(),
		})
		return
	}

	response, err := h.controllerService.Heartbeat(c.Request.Context(), controllerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error: " + err.Error(),
		})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusNotFound, response)
	}
}

// GetControllerStatus returns the status of a specific controller
func (h *ControllerHandler) GetControllerStatus(c *gin.Context) {
	controllerID := c.Param("id")
	if controllerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Controller ID is required",
		})
		return
	}

	status, err := h.controllerService.GetControllerStatus(c.Request.Context(), controllerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error: " + err.Error(),
		})
		return
	}

	if status == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Controller not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"controller": status,
	})
}

// GetAllControllers returns all registered controllers
func (h *ControllerHandler) GetAllControllers(c *gin.Context) {
	controllers, err := h.controllerService.GetAllControllers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error: " + err.Error(),
		})
		return
	}

	// Ensure we always return an array, even if empty
	if controllers == nil {
		controllers = []*models.ControllerStatus{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"controllers": controllers,
	})
}
