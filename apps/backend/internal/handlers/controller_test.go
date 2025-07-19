package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupControllerHandlerTest(t *testing.T) (*ControllerHandler, *gin.Engine) {
	// Setup database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Controller{})
	require.NoError(t, err)

	// Setup config
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

	// Setup services
	jwtService := services.NewJWTService(cfg)
	controllerService := services.NewControllerService(db, cfg, jwtService)

	// Setup handler
	handler := NewControllerHandler(controllerService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	return handler, router
}

func TestControllerHandler_Handshake_Success(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.POST("/handshake", handler.Handshake)

	// Create request
	reqBody := models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Test Cluster",
		Version:     "1.0.0",
		Nonce:       "test-nonce-123",
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/handshake", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HandshakeResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.NotEmpty(t, response.ControllerID)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "Controller registered successfully", response.Message)
	assert.Equal(t, "/api/controller/heartbeat", response.HeartbeatURL)
	assert.Equal(t, 300, response.HeartbeatTTL)
}

func TestControllerHandler_Handshake_InvalidRequest(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.POST("/handshake", handler.Handshake)

	// Create invalid request (missing required fields)
	reqBody := map[string]interface{}{
		"cluster_id": "test-cluster-1",
		// Missing cluster_name, version, and nonce
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/handshake", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["message"].(string), "Invalid request format")
}

func TestControllerHandler_Handshake_ExistingController(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.POST("/handshake", handler.Handshake)

	// First handshake
	reqBody1 := models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Test Cluster",
		Version:     "1.0.0",
		Nonce:       "test-nonce-123",
	}

	reqBytes1, err := json.Marshal(reqBody1)
	require.NoError(t, err)

	req1 := httptest.NewRequest("POST", "/handshake", bytes.NewBuffer(reqBytes1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)

	// Second handshake with same cluster ID but different details
	reqBody2 := models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Updated Cluster Name",
		Version:     "1.1.0",
		Nonce:       "test-nonce-456",
	}

	reqBytes2, err := json.Marshal(reqBody2)
	require.NoError(t, err)

	req2 := httptest.NewRequest("POST", "/handshake", bytes.NewBuffer(reqBytes2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Assert response
	assert.Equal(t, http.StatusOK, w2.Code)

	var response models.HandshakeResponse
	err = json.Unmarshal(w2.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.Equal(t, "Controller re-registered successfully", response.Message)
}

func TestControllerHandler_Heartbeat_Success(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup routes
	router.POST("/handshake", handler.Handshake)
	router.POST("/heartbeat", handler.Heartbeat)

	// First, perform handshake to get a token
	handshakeReq := models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Test Cluster",
		Version:     "1.0.0",
		Nonce:       "test-nonce-123",
	}

	handshakeBytes, err := json.Marshal(handshakeReq)
	require.NoError(t, err)

	handshakeHTTPReq := httptest.NewRequest("POST", "/handshake", bytes.NewBuffer(handshakeBytes))
	handshakeHTTPReq.Header.Set("Content-Type", "application/json")
	handshakeW := httptest.NewRecorder()

	router.ServeHTTP(handshakeW, handshakeHTTPReq)

	assert.Equal(t, http.StatusOK, handshakeW.Code)

	var handshakeResp models.HandshakeResponse
	err = json.Unmarshal(handshakeW.Body.Bytes(), &handshakeResp)
	require.NoError(t, err)

	// Now send heartbeat with the token
	heartbeatReq := models.HeartbeatRequest{
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

	heartbeatBytes, err := json.Marshal(heartbeatReq)
	require.NoError(t, err)

	heartbeatHTTPReq := httptest.NewRequest("POST", "/heartbeat", bytes.NewBuffer(heartbeatBytes))
	heartbeatHTTPReq.Header.Set("Content-Type", "application/json")
	heartbeatHTTPReq.Header.Set("Authorization", "Bearer "+handshakeResp.Token)
	heartbeatW := httptest.NewRecorder()

	router.ServeHTTP(heartbeatW, heartbeatHTTPReq)

	// Assert response
	assert.Equal(t, http.StatusOK, heartbeatW.Code)

	var heartbeatResp models.HeartbeatResponse
	err = json.Unmarshal(heartbeatW.Body.Bytes(), &heartbeatResp)
	require.NoError(t, err)

	assert.True(t, heartbeatResp.Success)
	assert.Equal(t, "Heartbeat received", heartbeatResp.Message)
}

func TestControllerHandler_Heartbeat_MissingAuthHeader(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.POST("/heartbeat", handler.Heartbeat)

	// Create request without authorization header
	reqBody := models.HeartbeatRequest{
		Status:  "active",
		Message: "Controller is running",
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/heartbeat", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Missing authorization header", response["message"].(string))
}

func TestControllerHandler_Heartbeat_InvalidAuthHeader(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.POST("/heartbeat", handler.Heartbeat)

	// Create request with invalid authorization header
	reqBody := models.HeartbeatRequest{
		Status:  "active",
		Message: "Controller is running",
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/heartbeat", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat token123") // Wrong format

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Invalid authorization header format", response["message"].(string))
}

func TestControllerHandler_Heartbeat_InvalidToken(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.POST("/heartbeat", handler.Heartbeat)

	// Create request with invalid token
	reqBody := models.HeartbeatRequest{
		Status:  "active",
		Message: "Controller is running",
	}

	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/heartbeat", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Invalid or expired token", response["message"].(string))
}

func TestControllerHandler_GetControllerStatus_Success(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup routes
	router.POST("/handshake", handler.Handshake)
	router.GET("/controllers/:id", handler.GetControllerStatus)

	// First, create a controller via handshake
	handshakeReq := models.HandshakeRequest{
		ClusterID:   "test-cluster-1",
		ClusterName: "Test Cluster",
		Version:     "1.0.0",
		Nonce:       "test-nonce-123",
	}

	handshakeBytes, err := json.Marshal(handshakeReq)
	require.NoError(t, err)

	handshakeHTTPReq := httptest.NewRequest("POST", "/handshake", bytes.NewBuffer(handshakeBytes))
	handshakeHTTPReq.Header.Set("Content-Type", "application/json")
	handshakeW := httptest.NewRecorder()

	router.ServeHTTP(handshakeW, handshakeHTTPReq)

	var handshakeResp models.HandshakeResponse
	err = json.Unmarshal(handshakeW.Body.Bytes(), &handshakeResp)
	require.NoError(t, err)

	// Now get the controller status
	req := httptest.NewRequest("GET", "/controllers/"+handshakeResp.ControllerID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	controllerData := response["controller"].(map[string]interface{})
	assert.Equal(t, handshakeResp.ControllerID, controllerData["id"])
	assert.Equal(t, "test-cluster-1", controllerData["cluster_id"])
	assert.Equal(t, "Test Cluster", controllerData["cluster_name"])
	assert.Equal(t, "1.0.0", controllerData["version"])
	assert.Equal(t, "active", controllerData["status"])
	assert.True(t, controllerData["is_online"].(bool))
}

func TestControllerHandler_GetControllerStatus_NotFound(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.GET("/controllers/:id", handler.GetControllerStatus)

	// Request non-existent controller
	req := httptest.NewRequest("GET", "/controllers/non-existent-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Controller not found", response["message"].(string))
}

func TestControllerHandler_GetControllerStatus_MissingID(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup route
	router.GET("/controllers/:id", handler.GetControllerStatus)

	// Request with empty ID parameter - this will be handled by Gin as a 404
	req := httptest.NewRequest("GET", "/controllers/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Gin will return 404 for this case since the route doesn't match
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestControllerHandler_GetAllControllers_Success(t *testing.T) {
	handler, router := setupControllerHandlerTest(t)

	// Setup routes
	router.POST("/handshake", handler.Handshake)
	router.GET("/controllers", handler.GetAllControllers)

	// Create multiple controllers via handshake
	clusters := []models.HandshakeRequest{
		{
			ClusterID:   "cluster-1",
			ClusterName: "Cluster 1",
			Version:     "1.0.0",
			Nonce:       "nonce-1",
		},
		{
			ClusterID:   "cluster-2",
			ClusterName: "Cluster 2",
			Version:     "1.1.0",
			Nonce:       "nonce-2",
		},
	}

	for _, cluster := range clusters {
		handshakeBytes, err := json.Marshal(cluster)
		require.NoError(t, err)

		handshakeHTTPReq := httptest.NewRequest("POST", "/handshake", bytes.NewBuffer(handshakeBytes))
		handshakeHTTPReq.Header.Set("Content-Type", "application/json")
		handshakeW := httptest.NewRecorder()

		router.ServeHTTP(handshakeW, handshakeHTTPReq)
		assert.Equal(t, http.StatusOK, handshakeW.Code)
	}

	// Now get all controllers
	req := httptest.NewRequest("GET", "/controllers", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	controllers := response["controllers"].([]interface{})
	assert.Len(t, controllers, 2)

	// Verify both controllers are present
	clusterIDs := make(map[string]bool)
	for _, controller := range controllers {
		controllerData := controller.(map[string]interface{})
		clusterIDs[controllerData["cluster_id"].(string)] = true
		assert.True(t, controllerData["is_online"].(bool))
	}

	assert.True(t, clusterIDs["cluster-1"])
	assert.True(t, clusterIDs["cluster-2"])
}
