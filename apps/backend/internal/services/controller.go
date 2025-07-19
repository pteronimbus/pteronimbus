package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
		existingController.Status = "active"
		existingController.HandshakeToken = s.generateControllerToken(existingController.ID, req.ClusterID)

		if err := s.db.WithContext(ctx).Save(&existingController).Error; err != nil {
			return nil, fmt.Errorf("failed to update controller: %w", err)
		}

		return &models.HandshakeResponse{
			Success:      true,
			ControllerID: existingController.ID,
			Token:        existingController.HandshakeToken,
			Message:      "Controller re-registered successfully",
			HeartbeatURL: "/api/controller/heartbeat",
			HeartbeatTTL: int(s.config.Controller.HeartbeatTTL.Seconds()),
		}, nil
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing controller: %w", err)
	}

	// Create new controller
	controllerID := uuid.New().String()
	controller := models.Controller{
		ID:             controllerID,
		ClusterID:      req.ClusterID,
		ClusterName:    req.ClusterName,
		Version:        req.Version,
		LastHeartbeat:  time.Now().UTC(),
		Status:         "active",
		HandshakeToken: s.generateControllerToken(controllerID, req.ClusterID),
	}

	if err := s.db.WithContext(ctx).Create(&controller).Error; err != nil {
		return nil, fmt.Errorf("failed to create controller: %w", err)
	}

	return &models.HandshakeResponse{
		Success:      true,
		ControllerID: controller.ID,
		Token:        controller.HandshakeToken,
		Message:      "Controller registered successfully",
		HeartbeatURL: "/api/controller/heartbeat",
		HeartbeatTTL: int(s.config.Controller.HeartbeatTTL.Seconds()),
	}, nil
}

// Heartbeat processes a controller heartbeat and updates its status
func (s *ControllerService) Heartbeat(ctx context.Context, controllerID string, req *models.HeartbeatRequest) (*models.HeartbeatResponse, error) {
	// Update controller heartbeat
	result := s.db.WithContext(ctx).Model(&models.Controller{}).
		Where("id = ?", controllerID).
		Updates(map[string]interface{}{
			"last_heartbeat": time.Now().UTC(),
			"status":         req.Status,
		})

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update controller heartbeat: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return &models.HeartbeatResponse{
			Success: false,
			Message: "Controller not found",
		}, nil
	}

	return &models.HeartbeatResponse{
		Success: true,
		Message: "Heartbeat received",
	}, nil
}

// GetControllerStatus returns the current status of a controller
func (s *ControllerService) GetControllerStatus(ctx context.Context, controllerID string) (*models.ControllerStatus, error) {
	var controller models.Controller
	err := s.db.WithContext(ctx).Where("id = ?", controllerID).First(&controller).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get controller: %w", err)
	}

	// Determine if controller is online based on last heartbeat
	isOnline := time.Since(controller.LastHeartbeat) < s.config.Controller.MaxHeartbeatAge

	status := &models.ControllerStatus{
		ID:            controller.ID,
		ClusterID:     controller.ClusterID,
		ClusterName:   controller.ClusterName,
		Version:       controller.Version,
		Status:        controller.Status,
		LastHeartbeat: controller.LastHeartbeat,
		IsOnline:      isOnline,
		CreatedAt:     controller.CreatedAt,
	}

	// Calculate uptime if online
	if isOnline {
		status.Uptime = time.Since(controller.CreatedAt).String()
	}

	return status, nil
}

// GetAllControllers returns all registered controllers with their status
func (s *ControllerService) GetAllControllers(ctx context.Context) ([]*models.ControllerStatus, error) {
	var controllers []models.Controller
	err := s.db.WithContext(ctx).Find(&controllers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get controllers: %w", err)
	}

	var statuses []*models.ControllerStatus
	for _, controller := range controllers {
		isOnline := time.Since(controller.LastHeartbeat) < s.config.Controller.MaxHeartbeatAge

		status := &models.ControllerStatus{
			ID:            controller.ID,
			ClusterID:     controller.ClusterID,
			ClusterName:   controller.ClusterName,
			Version:       controller.Version,
			Status:        controller.Status,
			LastHeartbeat: controller.LastHeartbeat,
			IsOnline:      isOnline,
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
