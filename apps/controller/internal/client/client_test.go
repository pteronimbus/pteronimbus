package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*httptest.Server, *BackendClient) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/controller/handshake":
			if r.Method != "POST" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Parse handshake request
			var req HandshakeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Validate required fields
			if req.ClusterID == "" || req.ClusterName == "" || req.Version == "" || req.Nonce == "" {
				http.Error(w, "Missing required fields", http.StatusBadRequest)
				return
			}

			// Return successful handshake response
			resp := HandshakeResponse{
				Success:       true,
				ControllerID:  "test-controller-id",
				Token:         "test-jwt-token",
				Message:       "Controller registered successfully",
				HeartbeatURL:  "/api/controller/heartbeat",
				HeartbeatTTL:  300,
			}

			json.NewEncoder(w).Encode(resp)

		case "/api/controller/heartbeat":
			if r.Method != "POST" {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Check authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			if authHeader != "Bearer test-jwt-token" {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Parse heartbeat request
			var req HeartbeatRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Return successful heartbeat response
			resp := HeartbeatResponse{
				Success: true,
				Message: "Heartbeat received",
			}

			json.NewEncoder(w).Encode(resp)

		default:
			http.NotFound(w, r)
		}
	}))

	// Create client
	client := NewBackendClient(server.URL, "test-cluster", "Test Cluster", "1.0.0")

	return server, client
}

func TestBackendClient_Handshake_Success(t *testing.T) {
	server, client := setupTestServer(t)
	defer server.Close()

	ctx := context.Background()

	err := client.Handshake(ctx)
	require.NoError(t, err)

	// Verify client state
	assert.Equal(t, "test-controller-id", client.GetControllerID())
	assert.Equal(t, "test-jwt-token", client.GetAuthToken())
	assert.Equal(t, "/api/controller/heartbeat", client.heartbeatURL)
	assert.Equal(t, 300, client.GetHeartbeatTTL())
}

func TestBackendClient_Handshake_ServerError(t *testing.T) {
	// Create server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewBackendClient(server.URL, "test-cluster", "Test Cluster", "1.0.0")
	ctx := context.Background()

	err := client.Handshake(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "handshake failed with status: 500")
}

func TestBackendClient_Handshake_InvalidResponse(t *testing.T) {
	// Create server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewBackendClient(server.URL, "test-cluster", "Test Cluster", "1.0.0")
	ctx := context.Background()

	err := client.Handshake(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode handshake response")
}

func TestBackendClient_Handshake_UnsuccessfulResponse(t *testing.T) {
	// Create server that returns unsuccessful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := HandshakeResponse{
			Success: false,
			Message: "Handshake failed",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewBackendClient(server.URL, "test-cluster", "Test Cluster", "1.0.0")
	ctx := context.Background()

	err := client.Handshake(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Handshake failed")
}

func TestBackendClient_Heartbeat_Success(t *testing.T) {
	server, client := setupTestServer(t)
	defer server.Close()

	// First perform handshake
	ctx := context.Background()
	err := client.Handshake(ctx)
	require.NoError(t, err)

	// Now send heartbeat
	err = client.Heartbeat(ctx, "active", "Controller is running")
	require.NoError(t, err)
}

func TestBackendClient_Heartbeat_NotAuthenticated(t *testing.T) {
	server, client := setupTestServer(t)
	defer server.Close()

	// Try to send heartbeat without handshake
	ctx := context.Background()
	err := client.Heartbeat(ctx, "active", "Controller is running")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestBackendClient_Heartbeat_ServerError(t *testing.T) {
	server, client := setupTestServer(t)
	defer server.Close()

	// First perform handshake
	ctx := context.Background()
	err := client.Handshake(ctx)
	require.NoError(t, err)

	// Create a new server that returns error for heartbeat
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	// Update client to use error server
	client.baseURL = errorServer.URL

	err = client.Heartbeat(ctx, "active", "Controller is running")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "heartbeat failed with status: 500")
}

func TestBackendClient_Heartbeat_Unauthorized(t *testing.T) {
	server, client := setupTestServer(t)
	defer server.Close()

	// First perform handshake
	ctx := context.Background()
	err := client.Handshake(ctx)
	require.NoError(t, err)

	// Create a new server that returns unauthorized for heartbeat
	unauthServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer unauthServer.Close()

	// Update client to use unauthorized server
	client.baseURL = unauthServer.URL

	err = client.Heartbeat(ctx, "active", "Controller is running")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "heartbeat failed with status: 401")
}

func TestBackendClient_Heartbeat_InvalidResponse(t *testing.T) {
	server, client := setupTestServer(t)
	defer server.Close()

	// First perform handshake
	ctx := context.Background()
	err := client.Handshake(ctx)
	require.NoError(t, err)

	// Create a new server that returns invalid JSON
	invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer invalidServer.Close()

	// Update client to use invalid server
	client.baseURL = invalidServer.URL

	err = client.Heartbeat(ctx, "active", "Controller is running")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode heartbeat response")
}

func TestBackendClient_Heartbeat_UnsuccessfulResponse(t *testing.T) {
	server, client := setupTestServer(t)
	defer server.Close()

	// First perform handshake
	ctx := context.Background()
	err := client.Handshake(ctx)
	require.NoError(t, err)

	// Create a new server that returns unsuccessful response
	unsuccessfulServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := HeartbeatResponse{
			Success: false,
			Message: "Heartbeat failed",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer unsuccessfulServer.Close()

	// Update client to use unsuccessful server
	client.baseURL = unsuccessfulServer.URL

	err = client.Heartbeat(ctx, "active", "Controller is running")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Heartbeat failed")
}

func TestBackendClient_NewBackendClient(t *testing.T) {
	client := NewBackendClient("http://localhost:8080", "test-cluster", "Test Cluster", "1.0.0")

	assert.Equal(t, "http://localhost:8080", client.baseURL)
	assert.Equal(t, "test-cluster", client.clusterID)
	assert.Equal(t, "Test Cluster", client.clusterName)
	assert.Equal(t, "1.0.0", client.version)
	assert.NotNil(t, client.httpClient)
	assert.Equal(t, 30*time.Second, client.httpClient.Timeout)
}

func TestBackendClient_GetControllerID_NotSet(t *testing.T) {
	client := NewBackendClient("http://localhost:8080", "test-cluster", "Test Cluster", "1.0.0")
	assert.Equal(t, "", client.GetControllerID())
}

func TestBackendClient_GetAuthToken_NotSet(t *testing.T) {
	client := NewBackendClient("http://localhost:8080", "test-cluster", "Test Cluster", "1.0.0")
	assert.Equal(t, "", client.GetAuthToken())
}

func TestBackendClient_GetHeartbeatTTL_NotSet(t *testing.T) {
	client := NewBackendClient("http://localhost:8080", "test-cluster", "Test Cluster", "1.0.0")
	assert.Equal(t, 0, client.GetHeartbeatTTL())
} 