package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/middleware"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService services.AuthServiceInterface
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login initiates Discord OAuth2 flow
func (h *AuthHandler) Login(c *gin.Context) {
	// Generate state parameter for CSRF protection
	state := uuid.New().String()
	
	// Store state in session/cookie for validation (simplified for now)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true) // 10 minutes

	// Get Discord authorization URL
	authURL := h.authService.GetAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// Callback handles Discord OAuth2 callback
func (h *AuthHandler) Callback(c *gin.Context) {
	// Get code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Authorization code is required",
		})
		return
	}

	if state == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "State parameter is required",
		})
		return
	}

	// Validate state parameter (CSRF protection)
	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid state parameter",
		})
		return
	}

	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// Handle the callback
	authResponse, err := h.authService.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "DISCORD_API_ERROR",
			Message: "Failed to authenticate with Discord",
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// Refresh refreshes access token using refresh token
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid request body",
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Refresh token is required",
		})
		return
	}

	authResponse, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "Failed to refresh token",
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// Me returns current user information
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIError{
			Code:    "UNAUTHORIZED",
			Message: "User not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Logout invalidates the current session
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Authorization header required",
		})
		return
	}

	// Extract token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid authorization header format",
		})
		return
	}

	token := parts[1]

	err := h.authService.Logout(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIError{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to logout",
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}