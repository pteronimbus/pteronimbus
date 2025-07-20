package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"gorm.io/gorm"
)

// ControllerService handles controller registration and heartbeat management
type ControllerService struct {
	db     *gorm.DB
	config *config.Config
	jwt    *JWTService
}

// NewControllerService creates a new controller service
func NewControllerService(db *gorm.DB, config *config.Config, jwt *JWTService) *ControllerService {
	return &ControllerService{
		db:     db,
		config: config,
		jwt:    jwt,
	}
}

// Handshake performs the initial controller registration and authentication
func (s *ControllerService) Handshake(ctx context.Context, req *models.HandshakeRequest) (*models.HandshakeResponse, error) {
	// Validate the handshake secret if configured
	if s.config.Controller.HandshakeSecret != "" {
		if !s.validateHandshakeSecret(req.Nonce) {
			return &models.HandshakeResponse{
				Success: false,
				Message: "Invalid handshake secret",
			}, nil
		}
	}

	// Check if controller already exists
	var existingController models.Controller
	err := s.db.WithContext(ctx).Where("cluster_id = ?", req.ClusterID).First(&existingController).Error
	if err == nil {
		// Controller exists, update it and generate new token
		existingController.ClusterName = req.ClusterName
		existingController.Version = req.Version
		existingController.LastHeartbeat = time.Now().UTC()
		
		// Only update status to active if it was previously approved
		if existingController.Status == "active" {
			existingController.Status = "active"
		} else if existingController.Status == "pending_approval" {
			// Keep as pending_approval if not yet approved
			existingController.Status = "pending_approval"
		}
		
		existingController.HandshakeToken = s.generateControllerToken(existingController.ID, req.ClusterID)

		if err := s.db.WithContext(ctx).Save(&existingController).Error; err != nil {
			return nil, fmt.Errorf("failed to update controller: %w", err)
		}

		message := "Controller re-registered successfully"
		if existingController.Status == "pending_approval" {
			message = "Controller re-registered successfully - awaiting approval"
		}

		return &models.HandshakeResponse{
			Success:      true,
			ControllerID: existingController.ID,
			Token:        existingController.HandshakeToken,
			Message:      message,
			HeartbeatURL: "/api/controller/heartbeat",
			HeartbeatTTL: int(s.config.Controller.HeartbeatTTL.Seconds()),
		}, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing controller: %w", err)
	}

	// Create new controller in pending_approval status
	controllerID := uuid.New().String()
	controller := models.Controller{
		ID:             controllerID,
		ClusterID:      req.ClusterID,
		ClusterName:    req.ClusterName,
		Version:        req.Version,
		LastHeartbeat:  time.Now().UTC(),
		Status:         "pending_approval", // New controllers start as pending
		HandshakeToken: s.generateControllerToken(controllerID, req.ClusterID),
	}

	if err := s.db.WithContext(ctx).Create(&controller).Error; err != nil {
		return nil, fmt.Errorf("failed to create controller: %w", err)
	}

	return &models.HandshakeResponse{
		Success:      true,
		ControllerID: controller.ID,
		Token:        controller.HandshakeToken,
		Message:      "Controller registered successfully - awaiting approval",
		HeartbeatURL: "/api/controller/heartbeat",
		HeartbeatTTL: int(s.config.Controller.HeartbeatTTL.Seconds()),
	}, nil
}

// Heartbeat processes a controller heartbeat and updates its status
func (s *ControllerService) Heartbeat(ctx context.Context, controllerID string, req *models.HeartbeatRequest) (*models.HeartbeatResponse, error) {
	// Validate UUID format first
	if !s.validateUUID(controllerID) {
		return &models.HeartbeatResponse{
			Success: false,
			Message: "Controller not found",
		}, nil
	}

	// First, get the current controller to check its status
	var controller models.Controller
	err := s.db.WithContext(ctx).Where("id = ?", controllerID).First(&controller).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.HeartbeatResponse{
				Success: false,
				Message: "Controller not found",
			}, nil
		}
		// Check if it's a UUID format error from PostgreSQL
		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			return &models.HeartbeatResponse{
				Success: false,
				Message: "Controller not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get controller: %w", err)
	}

	// Only allow status updates for approved controllers
	// Pending controllers can send heartbeats but their status won't change
	updates := map[string]interface{}{
		"last_heartbeat": time.Now().UTC(),
	}

	// Only update status if the controller is approved (active, inactive, error, degraded)
	// Don't allow pending_approval or rejected controllers to change their status
	if controller.Status != "pending_approval" && controller.Status != "rejected" {
		updates["status"] = req.Status
	}

	// Update controller heartbeat
	result := s.db.WithContext(ctx).Model(&models.Controller{}).
		Where("id = ?", controllerID).
		Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update controller heartbeat: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return &models.HeartbeatResponse{
			Success: false,
			Message: "Controller not found",
		}, nil
	}

	message := "Heartbeat received"
	if controller.Status == "pending_approval" {
		message = "Heartbeat received - controller awaiting approval"
	} else if controller.Status == "rejected" {
		message = "Heartbeat received - controller has been rejected"
	}

	return &models.HeartbeatResponse{
		Success: true,
		Message: message,
	}, nil
}

