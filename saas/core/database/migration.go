package database

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"wacast/core/utils"
	"go.uber.org/zap"
)

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// MigrationRunner manages database migrations
type MigrationRunner struct {
	db         *Database
	migrations []Migration
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *Database) *MigrationRunner {
	return &MigrationRunner{
		db:         db,
		migrations: []Migration{},
	}
}

// LoadMigrationsFromDirectory loads migration files from a directory
func (mr *MigrationRunner) LoadMigrationsFromDirectory(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		utils.Error("Failed to read migrations directory", zap.Error(err), zap.String("path", dirPath))
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(dirPath, file.Name())

			// Read file content
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				utils.Error("Failed to read migration file", zap.Error(err), zap.String("file", file.Name()))
				continue
			}

			// Parse migration version from filename
			// Expected format: XXX_migration_name.sql
			parts := strings.Split(file.Name(), "_")
			if len(parts) < 2 {
				utils.Warn("Invalid migration filename format", zap.String("file", file.Name()))
				continue
			}

			version := parts[0]
			name := strings.TrimSuffix(file.Name(), ".sql")

			migration := Migration{
				Version: version,
				Name:    name,
				UpSQL:   string(content),
				DownSQL: "", // We'll implement rollback later
			}

			mr.migrations = append(mr.migrations, migration)

			utils.Debug("Loaded migration", zap.String("version", version), zap.String("name", name))
		}
	}

	return nil
}

// InitMigrationTable creates migration tracking table if not exists
func (mr *MigrationRunner) InitMigrationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL
		)
	`

	_, err := mr.db.Exec(query)
	if err != nil {
		utils.Error("Failed to create migrations table", zap.Error(err))
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	utils.Debug("Migrations table initialized")
	return nil
}

// GetAppliedMigrations retrieves all applied migrations
func (mr *MigrationRunner) GetAppliedMigrations() (map[string]time.Time, error) {
	query := `SELECT version, applied_at FROM migrations ORDER BY version`

	rows, err := mr.db.Query(query)
	if err != nil {
		utils.Error("Failed to get applied migrations", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]time.Time)
	for rows.Next() {
		var version string
		var appliedAt time.Time

		err := rows.Scan(&version, &appliedAt)
		if err != nil {
			utils.Error("Failed to scan migration", zap.Error(err))
			continue
		}

		applied[version] = appliedAt
	}

	return applied, nil
}

// RunMigrations runs all pending migrations
func (mr *MigrationRunner) RunMigrations() error {
	utils.Info("Starting database migrations")

	// Initialize migration table
	err := mr.InitMigrationTable()
	if err != nil {
		return err
	}

	// Get already applied migrations
	applied, err := mr.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Sort migrations by version
	sortMigrations(mr.migrations)

	// Run pending migrations
	for _, migration := range mr.migrations {
		if _, exists := applied[migration.Version]; exists {
			utils.Debug("Migration already applied", zap.String("version", migration.Version))
			continue
		}

		utils.Info("Running migration", zap.String("version", migration.Version), zap.String("name", migration.Name))

		// Execute migration
		_, err := mr.db.Exec(migration.UpSQL)
		if err != nil {
			utils.Error("Failed to apply migration", zap.Error(err), zap.String("version", migration.Version))
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		// Record migration as applied
		recordQuery := `INSERT INTO migrations (version, name, applied_at) VALUES ($1, $2, $3)`
		_, err = mr.db.Exec(recordQuery, migration.Version, migration.Name, time.Now())
		if err != nil {
			utils.Error("Failed to record migration", zap.Error(err), zap.String("version", migration.Version))
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		utils.Info("Migration applied successfully", zap.String("version", migration.Version))
	}

	utils.Info("Database migrations completed successfully")
	return nil
}

// Helper function to sort migrations by version
func sortMigrations(migrations []Migration) {
	// Simple bubble sort for migration versions
	// In production, you might want to use a more sophisticated version comparison
	for i := 0; i < len(migrations); i++ {
		for j := i + 1; j < len(migrations); j++ {
			if migrations[i].Version > migrations[j].Version {
				migrations[i], migrations[j] = migrations[j], migrations[i]
			}
		}
	}
}

// PrintMigrationStatus prints migration status
func (mr *MigrationRunner) PrintMigrationStatus() error {
	applied, err := mr.GetAppliedMigrations()
	if err != nil {
		return err
	}

	utils.Info("Applied migrations:", zap.Int("count", len(applied)))
	for version, appliedAt := range applied {
		utils.Info("  - "+version, zap.Time("applied_at", appliedAt))
	}

	return nil
}

// GetMigrations returns loaded migrations
func (mr *MigrationRunner) GetMigrations() []Migration {
	return mr.migrations
}
