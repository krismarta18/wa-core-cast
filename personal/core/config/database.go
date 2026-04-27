package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
		Host:              getEnv("DB_HOST", ""),
		Port:              port,
		User:              getEnv("DB_USER", ""),
		Password:          getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "wacast"),
		SSLMode:           getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:      maxOpenConns,
		MaxIdleConns:      maxIdleConns,
		ConnMaxLifetime:   connMaxLifetime,
		ConnMaxIdleTime:   connMaxIdleTime,
		ConnectionTimeout: connectionTimeout,
	}
}

// SaveDatabaseConfigToEnv persists database settings to the .env file
func SaveDatabaseConfigToEnv(cfg *DatabaseConfig) error {
	envFile := ".env"
	input, err := os.ReadFile(envFile)
	if err != nil {
		if os.IsNotExist(err) {
			input = []byte("")
		} else {
			return fmt.Errorf("failed to read .env file: %w", err)
		}
	}

	lines := strings.Split(string(input), "\n")
	envMap := map[string]string{
		"DB_HOST":               cfg.Host,
		"DB_PORT":               strconv.Itoa(cfg.Port),
		"DB_USER":               cfg.User,
		"DB_PASSWORD":           cfg.Password,
		"DB_NAME":               cfg.DBName,
		"DB_SSL_MODE":           cfg.SSLMode,
		"DB_MAX_OPEN_CONNS":      strconv.Itoa(cfg.MaxOpenConns),
		"DB_MAX_IDLE_CONNS":      strconv.Itoa(cfg.MaxIdleConns),
		"DB_CONN_MAX_LIFETIME":  strconv.Itoa(cfg.ConnMaxLifetime),
		"DB_CONN_MAX_IDLE_TIME":  strconv.Itoa(cfg.ConnMaxIdleTime),
		"DB_CONNECTION_TIMEOUT": strconv.Itoa(cfg.ConnectionTimeout),
	}

	for key, value := range envMap {
		found := false
		for i, line := range lines {
			if strings.HasPrefix(line, key+"=") {
				lines[i] = fmt.Sprintf("%s=%s", key, value)
				found = true
				break
			}
		}
		if !found {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(envFile, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
