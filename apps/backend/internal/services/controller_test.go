package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupControllerTestDB(t *testing.T) (*gorm.DB, func()) {
	db, cleanup := testutils.SetupTestDatabase(t)

	err := db.AutoMigrate(&models.Controller{})
	require.NoError(t, err)

	return db, cleanup
}

func setupControllerService(t *testing.T) (*ControllerService, *gorm.DB, func()) {
	db, cleanup := setupControllerTestDB(t)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
			Issuer: "pteronimbus-test",
		},
		Controller: config.ControllerConfig{
			HandshakeSecret: "",
			HeartbeatTTL:    time.Minute * 5,
			MaxHeartbeatAge: time.Minute * 10,
		},
	}

	jwtService := NewJWTService(cfg)
	service := NewControllerService(db, cfg, jwtService)

	return service, db, cleanup
}

func TestControllerService_Handshake_NewController(t *testing.T) {
	service, _, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	req := &models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Test Cluster",
		Version:     "1.0.0",
		Nonce:       "test-nonce-123",
	}

	resp, err := service.Handshake(ctx, req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.NotEmpty(t, resp.ControllerID)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "Controller registered successfully - awaiting approval", resp.Message)
	assert.Equal(t, "/api/controller/heartbeat", resp.HeartbeatURL)
	assert.Equal(t, 300, resp.HeartbeatTTL) // 5 minutes in seconds

	// Verify controller was created with pending_approval status
	var controller models.Controller
	err = service.db.WithContext(ctx).Where("cluster_id = ?", req.ClusterID).First(&controller).Error
	require.NoError(t, err)
	assert.Equal(t, "pending_approval", controller.Status)
}

func TestControllerService_Handshake_ExistingController(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create an existing controller
	existingController := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174000",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Old Cluster Name",
		Version:        "0.9.0",
		LastHeartbeat:  time.Now().UTC().Add(-time.Hour),
		Status:         "inactive",
		HandshakeToken: "old-token",
	}
	err := db.Create(&existingController).Error
	require.NoError(t, err)

	req := &models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Updated Cluster Name",
		Version:     "1.0.0",
		Nonce:       "test-nonce-456",
	}

	resp, err := service.Handshake(ctx, req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, existingController.ID, resp.ControllerID)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "Controller re-registered successfully", resp.Message)

	// Verify the controller was updated in the database
	var updatedController models.Controller
	err = db.Where("cluster_id = ?", "test-cluster-1").First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated Cluster Name", updatedController.ClusterName)
	assert.Equal(t, "1.0.0", updatedController.Version)
	assert.Equal(t, "inactive", updatedController.Status) // Status should remain inactive, not change to active
}

func TestControllerService_Handshake_ExistingApprovedController(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create an existing approved controller
	existingController := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174001",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Old Cluster Name",
		Version:        "0.9.0",
		LastHeartbeat:  time.Now().UTC().Add(-time.Hour),
		Status:         "active",
		HandshakeToken: "old-token",
	}
	err := db.Create(&existingController).Error
	require.NoError(t, err)

	req := &models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Updated Cluster Name",
		Version:     "1.0.0",
		Nonce:       "test-nonce-456",
	}

	resp, err := service.Handshake(ctx, req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, existingController.ID, resp.ControllerID)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "Controller re-registered successfully", resp.Message)

	// Verify the controller was updated in the database
	var updatedController models.Controller
	err = db.Where("cluster_id = ?", "test-cluster-1").First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "Updated Cluster Name", updatedController.ClusterName)
	assert.Equal(t, "1.0.0", updatedController.Version)
	assert.Equal(t, "active", updatedController.Status) // Status should remain active
}

func TestControllerService_Handshake_WithSecret(t *testing.T) {
	service, _, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Update the service config to include a handshake secret
	service.config.Controller.HandshakeSecret = "test-secret"

	req := &models.HandshakeRequest{
		ClusterID:   "test-cluster-2",
		ClusterName: "Test Cluster",
		Version:     "1.0.0",
		Nonce:       "test-nonce-789",
	}

	resp, err := service.Handshake(ctx, req)
	require.NoError(t, err)
	assert.True(t, resp.Success) // Current implementation accepts any nonce when secret is configured
}