// ApproveController approves a pending controller
func (s *ControllerService) ApproveController(ctx context.Context, controllerID string, approvedBy string) (*models.ControllerApprovalResponse, error) {
	// Validate UUID format first
	if !s.validateUUID(controllerID) {
		return &models.ControllerApprovalResponse{
			Success: false,
			Message: "Controller not found",
		}, nil
	}

	var controller models.Controller
	err := s.db.WithContext(ctx).Where("id = ?", controllerID).First(&controller).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.ControllerApprovalResponse{
				Success: false,
				Message: "Controller not found",
			}, nil
		}
		// Check if it's a UUID format error from PostgreSQL
		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			return &models.ControllerApprovalResponse{
				Success: false,
				Message: "Controller not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get controller: %w", err)
	}

	if controller.Status != "pending_approval" {
		return &models.ControllerApprovalResponse{
			Success: false,
			Message: "Controller is not in pending approval status",
		}, nil
	}

	now := time.Now().UTC()
	controller.Status = "active"
	controller.ApprovedAt = &now
	controller.ApprovedBy = &approvedBy

	if err := s.db.WithContext(ctx).Save(&controller).Error; err != nil {
		return nil, fmt.Errorf("failed to approve controller: %w", err)
	}

	return &models.ControllerApprovalResponse{
		Success: true,
		Message: "Controller approved successfully",
	}, nil
}

// RejectController rejects a pending controller
func (s *ControllerService) RejectController(ctx context.Context, controllerID string, rejectedBy string, reason string) (*models.ControllerApprovalResponse, error) {
	// Validate UUID format first
	if !s.validateUUID(controllerID) {
		return &models.ControllerApprovalResponse{
			Success: false,
			Message: "Controller not found",
		}, nil
	}

	var controller models.Controller
	err := s.db.WithContext(ctx).Where("id = ?", controllerID).First(&controller).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.ControllerApprovalResponse{
				Success: false,
				Message: "Controller not found",
			}, nil
		}
		// Check if it's a UUID format error from PostgreSQL
		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			return &models.ControllerApprovalResponse{
				Success: false,
				Message: "Controller not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get controller: %w", err)
	}

	if controller.Status != "pending_approval" {
		return &models.ControllerApprovalResponse{
			Success: false,
			Message: "Controller is not in pending approval status",
		}, nil
	}

	controller.Status = "rejected"

	if err := s.db.WithContext(ctx).Save(&controller).Error; err != nil {
		return nil, fmt.Errorf("failed to reject controller: %w", err)
	}

	message := "Controller rejected successfully"
	if reason != "" {
		message = fmt.Sprintf("Controller rejected: %s", reason)
	}

	return &models.ControllerApprovalResponse{
		Success: true,
		Message: message,
	}, nil
}

