package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server     ServerConfig
	Discord    DiscordConfig
	JWT        JWTConfig
	Redis      RedisConfig
	Database   DatabaseConfig
	Controller ControllerConfig
	RBAC       RBACConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	Environment  string
	AllowOrigins []string
}

// DiscordConfig holds Discord OAuth2 configuration
type DiscordConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	BotToken     string
	APIBaseURL   string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ControllerConfig holds controller integration configuration
type ControllerConfig struct {
	HandshakeSecret string
	HeartbeatTTL    time.Duration
	MaxHeartbeatAge time.Duration
}

// RBACConfig holds RBAC system configuration
type RBACConfig struct {
	SuperAdminDiscordID string
	RoleSyncTTL         time.Duration
	GuildCacheTTL       time.Duration
	GracePeriod         time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Host:         getEnv("HOST", "0.0.0.0"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			AllowOrigins: getAllowedOrigins(),
		},
		Discord: DiscordConfig{
			ClientID:     getEnv("DISCORD_CLIENT_ID", ""),
			ClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("DISCORD_REDIRECT_URL", "http://localhost:8080/auth/callback"),
			BotToken:     getEnv("DISCORD_BOT_TOKEN", ""),
			APIBaseURL:   "https://discord.com/api/v10",
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: time.Hour * 24 * 7, // 7 days
			Issuer:          getEnv("JWT_ISSUER", "pteronimbus"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "pteronimbus"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Controller: ControllerConfig{
			HandshakeSecret: getEnv("CONTROLLER_HANDSHAKE_SECRET", ""),
			HeartbeatTTL:    time.Minute * 5,  // 5 minutes
			MaxHeartbeatAge: time.Minute * 10, // 10 minutes
		},
		RBAC: RBACConfig{
			SuperAdminDiscordID: getEnv("SUPER_ADMIN_DISCORD_ID", ""),
			RoleSyncTTL:         time.Minute * 5,  // 5 minutes
			GuildCacheTTL:       time.Minute * 5,  // 5 minutes
			GracePeriod:         time.Minute * 2,  // 2 minutes for security
		},
	}

	return config
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

// getAllowedOrigins returns the list of allowed CORS origins
func getAllowedOrigins() []string {
	// Get the primary frontend URL
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	// For Docker environments, we might need to allow both localhost and container names
	allowedOrigins := []string{frontendURL}

	// Add additional origins if specified via environment variable
	if additionalOrigins := getEnv("ADDITIONAL_CORS_ORIGINS", ""); additionalOrigins != "" {
		// Split by comma and add to allowed origins
		for _, origin := range splitAndTrim(additionalOrigins, ",") {
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	return allowedOrigins
}

// splitAndTrim splits a string by delimiter and trims whitespace
func splitAndTrim(s, delimiter string) []string {
	if s == "" {
		return []string{}
	}

	parts := make([]string, 0)
	for _, part := range strings.Split(s, delimiter) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}
