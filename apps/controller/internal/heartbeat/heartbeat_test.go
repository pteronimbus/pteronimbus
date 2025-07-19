package heartbeat

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockBackendClient is a mock implementation of the backend client for testing
type MockBackendClient struct {
	handshakeCalled bool
	heartbeatCalled bool
	heartbeatCount  int
	shouldFail      bool
	lastStatus      string
	lastMessage     string
}

func (m *MockBackendClient) Handshake(ctx context.Context) error {
	m.handshakeCalled = true
	if m.shouldFail {
		return assert.AnError
	}
	return nil
}

func (m *MockBackendClient) Heartbeat(ctx context.Context, status, message string) error {
	m.heartbeatCalled = true
	m.heartbeatCount++
	m.lastStatus = status
	m.lastMessage = message
	if m.shouldFail {
		return assert.AnError
	}
	return nil
}

func (m *MockBackendClient) GetControllerID() string {
	return "test-controller-id"
}

func (m *MockBackendClient) GetAuthToken() string {
	return "test-token"
}

func (m *MockBackendClient) GetHeartbeatTTL() int {
	return 5
}

func TestManager_NewManager(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Second * 5

	manager := NewManager(mockClient, interval)

	assert.Equal(t, mockClient, manager.client)
	assert.Equal(t, interval, manager.interval)
	assert.NotNil(t, manager.stopChan)
	assert.False(t, manager.isRunning)
}

func TestManager_Start_Success(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Millisecond * 100 // Short interval for testing
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	err := manager.Start(ctx)
	require.NoError(t, err)

	assert.True(t, manager.isRunning)
	assert.False(t, manager.startTime.IsZero())

	// Wait a bit for the heartbeat to be sent
	time.Sleep(time.Millisecond * 150)

	// Stop the manager
	manager.Stop()

	// Verify that heartbeat was called
	assert.True(t, mockClient.heartbeatCalled)
	assert.Equal(t, "active", mockClient.lastStatus)
	assert.Equal(t, "Controller is running", mockClient.lastMessage)
}

func TestManager_Start_AlreadyRunning(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Second * 5
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	// Start first time
	err := manager.Start(ctx)
	require.NoError(t, err)
	assert.True(t, manager.isRunning)

	// Try to start again
	err = manager.Start(ctx)
	require.NoError(t, err) // Should not return error, just ignore

	manager.Stop()
}

func TestManager_Stop_NotRunning(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Second * 5
	manager := NewManager(mockClient, interval)

	// Stop when not running should not panic
	manager.Stop()
	assert.False(t, manager.isRunning)
}

func TestManager_Stop_Running(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Millisecond * 100
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)
	assert.True(t, manager.isRunning)

	// Stop the manager
	manager.Stop()
	assert.False(t, manager.isRunning)
}

func TestManager_IsRunning(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Second * 5
	manager := NewManager(mockClient, interval)

	// Initially not running
	assert.False(t, manager.IsRunning())

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)
	assert.True(t, manager.IsRunning())

	// Stop the manager
	manager.Stop()
	assert.False(t, manager.IsRunning())
}

func TestManager_GetUptime(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Second * 5
	manager := NewManager(mockClient, interval)

	// Initially zero uptime (before starting)
	assert.Equal(t, time.Duration(0), manager.GetUptime())

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(time.Millisecond * 50)

	// Check uptime
	uptime := manager.GetUptime()
	assert.True(t, uptime > 0)
	assert.True(t, uptime < time.Second) // Should be less than 1 second

	manager.Stop()
}

func TestManager_HeartbeatLoop_ContextCancelled(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Millisecond * 100
	manager := NewManager(mockClient, interval)

	ctx, cancel := context.WithCancel(context.Background())

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)

	// Wait a bit for initial heartbeat
	time.Sleep(time.Millisecond * 50)

	// Cancel context
	cancel()

	// Wait a bit for the loop to stop
	time.Sleep(time.Millisecond * 100)

	// Verify that heartbeat was called at least once
	assert.True(t, mockClient.heartbeatCalled)
	assert.True(t, mockClient.heartbeatCount >= 1)
}

func TestManager_HeartbeatLoop_StopChannel(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Millisecond * 100
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)

	// Wait a bit for initial heartbeat
	time.Sleep(time.Millisecond * 50)

	// Stop via stop channel
	manager.Stop()

	// Wait a bit for the loop to stop
	time.Sleep(time.Millisecond * 100)

	// Verify that heartbeat was called at least once
	assert.True(t, mockClient.heartbeatCalled)
	assert.True(t, mockClient.heartbeatCount >= 1)
}

func TestManager_HeartbeatLoop_ClientFailure(t *testing.T) {
	mockClient := &MockBackendClient{shouldFail: true}
	interval := time.Millisecond * 100
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)

	// Wait a bit for heartbeat attempts
	time.Sleep(time.Millisecond * 150)

	// Stop the manager
	manager.Stop()

	// Verify that heartbeat was attempted
	assert.True(t, mockClient.heartbeatCalled)
	assert.True(t, mockClient.heartbeatCount >= 1)
}

func TestManager_HeartbeatLoop_MultipleHeartbeats(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Millisecond * 50 // Very short interval for testing
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)

	// Wait for multiple heartbeats
	time.Sleep(time.Millisecond * 200)

	// Stop the manager
	manager.Stop()

	// Verify that multiple heartbeats were sent
	assert.True(t, mockClient.heartbeatCalled)
	assert.True(t, mockClient.heartbeatCount >= 3) // Should have sent at least 3 heartbeats
}

func TestManager_HeartbeatLoop_InitialHeartbeatImmediate(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Second * 5 // Long interval
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)

	// Wait a short time - should have sent initial heartbeat immediately
	time.Sleep(time.Millisecond * 50)

	// Stop the manager
	manager.Stop()

	// Verify that initial heartbeat was sent immediately
	assert.True(t, mockClient.heartbeatCalled)
	assert.Equal(t, 1, mockClient.heartbeatCount)
}

func TestManager_HeartbeatLoop_StatusAndMessage(t *testing.T) {
	mockClient := &MockBackendClient{}
	interval := time.Millisecond * 100
	manager := NewManager(mockClient, interval)

	ctx := context.Background()

	// Start the manager
	err := manager.Start(ctx)
	require.NoError(t, err)

	// Wait for heartbeat
	time.Sleep(time.Millisecond * 150)

	// Stop the manager
	manager.Stop()

	// Verify the status and message
	assert.Equal(t, "active", mockClient.lastStatus)
	assert.Equal(t, "Controller is running", mockClient.lastMessage)
}
