package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
)

// PermissionMiddleware handles permission-based access control
type PermissionMiddleware struct {
	rbacService *services.RBACService
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware(rbacService *services.RBACService) *PermissionMiddleware {
	return &PermissionMiddleware{
		rbacService: rbacService,
	}
}

// RequirePermission middleware ensures user has specific permission in tenant
func (pm *PermissionMiddleware) RequirePermission(permission string) gin.HandlerFunc {
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
		hasPermission, err := pm.rbacService.HasPermission(c.Request.Context(), userModel.ID, tenantID.(string), permission)
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

// RequireAnyPermission middleware ensures user has at least one of the specified permissions
func (pm *PermissionMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from context
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

		// Check if user has any of the required permissions
		for _, permission := range permissions {
			hasPermission, err := pm.rbacService.HasPermission(c.Request.Context(), userModel.ID, tenantID.(string), permission)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.APIError{
					Code:    "INTERNAL_ERROR",
					Message: "Failed to check permission",
					Details: map[string]interface{}{"error": err.Error()},
				})
				c.Abort()
				return
			}

			if hasPermission {
				c.Next()
				return
			}
		}

		// User doesn't have any of the required permissions
		c.JSON(http.StatusForbidden, models.APIError{
			Code:    "INSUFFICIENT_PERMISSIONS",
			Message: "Insufficient permissions for this operation",
			Details: map[string]interface{}{"required_permissions": permissions},
		})
		c.Abort()
	}
}

// RequireAllPermissions middleware ensures user has all of the specified permissions
func (pm *PermissionMiddleware) RequireAllPermissions(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from context
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

		// Check if user has all of the required permissions
		for _, permission := range permissions {
			hasPermission, err := pm.rbacService.HasPermission(c.Request.Context(), userModel.ID, tenantID.(string), permission)
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
					Details: map[string]interface{}{"required_permissions": permissions, "missing_permission": permission},
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireSuperAdmin middleware ensures user is a super admin
func (pm *PermissionMiddleware) RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		// Check if user is super admin
		isSuperAdmin, err := pm.rbacService.IsSuperAdmin(c.Request.Context(), userModel.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to check super admin status",
				Details: map[string]interface{}{"error": err.Error()},
			})
			c.Abort()
			return
		}

		if !isSuperAdmin {
			c.JSON(http.StatusForbidden, models.APIError{
				Code:    "INSUFFICIENT_PERMISSIONS",
				Message: "Super admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalPermission middleware adds permission context but doesn't require it
func (pm *PermissionMiddleware) OptionalPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get tenant ID from context
		tenantID, exists := c.Get("tenant_id")
		if !exists {
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

		// Check if user has the permission
		hasPermission, err := pm.rbacService.HasPermission(c.Request.Context(), userModel.ID, tenantID.(string), permission)
		if err != nil {
			// Log error but don't fail the request
			c.Error(err)
			c.Next()
			return
		}

		// Set permission context
		c.Set("has_permission_"+permission, hasPermission)

		c.Next()
	}
} 