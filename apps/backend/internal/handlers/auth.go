package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/middleware"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// stateEntry represents a stored OAuth state with expiration
type stateEntry struct {
	value     string
	expiresAt time.Time
}

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService services.AuthServiceInterface
	stateStore  map[string]stateEntry
	stateMutex  sync.RWMutex
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	h := &AuthHandler{
		authService: authService,
		stateStore:  make(map[string]stateEntry),
	}
	
	// Start cleanup goroutine for expired states
	go h.cleanupExpiredStates()
	
	return h
}

// cleanupExpiredStates periodically removes expired state entries
func (h *AuthHandler) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		h.stateMutex.Lock()
		now := time.Now()
		for state, entry := range h.stateStore {
			if now.After(entry.expiresAt) {
				delete(h.stateStore, state)
			}
		}
		h.stateMutex.Unlock()
	}
}

// Login initiates Discord OAuth2 flow
func (h *AuthHandler) Login(c *gin.Context) {
	// Generate state parameter for CSRF protection
	state := uuid.New().String()
	
	// Store state in memory with expiration (10 minutes)
	h.stateMutex.Lock()
	h.stateStore[state] = stateEntry{
		value:     state,
		expiresAt: time.Now().Add(10 * time.Minute),
	}
	h.stateMutex.Unlock()

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
	h.stateMutex.RLock()
	storedEntry, exists := h.stateStore[state]
	h.stateMutex.RUnlock()
	
	if !exists {
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid state parameter",
			Details: map[string]interface{}{
				"error": "State not found or expired",
				"received_state": state,
			},
		})
		return
	}
	
	// Check if state has expired
	if time.Now().After(storedEntry.expiresAt) {
		// Clean up expired state
		h.stateMutex.Lock()
		delete(h.stateStore, state)
		h.stateMutex.Unlock()
		
		c.JSON(http.StatusBadRequest, models.APIError{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid state parameter",
			Details: map[string]interface{}{
				"error": "State expired",
				"received_state": state,
			},
		})
		return
	}
	
	// Remove the used state to prevent replay attacks
	h.stateMutex.Lock()
	delete(h.stateStore, state)
	h.stateMutex.Unlock()

	// Handle the callback
	authResponse, err := h.authService.HandleCallback(c.Request.Context(), code)
	if err != nil {
		// Redirect to login page with error message
		frontendURL := h.getFrontendURL(c)
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error=discord_auth_failed")
		return
	}

	// Redirect to frontend callback with tokens as query parameters
	frontendURL := h.getFrontendURL(c)
	callbackURL := frontendURL + "/auth/callback" +
		"?access_token=" + authResponse.AccessToken +
		"&refresh_token=" + authResponse.RefreshToken +
		"&expires_in=" + fmt.Sprintf("%d", authResponse.ExpiresIn)
	
	c.Redirect(http.StatusTemporaryRedirect, callbackURL)
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

// getFrontendURL returns the frontend URL from configuration
func (h *AuthHandler) getFrontendURL(c *gin.Context) string {
	// Default to localhost for development
	return "http://localhost:3000"
}