func TestControllerService_Heartbeat_Success(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create a controller first
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174002",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  time.Now().UTC().Add(-time.Hour),
		Status:         "active",
		HandshakeToken: "test-token",
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	req := &models.HeartbeatRequest{
		Status:  "active",
		Message: "Controller is running",
		Metrics: map[string]string{
			"uptime": "1h30m",
		},
		Resources: map[string]int64{
			"memory_usage": 1024,
			"cpu_usage":    50,
		},
	}

	resp, err := service.Heartbeat(ctx, controller.ID, req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Heartbeat received", resp.Message)

	// Verify the controller was updated
	var updatedController models.Controller
	err = db.Where("id = ?", controller.ID).First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "active", updatedController.Status)
}

func TestControllerService_Heartbeat_PendingController(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create a pending controller
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174003",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  time.Now().UTC().Add(-time.Hour),
		Status:         "pending_approval",
		HandshakeToken: "test-token",
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	req := &models.HeartbeatRequest{
		Status:  "active", // Controller tries to set status to active
		Message: "Controller is running",
	}

	resp, err := service.Heartbeat(ctx, controller.ID, req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Heartbeat received - controller awaiting approval", resp.Message)

	// Verify the controller status remains pending_approval
	var updatedController models.Controller
	err = db.Where("id = ?", controller.ID).First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "pending_approval", updatedController.Status) // Status should not change
}

func TestControllerService_Heartbeat_ControllerNotFound(t *testing.T) {
	service, _, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	req := &models.HeartbeatRequest{
		Status:  "active",
		Message: "Controller is running",
	}

	resp, err := service.Heartbeat(ctx, "non-existent-id", req)
	require.NoError(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "Controller not found", resp.Message)
}

func TestControllerService_GetControllerStatus_Online(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create a controller with recent heartbeat
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174004",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  time.Now().UTC(),
		Status:         "active",
		HandshakeToken: "test-token",
		CreatedAt:      time.Now().UTC().Add(-time.Hour),
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	status, err := service.GetControllerStatus(ctx, controller.ID)
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, controller.ID, status.ID)
	assert.Equal(t, controller.ClusterID, status.ClusterID)
	assert.Equal(t, controller.ClusterName, status.ClusterName)
	assert.Equal(t, controller.Version, status.Version)
	assert.Equal(t, controller.Status, status.Status)
	assert.True(t, status.IsOnline)
	assert.NotEmpty(t, status.Uptime)
}

func TestControllerService_GetControllerStatus_Offline(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create a controller with old heartbeat
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174005",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  time.Now().UTC().Add(-time.Hour * 2), // 2 hours ago
		Status:         "active",
		HandshakeToken: "test-token",
		CreatedAt:      time.Now().UTC().Add(-time.Hour * 3),
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	status, err := service.GetControllerStatus(ctx, controller.ID)
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.False(t, status.IsOnline)
	assert.Empty(t, status.Uptime)
}

func TestControllerService_GetControllerStatus_NotFound(t *testing.T) {
	service, _, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	status, err := service.GetControllerStatus(ctx, "non-existent-id")
	require.NoError(t, err)
	assert.Nil(t, status)
}

func TestControllerService_GetAllControllers(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create multiple controllers
	controllers := []models.Controller{
		{
			ID:             "123e4567-e89b-12d3-a456-426614174006",
			ClusterID:      "cluster-1",
			ClusterName:    "Cluster 1",
			Version:        "1.0.0",
			LastHeartbeat:  time.Now().UTC(),
			Status:         "active",
			HandshakeToken: "token-1",
			CreatedAt:      time.Now().UTC().Add(-time.Hour),
		},
		{
			ID:             "123e4567-e89b-12d3-a456-426614174007",
			ClusterID:      "cluster-2",
			ClusterName:    "Cluster 2",
			Version:        "1.1.0",
			LastHeartbeat:  time.Now().UTC().Add(-time.Hour * 2),
			Status:         "error",
			HandshakeToken: "token-2",
			CreatedAt:      time.Now().UTC().Add(-time.Hour * 2),
		},
	}

	for _, controller := range controllers {
		err := db.Create(&controller).Error
		require.NoError(t, err)
	}

	statuses, err := service.GetAllControllers(ctx)
	require.NoError(t, err)
	assert.Len(t, statuses, 2)

	// Find the online controller
	var onlineController *models.ControllerStatus
	for _, status := range statuses {
		if status.ClusterID == "cluster-1" {
			onlineController = status
			break
		}
	}

	assert.NotNil(t, onlineController)
	assert.True(t, onlineController.IsOnline)
	assert.NotEmpty(t, onlineController.Uptime)

	// Find the offline controller
	var offlineController *models.ControllerStatus
	for _, status := range statuses {
		if status.ClusterID == "cluster-2" {
			offlineController = status
			break
		}
	}

	assert.NotNil(t, offlineController)
	assert.False(t, offlineController.IsOnline)
	assert.Empty(t, offlineController.Uptime)
}

