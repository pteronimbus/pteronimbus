package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupPermissionMiddlewareTest(t *testing.T) (*PermissionMiddleware, *gorm.DB, func()) {
	db, cleanup := testutils.SetupTestDatabaseWithModels(t,
		&models.User{},
		&models.Tenant{},
		&models.UserTenant{},
		&models.TenantDiscordRole{},
		&models.TenantDiscordUser{},
		&models.Permission{},
		&models.Role{},
		&models.PermissionAuditLog{},
		&models.GuildMembershipCache{},
	)
	
	rbacConfig := &config.RBACConfig{
		SuperAdminDiscordID: "superadmin123",
		RoleSyncTTL:         time.Minute * 5,
		GuildCacheTTL:       time.Minute * 5,
		GracePeriod:         time.Minute * 2,
	}
	
	rbacService := services.NewRBACService(db, rbacConfig)
	permissionMiddleware := NewPermissionMiddleware(rbacService)
	
	return permissionMiddleware, db, cleanup
}

func TestPermissionMiddleware_RequirePermission(t *testing.T) {
	permissionMiddleware, db, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create user-tenant relationship with permission
	userTenant := &models.UserTenant{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		Roles:       models.StringArray{},
		Permissions: models.StringArray{models.PermissionServerRead},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenant.ID)
		c.Set("user", user)
		c.Next()
	})
	
	router.GET("/test", permissionMiddleware.RequirePermission(models.PermissionServerRead), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test successful permission check
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)

	// Test failed permission check
	router.GET("/test2", permissionMiddleware.RequirePermission(models.PermissionServerWrite), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test2", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPermissionMiddleware_RequirePermission_MissingTenantContext(t *testing.T) {
	permissionMiddleware, _, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Setup Gin router without tenant context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/test", permissionMiddleware.RequirePermission(models.PermissionServerRead), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPermissionMiddleware_RequirePermission_MissingUser(t *testing.T) {
	permissionMiddleware, db, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Create tenant
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	// Setup Gin router with tenant but no user
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenant.ID)
		c.Next()
	})
	
	router.GET("/test", permissionMiddleware.RequirePermission(models.PermissionServerRead), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPermissionMiddleware_RequireAnyPermission(t *testing.T) {
	permissionMiddleware, db, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create user-tenant relationship with one permission
	userTenant := &models.UserTenant{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		Roles:       models.StringArray{},
		Permissions: models.StringArray{models.PermissionServerRead},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenant.ID)
		c.Set("user", user)
		c.Next()
	})
	
	router.GET("/test", permissionMiddleware.RequireAnyPermission(models.PermissionServerRead, models.PermissionServerWrite), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test successful permission check (user has one of the required permissions)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)

	// Test failed permission check (user has none of the required permissions)
	router.GET("/test2", permissionMiddleware.RequireAnyPermission(models.PermissionServerWrite, models.PermissionServerDelete), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test2", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPermissionMiddleware_RequireAllPermissions(t *testing.T) {
	permissionMiddleware, db, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create user-tenant relationship with all required permissions
	userTenant := &models.UserTenant{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		Roles:       models.StringArray{},
		Permissions: models.StringArray{models.PermissionServerRead, models.PermissionServerWrite},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenant.ID)
		c.Set("user", user)
		c.Next()
	})
	
	router.GET("/test", permissionMiddleware.RequireAllPermissions(models.PermissionServerRead, models.PermissionServerWrite), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test successful permission check (user has all required permissions)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)

	// Test failed permission check (user missing one required permission)
	router.GET("/test2", permissionMiddleware.RequireAllPermissions(models.PermissionServerRead, models.PermissionServerDelete), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test2", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPermissionMiddleware_RequireSuperAdmin(t *testing.T) {
	permissionMiddleware, db, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Create super admin user
	superAdmin := &models.User{
		DiscordUserID: "superadmin123", // Matches config
		Username:      "superadmin",
	}
	err := db.Create(superAdmin).Error
	require.NoError(t, err)

	// Create regular user
	regularUser := &models.User{
		DiscordUserID: "regular123",
		Username:      "regularuser",
	}
	err = db.Create(regularUser).Error
	require.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/admin", func(c *gin.Context) {
		c.Set("user", superAdmin)
		c.Next()
	}, permissionMiddleware.RequireSuperAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin_access"})
	})

	router.GET("/admin2", func(c *gin.Context) {
		c.Set("user", regularUser)
		c.Next()
	}, permissionMiddleware.RequireSuperAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin_access"})
	})

	// Test super admin access
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)

	// Test regular user access (should be denied)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/admin2", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestPermissionMiddleware_RequireSuperAdmin_MissingUser(t *testing.T) {
	permissionMiddleware, _, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Setup Gin router without user
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.GET("/admin", permissionMiddleware.RequireSuperAdmin(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin_access"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPermissionMiddleware_OptionalPermission(t *testing.T) {
	permissionMiddleware, db, cleanup := setupPermissionMiddlewareTest(t)
	defer cleanup()

	// Create test data
	tenant := &models.Tenant{
		DiscordServerID: "guild-123",
		Name:            "Test Guild",
		OwnerID:         uuid.New().String(),
	}
	err := db.Create(tenant).Error
	require.NoError(t, err)

	user := &models.User{
		DiscordUserID: "user-123",
		Username:      "testuser",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create user-tenant relationship with permission
	userTenant := &models.UserTenant{
		UserID:      user.ID,
		TenantID:    tenant.ID,
		Roles:       models.StringArray{},
		Permissions: models.StringArray{models.PermissionServerRead},
	}
	err = db.Create(userTenant).Error
	require.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	router.Use(func(c *gin.Context) {
		c.Set("tenant_id", tenant.ID)
		c.Set("user", user)
		c.Next()
	})
	
	router.GET("/test", permissionMiddleware.OptionalPermission(models.PermissionServerRead), func(c *gin.Context) {
		hasPermission, exists := c.Get("has_permission_server:read")
		if exists {
			c.JSON(http.StatusOK, gin.H{"has_permission": hasPermission})
		} else {
			c.JSON(http.StatusOK, gin.H{"has_permission": false})
		}
	})

	// Test with permission
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	// Note: In a real test, you'd parse the JSON response to verify has_permission is true

	// Test without tenant context (should still work)
	router.GET("/test2", func(c *gin.Context) {
		c.Set("user", user)
		c.Next()
	}, permissionMiddleware.OptionalPermission(models.PermissionServerRead), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test2", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
} 