package services

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/pteronimbus/pteronimbus/apps/backend/internal/config"
	"github.com/pteronimbus/pteronimbus/apps/backend/internal/models"
)

// DatabaseService handles database operations
type DatabaseService struct {
	DB *gorm.DB
}

// NewDatabaseService creates a new database service
func NewDatabaseService(cfg *config.DatabaseConfig) (*DatabaseService, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &DatabaseService{DB: db}, nil
}

// AutoMigrate runs database migrations
func (ds *DatabaseService) AutoMigrate() error {
	log.Println("Running database migrations...")
	
	err := ds.DB.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Tenant{},
		&models.UserTenant{},
		&models.TenantDiscordRole{},
		&models.TenantDiscordUser{},
		&models.GameServer{},
		&models.Controller{},
		&models.Permission{},
		&models.Role{},
		&models.SystemRole{},
		&models.UserSystemRole{},
		&models.PermissionAuditLog{},
		&models.GuildMembershipCache{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (ds *DatabaseService) Close() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the GORM database instance
func (ds *DatabaseService) GetDB() *gorm.DB {
	return ds.DB
}