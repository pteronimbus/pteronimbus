package models

import (
	"time"
)

// Controller represents a registered Kubernetes controller
type Controller struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid"`
	ClusterID      string    `json:"cluster_id" gorm:"uniqueIndex;not null"`
	ClusterName    string    `json:"cluster_name" gorm:"not null"`
	Version        string    `json:"version" gorm:"not null"`
	LastHeartbeat  time.Time `json:"last_heartbeat" gorm:"not null"`
	Status         string    `json:"status" gorm:"not null;default:'pending_approval'"` // pending_approval, active, inactive, error, degraded, rejected
	HandshakeToken string    `json:"-" gorm:"not null"`                       // JWT token for secure communication
	ApprovedAt     *time.Time `json:"approved_at,omitempty" gorm:"index"`     // When the controller was approved
	ApprovedBy     *string    `json:"approved_by,omitempty" gorm:"index"`     // User ID who approved the controller
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// HandshakeRequest represents a controller handshake request
type HandshakeRequest struct {
	ClusterID   string `json:"cluster_id" binding:"required"`
	ClusterName string `json:"cluster_name" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Nonce       string `json:"nonce" binding:"required"` // Random nonce for replay protection
}

// HandshakeResponse represents a controller handshake response
type HandshakeResponse struct {
	Success      bool   `json:"success"`
	ControllerID string `json:"controller_id,omitempty"`
	Token        string `json:"token,omitempty"`
	Message      string `json:"message,omitempty"`
	HeartbeatURL string `json:"heartbeat_url,omitempty"`
	HeartbeatTTL int    `json:"heartbeat_ttl,omitempty"` // in seconds
}

// HeartbeatRequest represents a controller heartbeat request
type HeartbeatRequest struct {
	Status    string            `json:"status" binding:"required"` // active, error, degraded
	Message   string            `json:"message,omitempty"`
	Metrics   map[string]string `json:"metrics,omitempty"`
	Resources map[string]int64  `json:"resources,omitempty"`
}

// HeartbeatResponse represents a controller heartbeat response
type HeartbeatResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ControllerStatus represents the current status of a controller
type ControllerStatus struct {
	ID            string     `json:"id"`
	ClusterID     string     `json:"cluster_id"`
	ClusterName   string     `json:"cluster_name"`
	Version       string     `json:"version"`
	Status        string     `json:"status"`
	LastHeartbeat time.Time  `json:"last_heartbeat"`
	IsOnline      bool       `json:"is_online"`
	Uptime        string     `json:"uptime,omitempty"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	ApprovedBy    *string    `json:"approved_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ControllerApprovalRequest represents a request to approve or reject a controller
type ControllerApprovalRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"` // approve or reject
	Reason string `json:"reason,omitempty"`                               // Optional reason for rejection
}

// ControllerApprovalResponse represents the response to a controller approval action
type ControllerApprovalResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// AdminStats represents admin-level statistics
type AdminStats struct {
	TotalTenants      int64 `json:"total_tenants"`
	TotalUsers        int64 `json:"total_users"`
	TotalGameServers  int64 `json:"total_game_servers"`
	ActiveControllers int64 `json:"active_controllers"`
}
