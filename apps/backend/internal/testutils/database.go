package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDatabaseConfig holds configuration for test database setup
type TestDatabaseConfig struct {
	Image       string
	Password    string
	Database    string
	User        string
	Port        string
	WaitTimeout time.Duration
}

// DefaultTestDatabaseConfig returns a default configuration for test database
func DefaultTestDatabaseConfig() *TestDatabaseConfig {
	return &TestDatabaseConfig{
		Image:       "postgres:16",
		Password:    "password",
		Database:    "testdb",
		User:        "postgres",
		Port:        "5432/tcp",
		WaitTimeout: 30 * time.Second,
	}
}

// SetupTestDatabase creates a PostgreSQL test container and returns a GORM DB connection
// along with a cleanup function that should be called when the test is done
func SetupTestDatabase(t *testing.T) (*gorm.DB, func()) {
	return SetupTestDatabaseWithConfig(t, DefaultTestDatabaseConfig())
}

// SetupTestDatabaseWithConfig creates a PostgreSQL test container with custom configuration
func SetupTestDatabaseWithConfig(t *testing.T, config *TestDatabaseConfig) (*gorm.DB, func()) {
	ctx := context.Background()

	// Start the container using the postgres module
	container, err := postgrescontainer.RunContainer(ctx,
		testcontainers.WithImage(config.Image),
		postgrescontainer.WithDatabase(config.Database),
		postgrescontainer.WithUsername(config.User),
		postgrescontainer.WithPassword(config.Password),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(config.WaitTimeout),
		),
	)
	require.NoError(t, err, "Failed to start PostgreSQL container")

	// Get the container host and port
	host, err := container.Host(ctx)
	require.NoError(t, err, "Failed to get container host")

	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err, "Failed to get container port")

	// Create database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), config.User, config.Password, config.Database)

	// Connect to the database
	db, err := gorm.Open(postgresdriver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable SQL logging in tests
	})
	require.NoError(t, err, "Failed to connect to test database")

	// Test the connection
	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get underlying SQL DB")

	err = sqlDB.Ping()
	require.NoError(t, err, "Failed to ping test database")

	// Return cleanup function
	cleanup := func() {
		// Close database connection
		if sqlDB != nil {
			sqlDB.Close()
		}

		// Stop and remove container
		if container != nil {
			container.Terminate(ctx)
		}
	}

	return db, cleanup
}

// SetupTestDatabaseWithModels creates a PostgreSQL test container and auto-migrates the provided models
func SetupTestDatabaseWithModels(t *testing.T, models ...interface{}) (*gorm.DB, func()) {
	db, cleanup := SetupTestDatabase(t)

	// Auto-migrate the models
	err := db.AutoMigrate(models...)
	require.NoError(t, err, "Failed to auto-migrate models")

	return db, cleanup
}

// CleanupTestDatabase is a helper function to clean up test data
func CleanupTestDatabase(t *testing.T, db *gorm.DB, models ...interface{}) {
	for _, model := range models {
		err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error
		require.NoError(t, err, "Failed to cleanup test data for model")
	}
} 