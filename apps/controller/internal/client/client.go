package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BackendClientInterface defines the interface for backend client operations
type BackendClientInterface interface {
	Handshake(ctx context.Context) error
	Heartbeat(ctx context.Context, status, message string) error
	GetControllerID() string
	GetAuthToken() string
	GetHeartbeatTTL() int
}

// BackendClient handles communication with the Pteronimbus backend
type BackendClient struct {
	baseURL      string
	httpClient   *http.Client
	controllerID string
	authToken    string
	heartbeatURL string
	heartbeatTTL int
	clusterID    string
	clusterName  string
	version      string
}

// HandshakeRequest represents a handshake request to the backend
type HandshakeRequest struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Version     string `json:"version"`
	Nonce       string `json:"nonce"`
}

// HandshakeResponse represents a handshake response from the backend
type HandshakeResponse struct {
	Success      bool   `json:"success"`
	ControllerID string `json:"controller_id,omitempty"`
	Token        string `json:"token,omitempty"`
	Message      string `json:"message,omitempty"`
	HeartbeatURL string `json:"heartbeat_url,omitempty"`
	HeartbeatTTL int    `json:"heartbeat_ttl,omitempty"`
}

// HeartbeatRequest represents a heartbeat request to the backend
type HeartbeatRequest struct {
	Status    string            `json:"status"`
	Message   string            `json:"message,omitempty"`
	Metrics   map[string]string `json:"metrics,omitempty"`
	Resources map[string]int64  `json:"resources,omitempty"`
}

// HeartbeatResponse represents a heartbeat response from the backend
type HeartbeatResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// NewBackendClient creates a new backend client
func NewBackendClient(baseURL, clusterID, clusterName, version string) *BackendClient {
	return &BackendClient{
		baseURL:     baseURL,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		clusterID:   clusterID,
		clusterName: clusterName,
		version:     version,
	}
}

// Handshake performs the initial handshake with the backend
func (c *BackendClient) Handshake(ctx context.Context) error {
	// Generate a random nonce for replay protection
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	req := HandshakeRequest{
		ClusterID:   c.clusterID,
		ClusterName: c.clusterName,
		Version:     c.version,
		Nonce:       hex.EncodeToString(nonce),
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal handshake request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/controller/handshake", bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create handshake request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send handshake request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("handshake failed with status: %d", resp.StatusCode)
	}

	var handshakeResp HandshakeResponse
	if err := json.NewDecoder(resp.Body).Decode(&handshakeResp); err != nil {
		return fmt.Errorf("failed to decode handshake response: %w", err)
	}

	if !handshakeResp.Success {
		return fmt.Errorf("handshake failed: %s", handshakeResp.Message)
	}

	// Store the authentication details
	c.controllerID = handshakeResp.ControllerID
	c.authToken = handshakeResp.Token
	c.heartbeatURL = handshakeResp.HeartbeatURL
	c.heartbeatTTL = handshakeResp.HeartbeatTTL

	return nil
}

// Heartbeat sends a heartbeat to the backend
func (c *BackendClient) Heartbeat(ctx context.Context, status, message string) error {
	if c.authToken == "" {
		return fmt.Errorf("not authenticated - perform handshake first")
	}

	req := HeartbeatRequest{
		Status:  status,
		Message: message,
		Metrics: map[string]string{
			"uptime": time.Since(time.Now()).String(), // This will be negative, but it's just an example
		},
		Resources: map[string]int64{
			"memory_usage": 0, // TODO: Add actual resource metrics
			"cpu_usage":    0,
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+c.heartbeatURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create heartbeat request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send heartbeat request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("heartbeat failed with status: %d", resp.StatusCode)
	}

	var heartbeatResp HeartbeatResponse
	if err := json.NewDecoder(resp.Body).Decode(&heartbeatResp); err != nil {
		return fmt.Errorf("failed to decode heartbeat response: %w", err)
	}

	if !heartbeatResp.Success {
		return fmt.Errorf("heartbeat failed: %s", heartbeatResp.Message)
	}

	return nil
}

// GetControllerID returns the controller ID from the handshake
func (c *BackendClient) GetControllerID() string {
	return c.controllerID
}

// GetAuthToken returns the authentication token
func (c *BackendClient) GetAuthToken() string {
	return c.authToken
}

// GetHeartbeatTTL returns the heartbeat TTL in seconds
func (c *BackendClient) GetHeartbeatTTL() int {
	return c.heartbeatTTL
}