func TestControllerService_ValidateControllerToken_Valid(t *testing.T) {
	service, _, cleanup := setupControllerService(t)
	defer cleanup()

	// Generate a valid token
	controllerID := "123e4567-e89b-12d3-a456-426614174008"
	token := service.generateControllerToken(controllerID, "test-cluster")

	// Validate the token
	extractedID, err := service.ValidateControllerToken(token)
	require.NoError(t, err)
	assert.Equal(t, controllerID, extractedID)
}

func TestControllerService_ValidateControllerToken_Invalid(t *testing.T) {
	service, _, cleanup := setupControllerService(t)
	defer cleanup()

	// Test with invalid token
	_, err := service.ValidateControllerToken("invalid-token")
	assert.Error(t, err)
}

func TestControllerService_ValidateControllerToken_WrongType(t *testing.T) {
	service, _, cleanup := setupControllerService(t)
	defer cleanup()

	// Create a token with wrong type
	cfg := service.config
	jwtService := NewJWTService(cfg)

	// Generate a user token instead of controller token
	user := &models.User{
		ID:            "user-id",
		DiscordUserID: "discord-user-id",
		Username:      "testuser",
	}
	userToken, _, err := jwtService.GenerateAccessToken(user, "session-id")
	require.NoError(t, err)

	// Try to validate as controller token
	_, err = service.ValidateControllerToken(userToken)
	assert.Error(t, err)
	// The error could be either "invalid token type" or "token is expired" depending on timing
	assert.True(t, strings.Contains(err.Error(), "invalid token type") || strings.Contains(err.Error(), "token is expired"))
}

func TestControllerService_CleanupInactiveControllers(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create controllers with different heartbeat ages
	controllers := []models.Controller{
		{
			ID:             "123e4567-e89b-12d3-a456-426614174009",
			ClusterID:      "active-cluster",
			ClusterName:    "Active Cluster",
			Version:        "1.0.0",
			LastHeartbeat:  time.Now().UTC(),
			Status:         "active",
			HandshakeToken: "token-1",
		},
		{
			ID:             "123e4567-e89b-12d3-a456-426614174010",
			ClusterID:      "inactive-cluster",
			ClusterName:    "Inactive Cluster",
			Version:        "1.0.0",
			LastHeartbeat:  time.Now().UTC().Add(-time.Hour * 30), // Very old
			Status:         "inactive",
			HandshakeToken: "token-2",
		},
	}

	for _, controller := range controllers {
		err := db.Create(&controller).Error
		require.NoError(t, err)
	}

	// Run cleanup
	err := service.CleanupInactiveControllers(ctx)
	require.NoError(t, err)

	// Verify only the inactive controller was removed
	var remainingControllers []models.Controller
	err = db.Find(&remainingControllers).Error
	require.NoError(t, err)
	assert.Len(t, remainingControllers, 1)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174009", remainingControllers[0].ID)
}

func TestControllerService_ApproveController(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create a pending controller
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174011",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  time.Now().UTC(),
		Status:         "pending_approval",
		HandshakeToken: "test-token",
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	// Approve the controller
	resp, err := service.ApproveController(ctx, controller.ID, "test-user-id")
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Controller approved successfully", resp.Message)

	// Verify the controller was approved
	var updatedController models.Controller
	err = db.Where("id = ?", controller.ID).First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "active", updatedController.Status)
	assert.NotNil(t, updatedController.ApprovedAt)
	assert.Equal(t, "test-user-id", *updatedController.ApprovedBy)
}

func TestControllerService_RejectController(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create a pending controller
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174012",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  time.Now().UTC(),
		Status:         "pending_approval",
		HandshakeToken: "test-token",
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	// Reject the controller
	resp, err := service.RejectController(ctx, controller.ID, "test-user-id", "Test rejection reason")
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Controller rejected: Test rejection reason", resp.Message)

	// Verify the controller was rejected
	var updatedController models.Controller
	err = db.Where("id = ?", controller.ID).First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "rejected", updatedController.Status)
}

