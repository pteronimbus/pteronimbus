package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
	"gorm.io/gorm"
)

// AdminService handles admin-specific operations and permission checks
type AdminService struct {
	db *gorm.DB
}

// NewAdminService creates a new admin service
func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		db: db,
	}
}

// CheckSuperAdminAccess checks if a user has superadmin access
func (s *AdminService) CheckSuperAdminAccess(ctx context.Context, userID string) (bool, error) {
	// For now, we'll implement a simple superadmin check
	// In the future, this could be based on:
	// 1. A dedicated superadmin table
	// 2. Environment variables for superadmin Discord IDs
	// 3. Special Discord roles
	// 4. Database configuration

	// Get user details
	var user models.User
	err := s.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has superadmin permission in any tenant
	var userTenant models.UserTenant
	err = s.db.WithContext(ctx).
		Where("user_id = ? AND permissions @> ?", userID, gorm.Expr("ARRAY['superadmin']")).
		First(&userTenant).Error

	if err == nil {
		// User has superadmin permission in at least one tenant
		return true, nil
	}

	if err != gorm.ErrRecordNotFound {
		return false, fmt.Errorf("failed to check superadmin permissions: %w", err)
	}

	// Check for hardcoded superadmin Discord IDs (for development/testing)
	superadminDiscordIDs := []string{
		"197918357025062922", // thatguy164.
		// Add more superadmin Discord IDs here as needed
	}

	for _, adminID := range superadminDiscordIDs {
		if strings.EqualFold(user.DiscordUserID, adminID) {
			return true, nil
		}
	}

	return false, nil
}

// GetAdminStats returns admin-level statistics
func (s *AdminService) GetAdminStats(ctx context.Context) (*models.AdminStats, error) {
	var stats models.AdminStats

	// Count total tenants
	err := s.db.WithContext(ctx).Model(&models.Tenant{}).Count(&stats.TotalTenants).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count tenants: %w", err)
	}

	// Count total users
	err = s.db.WithContext(ctx).Model(&models.User{}).Count(&stats.TotalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Count total game servers
	err = s.db.WithContext(ctx).Model(&models.GameServer{}).Count(&stats.TotalGameServers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count game servers: %w", err)
	}

	// Count active controllers
	err = s.db.WithContext(ctx).Model(&models.Controller{}).Where("status = ?", "active").Count(&stats.ActiveControllers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count active controllers: %w", err)
	}

	return &stats, nil
}

// CleanupInactiveControllers removes controllers that haven't sent heartbeats
func (s *AdminService) CleanupInactiveControllers(ctx context.Context) error {
	// This is a simple cleanup - in production, you might want more sophisticated logic
	result := s.db.WithContext(ctx).Where("status = ?", "inactive").Delete(&models.Controller{})
	if result.Error != nil {
		return fmt.Errorf("failed to cleanup inactive controllers: %w", result.Error)
	}

	return nil
}