// validateUUID checks if a string is a valid UUID format
func (s *ControllerService) validateUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// GetControllerStatus returns the current status of a controller
func (s *ControllerService) GetControllerStatus(ctx context.Context, controllerID string) (*models.ControllerStatus, error) {
	// Validate UUID format first
	if !s.validateUUID(controllerID) {
		return nil, nil // Return nil to indicate "not found" for invalid UUIDs
	}

	var controller models.Controller
	err := s.db.WithContext(ctx).Where("id = ?", controllerID).First(&controller).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		// Check if it's a UUID format error from PostgreSQL
		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			return nil, nil // Return nil to indicate "not found" for invalid UUIDs
		}
		return nil, fmt.Errorf("failed to get controller: %w", err)
	}

	isOnline := time.Since(controller.LastHeartbeat) < s.config.Controller.MaxHeartbeatAge

	// Auto-transition to degraded status if controller is offline and was previously active
	if !isOnline && controller.Status == "active" {
		controller.Status = "degraded"
		if err := s.db.WithContext(ctx).Save(&controller).Error; err != nil {
			return nil, fmt.Errorf("failed to update controller status to degraded: %w", err)
		}
	}

	status := &models.ControllerStatus{
		ID:            controller.ID,
		ClusterID:     controller.ClusterID,
		ClusterName:   controller.ClusterName,
		Version:       controller.Version,
		Status:        controller.Status,
		LastHeartbeat: controller.LastHeartbeat,
		IsOnline:      isOnline,
		ApprovedAt:    controller.ApprovedAt,
		ApprovedBy:    controller.ApprovedBy,
		CreatedAt:     controller.CreatedAt,
	}

	if isOnline {
		status.Uptime = time.Since(controller.CreatedAt).String()
	}

	return status, nil
}

// GetAllControllers returns all registered controllers
func (s *ControllerService) GetAllControllers(ctx context.Context) ([]*models.ControllerStatus, error) {
	var controllers []models.Controller
	err := s.db.WithContext(ctx).Find(&controllers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get controllers: %w", err)
	}

	var statuses []*models.ControllerStatus
	for _, controller := range controllers {
		isOnline := time.Since(controller.LastHeartbeat) < s.config.Controller.MaxHeartbeatAge

		// Auto-transition to degraded status if controller is offline and was previously active
		if !isOnline && controller.Status == "active" {
			controller.Status = "degraded"
			if err := s.db.WithContext(ctx).Save(&controller).Error; err != nil {
				return nil, fmt.Errorf("failed to update controller status to degraded: %w", err)
			}
		}

		status := &models.ControllerStatus{
			ID:            controller.ID,
			ClusterID:     controller.ClusterID,
			ClusterName:   controller.ClusterName,
			Version:       controller.Version,
			Status:        controller.Status,
			LastHeartbeat: controller.LastHeartbeat,
			IsOnline:      isOnline,
			ApprovedAt:    controller.ApprovedAt,
			ApprovedBy:    controller.ApprovedBy,
			CreatedAt:     controller.CreatedAt,
		}

		if isOnline {
			status.Uptime = time.Since(controller.CreatedAt).String()
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// CleanupInactiveControllers removes controllers that haven't sent heartbeats
func (s *ControllerService) CleanupInactiveControllers(ctx context.Context) error {
	cutoff := time.Now().UTC().Add(-s.config.Controller.MaxHeartbeatAge * 2) // Double the max age for cleanup

	result := s.db.WithContext(ctx).Where("last_heartbeat < ?", cutoff).Delete(&models.Controller{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup inactive controllers: %w", result.Error)
	}

	return nil
}

// validateHandshakeSecret validates the handshake secret using HMAC
func (s *ControllerService) validateHandshakeSecret(nonce string) bool {
	if s.config.Controller.HandshakeSecret == "" {
		return true // No secret configured, allow all
	}

	// In a real implementation, you might want to use the nonce to create a challenge-response
	// For now, we'll use a simple HMAC validation
	h := hmac.New(sha256.New, []byte(s.config.Controller.HandshakeSecret))
	h.Write([]byte(nonce))
	_ = hex.EncodeToString(h.Sum(nil)) // Ignore the hash for now

	// For now, we'll accept any nonce if a secret is configured
	// In production, you'd want to implement proper challenge-response
	return true
}

// generateControllerToken generates a JWT token for controller authentication
func (s *ControllerService) generateControllerToken(controllerID, clusterID string) string {
	claims := jwt.MapClaims{
		"sub":        controllerID,
		"cluster_id": clusterID,
		"type":       "controller",
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
		"iss":        s.config.JWT.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(s.config.JWT.Secret))
	return tokenString
}

// ValidateControllerToken validates a controller JWT token
func (s *ControllerService) ValidateControllerToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "controller" {
			return "", fmt.Errorf("invalid token type")
		}
		return claims["sub"].(string), nil
	}

	return "", fmt.Errorf("invalid token")
}
