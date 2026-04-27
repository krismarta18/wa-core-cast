package config

import (
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	
	// Server
	ServerHost string
	ServerPort int
	
	// Database
	Database *DatabaseConfig
	
	// Logging
	LogLevel string
	
	// WhatsApp
	SessionTimeout int
	
	// Encryption
	EncryptionKey string

	// JWT
	JWTSecret             string
	JWTExpiryHours        int
	JWTRefreshExpiryHours int

	// Environment
	Environment string
}

// LoadConfig loads all configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	sessionTimeout, _ := strconv.Atoi(getEnv("WHATSAPP_SESSION_TIMEOUT", "300"))
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	jwtRefreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))

	config := &Config{
		ServiceName:    ServiceName,
		ServiceVersion: ServiceVersion,
		ServerHost:     getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:     serverPort,
		Database:       LoadDatabaseConfig(),
		LogLevel:       getEnv("LOG_LEVEL", "debug"),
		SessionTimeout: sessionTimeout,
		EncryptionKey:  getEnv("ENCRYPTION_KEY", "wacast-personal-secret-key-32chars"),
		JWTSecret:             getEnv("JWT_SECRET", "wacast-personal-jwt-secret-very-secure"),
		JWTExpiryHours:        jwtExpiry,
		JWTRefreshExpiryHours: jwtRefreshExpiry,
		Environment:    getEnv("ENVIRONMENT", "development"),
	}

	return config, nil
}

// IsProduction checks if running in production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetServerAddr returns server address
func (c *Config) GetServerAddr() string {
	return c.ServerHost + ":" + strconv.Itoa(c.ServerPort)
}
