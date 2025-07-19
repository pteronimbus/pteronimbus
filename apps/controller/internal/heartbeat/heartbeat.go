package heartbeat

import (
	"context"
	"log"
	"time"

	"github.com/pteronimbus/pteronimbus/apps/controller/internal/client"
)

// Manager handles periodic heartbeat sending
type Manager struct {
	client    client.BackendClientInterface
	interval  time.Duration
	stopChan  chan struct{}
	isRunning bool
	startTime time.Time
}

// NewManager creates a new heartbeat manager
func NewManager(client client.BackendClientInterface, interval time.Duration) *Manager {
	return &Manager{
		client:   client,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start begins the heartbeat loop
func (m *Manager) Start(ctx context.Context) error {
	if m.isRunning {
		return nil
	}

	m.startTime = time.Now()
	m.isRunning = true

	log.Printf("Starting heartbeat manager with interval: %v", m.interval)

	go m.heartbeatLoop(ctx)

	return nil
}

// Stop stops the heartbeat loop
func (m *Manager) Stop() {
	if !m.isRunning {
		return
	}

	log.Println("Stopping heartbeat manager...")
	close(m.stopChan)
	m.isRunning = false
}

// heartbeatLoop runs the main heartbeat loop
func (m *Manager) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Send initial heartbeat immediately
	if err := m.sendHeartbeat(ctx); err != nil {
		log.Printf("Initial heartbeat failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := m.sendHeartbeat(ctx); err != nil {
				log.Printf("Heartbeat failed: %v", err)
			}
		case <-m.stopChan:
			log.Println("Heartbeat manager stopped")
			return
		case <-ctx.Done():
			log.Println("Heartbeat manager context cancelled")
			return
		}
	}
}

// sendHeartbeat sends a single heartbeat to the backend
func (m *Manager) sendHeartbeat(ctx context.Context) error {
	uptime := time.Since(m.startTime)
	status := "active"
	message := "Controller is running"

	log.Printf("Sending heartbeat - uptime: %v", uptime)

	return m.client.Heartbeat(ctx, status, message)
}

// IsRunning returns whether the heartbeat manager is currently running
func (m *Manager) IsRunning() bool {
	return m.isRunning
}

// GetUptime returns the current uptime
func (m *Manager) GetUptime() time.Duration {
	if m.startTime.IsZero() {
		return time.Duration(0)
	}
	return time.Since(m.startTime)
}
