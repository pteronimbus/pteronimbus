package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// ControllerMiddleware provides authentication for controller endpoints
type ControllerMiddleware struct {
	controllerService *services.ControllerService
}

// NewControllerMiddleware creates a new controller middleware
func NewControllerMiddleware(controllerService *services.ControllerService) *ControllerMiddleware {
	return &ControllerMiddleware{
		controllerService: controllerService,
	}
}

// RequireControllerAuth returns a middleware that requires valid controller authentication
func (m *ControllerMiddleware) RequireControllerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract controller token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Missing authorization header",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		controllerID, err := m.controllerService.ValidateControllerToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store controller ID in context for later use
		c.Set("controller_id", controllerID)
		c.Next()
	}
}

// GetControllerIDFromContext extracts the controller ID from the Gin context
func GetControllerIDFromContext(c *gin.Context) (string, bool) {
	controllerID, exists := c.Get("controller_id")
	if !exists {
		return "", false
	}
	return controllerID.(string), true
}
