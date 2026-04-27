package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"wacast/core/config"
	"wacast/core/utils"

	_ "github.com/lib/pq"
	sqliteDriver "modernc.org/sqlite"
	"go.uber.org/zap"
)

func init() {
	sql.Register("sqlite3", &sqliteDriver.Driver{})
}

// Database holds the database connection pool
type Database struct {
	conn   *sql.DB
	driver string
}

// Global database instance
var DB *Database

// ExecContext executes a query with translation
func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.conn.ExecContext(ctx, d.translateQuery(query), args...)
}

// QueryContext executes a query with translation
func (d *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.conn.QueryContext(ctx, d.translateQuery(query), args...)
}

// QueryRowContext executes a query with translation
func (d *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.conn.QueryRowContext(ctx, d.translateQuery(query), args...)
}

// InitDatabase initializes database connection pool
func InitDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	driver := "postgres"
	dsn := cfg.GetDSN()

	if cfg.Host == "" {
		// Use SQLite as default for personal/local mode
		driver = "sqlite"
		dsn = "wacast.db"
		utils.Info("DATABASE: Mode Lokal Aktif (SQLite)", zap.String("file", dsn))
	} else {
		utils.Info("DATABASE: Mencoba koneksi PostgreSQL", zap.String("host", cfg.Host))
	}

	// Open database connection
	db, err := sql.Open(driver, dsn)
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

	// Special handling for SQLite
	if driver == "sqlite" {
		_, _ = db.Exec("PRAGMA foreign_keys = ON;")
		utils.Info("DATABASE: Foreign keys enabled for SQLite")
	}

	utils.Info("Database connected successfully",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	return &Database{conn: db, driver: driver}, nil
}

// UpdateConnection hot-swaps the underlying database connection pool
func (d *Database) UpdateConnection(cfg *config.DatabaseConfig) error {
	driver := "postgres"
	dsn := cfg.GetDSN()

	if cfg.Host == "" {
		driver = "sqlite"
		dsn = "wacast.db"
	}

	utils.Info("Updating database connection", zap.String("driver", driver), zap.String("dsn", dsn))

	// Open new connection
	newDb, err := sql.Open(driver, dsn)
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
	d.driver = driver

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

// DriverType returns "sqlite" or "postgres"
func (d *Database) DriverType() string {
	return d.driver
}

// Translate exports the query translation logic for use in transactions
func (d *Database) Translate(query string) string {
	return d.translateQuery(query)
}

// translateQuery converts Postgres-specific SQL to SQLite-compatible SQL
func (d *Database) translateQuery(query string) string {
	if d.driver != "sqlite" {
		return query
	}

	// 1. Replace NOW() with strftime in RFC3339 format to match Go
	translated := strings.ReplaceAll(query, "NOW()", "strftime('%Y-%m-%dT%H:%M:%SZ', 'now')")
	
	// 2. Replace ILIKE with LIKE (SQLite LIKE is case-insensitive for ASCII)
	translated = strings.ReplaceAll(translated, "ILIKE", "LIKE")

	// 3. Replace CURRENT_DATE with date('now')
	translated = strings.ReplaceAll(translated, "CURRENT_DATE", "date('now')")

	// 4. Handle INTERVAL 'X days' -> replaced with simple subtraction or nothing 
	// This is tricky without regex, but we'll do common ones
	translated = strings.ReplaceAll(translated, "INTERVAL '1 day'", "1")
	translated = strings.ReplaceAll(translated, "INTERVAL '7 days'", "7")

	// 5. Strip type casts like ::text, ::jsonb, ::int4, etc.
	// We handle the ones with names first
	casts := []string{"::text", "::jsonb", "::json", "::int4", "::int8", "::integer", "::timestamp", "::timestamptz", "::date", "::uuid", "::regclass"}
	for _, cast := range casts {
		translated = strings.ReplaceAll(translated, cast, "")
	}
	
	return translated
}

// QueryRow executes a query that returns a single row
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.conn.QueryRow(d.translateQuery(query), args...)
}

// Query executes a query that returns multiple rows
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.conn.Query(d.translateQuery(query), args...)
}

// Exec executes a query without returning rows
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.conn.Exec(d.translateQuery(query), args...)
}

// BeginTx starts a new transaction
func (d *Database) BeginTx() (*sql.Tx, error) {
	return d.conn.Begin()
}
