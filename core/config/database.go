package config

import (
	"fmt"
	"os"
	"strconv"
)

// DatabaseConfig holds all database related configuration
type DatabaseConfig struct {
	Host              string
	Port              int
	User              string
	Password          string
	DBName            string
	SSLMode           string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   int // in minutes
	ConnMaxIdleTime   int // in minutes
	ConnectionTimeout int // in seconds
}

// GetDSN returns PostgreSQL connection string
func (dc *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		dc.Host,
		dc.Port,
		dc.User,
		dc.Password,
		dc.DBName,
		dc.SSLMode,
		dc.ConnectionTimeout,
	)
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() *DatabaseConfig {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	connMaxLifetime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "5"))
	connMaxIdleTime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_IDLE_TIME", "2"))
	connectionTimeout, _ := strconv.Atoi(getEnv("DB_CONNECTION_TIMEOUT", "5"))

	return &DatabaseConfig{
		Host:              getEnv("DB_HOST", "localhost"),
		Port:              port,
		User:              getEnv("DB_USER", "postgres"),
		Password:          getEnv("DB_PASSWORD", "123456"),
		DBName:            getEnv("DB_NAME", "wacast"),
		SSLMode:           getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:      maxOpenConns,
		MaxIdleConns:      maxIdleConns,
		ConnMaxLifetime:   connMaxLifetime,
		ConnMaxIdleTime:   connMaxIdleTime,
		ConnectionTimeout: connectionTimeout,
	}
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
