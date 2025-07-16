package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// TenantMiddleware handles tenant context for API requests
type TenantMiddleware struct {
	tenantService services.TenantServiceInterface
}

// NewTenantMiddleware creates a new tenant middleware
func NewTenantMiddleware(tenantService services.TenantServiceInterface) *TenantMiddleware {
	return &TenantMiddleware{
		tenantService: tenantService,
	}
}

// RequireTenant middleware ensures a valid tenant context is present
func (tm *TenantMiddleware) RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from header or query parameter
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}

		if tenantID == "" {
			c.JSON(http.StatusBadRequest, models.APIError{
				Code:    "MISSING_TENANT",
				Message: "Tenant ID is required",
			})
			c.Abort()
			return
		}

		// Get authenticated user
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Code:    "UNAUTHORIZED",
				Message: "User not authenticated",
			})
			c.Abort()
			return
		}

		userModel := user.(*models.User)

		// Check if user has access to this tenant
		hasAccess, err := tm.tenantService.HasPermission(c.Request.Context(), userModel.ID, tenantID, "read")
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to check tenant access",
				Details: map[string]interface{}{"error": err.Error()},
			})
			c.Abort()
			return
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, models.APIError{
				Code:    "FORBIDDEN",
				Message: "Access denied to this tenant",
			})
			c.Abort()
			return
		}

		// Get tenant details
		tenant, err := tm.tenantService.GetTenant(c.Request.Context(), tenantID)
		if err != nil {
			c.JSON(http.StatusNotFound, models.APIError{
				Code:    "TENANT_NOT_FOUND",
				Message: "Tenant not found",
			})
			c.Abort()
			return
		}

		// Set tenant in context
		c.Set("tenant", tenant)
		c.Set("tenant_id", tenantID)

		c.Next()
	}
}

// RequirePermission middleware ensures user has specific permission in tenant
func (tm *TenantMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from context (should be set by RequireTenant middleware)
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.JSON(http.StatusBadRequest, models.APIError{
				Code:    "MISSING_TENANT_CONTEXT",
				Message: "Tenant context not found",
			})
			c.Abort()
			return
		}

		// Get authenticated user
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Code:    "UNAUTHORIZED",
				Message: "User not authenticated",
			})
			c.Abort()
			return
		}

		userModel := user.(*models.User)

		// Check if user has the required permission
		hasPermission, err := tm.tenantService.HasPermission(c.Request.Context(), userModel.ID, tenantID.(string), permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to check permission",
				Details: map[string]interface{}{"error": err.Error()},
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, models.APIError{
				Code:    "INSUFFICIENT_PERMISSIONS",
				Message: "Insufficient permissions for this operation",
				Details: map[string]interface{}{"required_permission": permission},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalTenant middleware adds tenant context if provided but doesn't require it
func (tm *TenantMiddleware) OptionalTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from header or query parameter
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}

		if tenantID == "" {
			c.Next()
			return
		}

		// Get authenticated user
		user, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		userModel := user.(*models.User)

		// Check if user has access to this tenant
		hasAccess, err := tm.tenantService.HasPermission(c.Request.Context(), userModel.ID, tenantID, "read")
		if err != nil || !hasAccess {
			c.Next()
			return
		}

		// Get tenant details
		tenant, err := tm.tenantService.GetTenant(c.Request.Context(), tenantID)
		if err != nil {
			c.Next()
			return
		}

		// Set tenant in context
		c.Set("tenant", tenant)
		c.Set("tenant_id", tenantID)

		c.Next()
	}
}

// TenantOwnerOnly middleware ensures only tenant owner can access
func (tm *TenantMiddleware) TenantOwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant from context (should be set by RequireTenant middleware)
		tenant, exists := c.Get("tenant")
		if !exists {
			c.JSON(http.StatusBadRequest, models.APIError{
				Code:    "MISSING_TENANT_CONTEXT",
				Message: "Tenant context not found",
			})
			c.Abort()
			return
		}

		// Get authenticated user
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIError{
				Code:    "UNAUTHORIZED",
				Message: "User not authenticated",
			})
			c.Abort()
			return
		}

		userModel := user.(*models.User)
		tenantModel := tenant.(*models.Tenant)

		// Check if user is the tenant owner
		if tenantModel.OwnerID != userModel.ID {
			c.JSON(http.StatusForbidden, models.APIError{
				Code:    "OWNER_ONLY",
				Message: "This operation is restricted to tenant owners",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}