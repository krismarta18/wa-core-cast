package database

import (
	"database/sql"
	"fmt"
	"time"

	"wacast/core/config"
	"wacast/core/utils"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Database holds the database connection pool
type Database struct {
	conn *sql.DB
}

// Global database instance
var DB *Database

// InitDatabase initializes database connection pool
func InitDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	utils.Debug("Initializing database connection", zap.String("host", cfg.Host))

	// Open database connection
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		utils.Error("Failed to open database", zap.Error(err))
		return &Database{conn: nil}, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Minute)

	// Test connection
	err = db.Ping()
	if err != nil {
		utils.Error("Failed to ping database", zap.Error(err))
		return &Database{conn: db}, fmt.Errorf("failed to ping database: %w", err)
	}

	utils.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	return &Database{conn: db}, nil
}

// UpdateConnection hot-swaps the underlying database connection pool
func (d *Database) UpdateConnection(cfg *config.DatabaseConfig) error {
	utils.Info("Updating database connection", zap.String("host", cfg.Host))

	// Open new connection
	newDb, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to open new database: %w", err)
	}

	// Test connection
	if err := newDb.Ping(); err != nil {
		newDb.Close()
		return fmt.Errorf("failed to ping new database: %w", err)
	}

	// Configure new pool
	newDb.SetMaxOpenConns(cfg.MaxOpenConns)
	newDb.SetMaxIdleConns(cfg.MaxIdleConns)
	newDb.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)
	newDb.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Minute)

	// Swap connection
	oldDb := d.conn
	d.conn = newDb

	// Close old connection in background (allow active queries to finish)
	if oldDb != nil {
		go func() {
			time.Sleep(5 * time.Second)
			oldDb.Close()
			utils.Debug("Old database connection closed")
		}()
	}

	utils.Info("Database connection updated successfully",
		zap.String("host", cfg.Host),
		zap.String("database", cfg.DBName),
	)

	return nil
}

// GetConnection returns the database connection
func (d *Database) GetConnection() *sql.DB {
	return d.conn
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

// HealthCheck checks if database is healthy
func (d *Database) HealthCheck() bool {
	if err := d.conn.Ping(); err != nil {
		utils.Error("Database health check failed", zap.Error(err))
		return false
	}
	return true
}

// GetStats returns connection pool statistics
func (d *Database) GetStats() sql.DBStats {
	return d.conn.Stats()
}

// Helper functions
// ============================================================================

// QueryRow executes a query that returns a single row
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

// Query executes a query that returns multiple rows
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.conn.Query(query, args...)
}

// Exec executes a query without returning rows
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.conn.Exec(query, args...)
}

// BeginTx starts a new transaction
func (d *Database) BeginTx() (*sql.Tx, error) {
	return d.conn.Begin()
}