func TestControllerService_GetControllerStatus_AutoDegraded(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create an active controller with an old heartbeat (offline)
	oldHeartbeat := time.Now().UTC().Add(-time.Hour) // 1 hour ago
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174013",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  oldHeartbeat,
		Status:         "active", // Currently active but offline
		HandshakeToken: "test-token",
		CreatedAt:      time.Now().UTC().Add(-time.Hour * 2),
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	// Get controller status - should auto-transition to degraded
	status, err := service.GetControllerStatus(ctx, controller.ID)
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "degraded", status.Status) // Should be degraded now
	assert.False(t, status.IsOnline)           // Should be offline

	// Verify the database was updated
	var updatedController models.Controller
	err = db.Where("id = ?", controller.ID).First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "degraded", updatedController.Status)
}

func TestControllerService_GetAllControllers_AutoDegraded(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create multiple controllers with different states
	oldHeartbeat := time.Now().UTC().Add(-time.Hour) // 1 hour ago
	recentHeartbeat := time.Now().UTC().Add(-time.Minute) // 1 minute ago

	controllers := []models.Controller{
		{
			ID:             "123e4567-e89b-12d3-a456-426614174014",
			ClusterID:      "cluster-1",
			ClusterName:    "Cluster 1",
			Version:        "1.0.0",
			LastHeartbeat:  oldHeartbeat,
			Status:         "active", // Will become degraded
			HandshakeToken: "token-1",
		},
		{
			ID:             "123e4567-e89b-12d3-a456-426614174015",
			ClusterID:      "cluster-2",
			ClusterName:    "Cluster 2",
			Version:        "1.0.0",
			LastHeartbeat:  recentHeartbeat,
			Status:         "active", // Will stay active
			HandshakeToken: "token-2",
		},
		{
			ID:             "123e4567-e89b-12d3-a456-426614174016",
			ClusterID:      "cluster-3",
			ClusterName:    "Cluster 3",
			Version:        "1.0.0",
			LastHeartbeat:  oldHeartbeat,
			Status:         "pending_approval", // Won't change
			HandshakeToken: "token-3",
		},
	}

	for _, c := range controllers {
		err := db.Create(&c).Error
		require.NoError(t, err)
	}

	// Get all controllers - should auto-transition offline active ones to degraded
	statuses, err := service.GetAllControllers(ctx)
	require.NoError(t, err)
	assert.Len(t, statuses, 3)

	// Find each controller and verify its status
	controller1 := findControllerByID(statuses, "123e4567-e89b-12d3-a456-426614174014")
	controller2 := findControllerByID(statuses, "123e4567-e89b-12d3-a456-426614174015")
	controller3 := findControllerByID(statuses, "123e4567-e89b-12d3-a456-426614174016")

	assert.NotNil(t, controller1)
	assert.Equal(t, "degraded", controller1.Status)
	assert.False(t, controller1.IsOnline)

	assert.NotNil(t, controller2)
	assert.Equal(t, "active", controller2.Status)
	assert.True(t, controller2.IsOnline)

	assert.NotNil(t, controller3)
	assert.Equal(t, "pending_approval", controller3.Status)
	assert.False(t, controller3.IsOnline)

	// Verify the database was updated for controller-1
	var updatedController models.Controller
	err = db.Where("id = ?", "123e4567-e89b-12d3-a456-426614174014").First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "degraded", updatedController.Status)
}

func TestControllerService_Heartbeat_DegradedToActive(t *testing.T) {
	service, db, cleanup := setupControllerService(t)
	defer cleanup()
	ctx := context.Background()

	// Create a degraded controller
	controller := models.Controller{
		ID:             "123e4567-e89b-12d3-a456-426614174017",
		ClusterID:      "test-cluster-1",
		ClusterName:    "Test Cluster",
		Version:        "1.0.0",
		LastHeartbeat:  time.Now().UTC().Add(-time.Hour),
		Status:         "degraded", // Currently degraded
		HandshakeToken: "test-token",
	}
	err := db.Create(&controller).Error
	require.NoError(t, err)

	// Send heartbeat with active status - should transition back to active
	req := &models.HeartbeatRequest{
		Status:  "active",
		Message: "Controller is running again",
	}

	resp, err := service.Heartbeat(ctx, controller.ID, req)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "Heartbeat received", resp.Message)

	// Verify the controller was updated to active
	var updatedController models.Controller
	err = db.Where("id = ?", controller.ID).First(&updatedController).Error
	require.NoError(t, err)
	assert.Equal(t, "active", updatedController.Status)
}

// Helper function to find controller by ID in status slice
func findControllerByID(statuses []*models.ControllerStatus, id string) *models.ControllerStatus {
	for _, status := range statuses {
		if status.ID == id {
			return status
		}
	}
	return nil
}
