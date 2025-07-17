package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// AuthMiddleware provides authentication middleware
type AuthMiddleware struct {
	authService services.AuthServiceInterface
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService services.AuthServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token and get user
		user, err := m.authService.ValidateAccessToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Invalid or expired token",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			})
			c.Abort()
			return
		}

		// We need to parse the token again to get the session ID
		// This is not ideal but necessary for the current architecture
		claims, err := m.authService.ParseTokenClaims(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Failed to parse token claims",
				Details: map[string]interface{}{
					"error": err.Error(),
				},
			})
			c.Abort()
			return
		}

		// Store user and session info in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("discord_user_id", user.DiscordUserID)
		c.Set("session_id", claims.SessionID)

		c.Next()
	}
}

// OptionalAuth middleware that optionally authenticates
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]

		// Validate token
		user, err := m.authService.ValidateAccessToken(context.Background(), token)
		if err != nil {
			// Don't abort, just continue without user
			c.Next()
			return
		}

		// Store user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("discord_user_id", user.DiscordUserID)

		c.Next()
	}
}

// GetUserFromContext gets the authenticated user from context
func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	u, ok := user.(*models.User)
	return u, ok
}