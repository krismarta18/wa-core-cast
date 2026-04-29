package database

import (
	"embed"
	"fmt"
	"strings"
	"time"

	"wacast/core/utils"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var EmbeddedMigrations embed.FS

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

// LoadMigrationsFromEmbedded loads migration files from embedded FS
func (mr *MigrationRunner) LoadMigrationsFromEmbedded() error {
	files, err := EmbeddedMigrations.ReadDir("migrations")
	if err != nil {
		utils.Error("Failed to read embedded migrations", zap.Error(err))
		return fmt.Errorf("failed to read embedded migrations: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			// Read file content
			content, err := EmbeddedMigrations.ReadFile("migrations/" + file.Name())
			if err != nil {
				utils.Error("Failed to read embedded migration file", zap.Error(err), zap.String("file", file.Name()))
				continue
			}

			// Parse migration version from filename
			parts := strings.Split(file.Name(), "_")
			if len(parts) < 2 {
				continue
			}

			version := parts[0]
			name := strings.TrimSuffix(file.Name(), ".sql")

			migration := Migration{
				Version: version,
				Name:    name,
				UpSQL:   string(content),
			}

			mr.migrations = append(mr.migrations, migration)
		}
	}

	return nil
}

// InitMigrationTable creates migration tracking table if not exists
func (mr *MigrationRunner) InitMigrationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			applied_at DATETIME NOT NULL
		)
	`
	if mr.db.DriverType() == "postgres" {
		query = `
			CREATE TABLE IF NOT EXISTS migrations (
				id SERIAL PRIMARY KEY,
				version VARCHAR(255) UNIQUE NOT NULL,
				name VARCHAR(255) NOT NULL,
				applied_at TIMESTAMPTZ NOT NULL
			)
		`
	}

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

		// Execute migration - Split by semicolon to handle multiple statements
		// This is necessary because many drivers only execute the first statement in Exec()
		sql := migration.UpSQL
		if mr.db.DriverType() == "sqlite" {
			sql = translateToSQLite(sql)
		}

		statements := strings.Split(sql, ";")
		for _, stmt := range statements {
			trimmedStmt := strings.TrimSpace(stmt)
			if trimmedStmt == "" {
				continue
			}
			
			_, err := mr.db.GetConnection().Exec(trimmedStmt)
			if err != nil {
				utils.Error(fmt.Sprintf("Migration Statement FAILED! Version: %s, Name: %s, Error: %v", 
					migration.Version, migration.Name, err), zap.String("sql", trimmedStmt))
				return fmt.Errorf("failed to apply migration statement in %s: %w", migration.Version, err)
			}
		}

		// Record migration as applied
		recordQuery := `INSERT INTO migrations (version, name, applied_at) VALUES ($1, $2, $3)`
		if mr.db.DriverType() == "sqlite" {
			recordQuery = `INSERT INTO migrations (version, name, applied_at) VALUES (?, ?, ?)`
		}
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

func translateToSQLite(sql string) string {
	var inAlterBlock bool
	var currentAlterTable string
	
	// 1. Remove Postgres specific schema prefixes and quotes
	sql = strings.ReplaceAll(sql, "\"public\".", "")
	sql = strings.ReplaceAll(sql, "public.", "")
	
	// 2. Remove Postgres specific DROP TYPE, CREATE TYPE, and EXTENSION (Multi-line)
	lines := strings.Split(sql, "\n")
	var filteredLines []string
	inTypeBlock := false
	inDoBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Start of a type block, DO block, SEQUENCE block, FUNCTION block, or whatsmeow table
		if strings.HasPrefix(trimmed, "DROP TYPE") || strings.HasPrefix(trimmed, "CREATE EXTENSION") || 
		   strings.HasPrefix(trimmed, "DROP SEQUENCE") || strings.HasPrefix(trimmed, "DROP FUNCTION") ||
		   strings.HasPrefix(trimmed, "ALTER SEQUENCE") || strings.HasPrefix(trimmed, "OWNED BY") ||
		   strings.HasPrefix(trimmed, "SELECT setval") || strings.Contains(trimmed, "whatsmeow_") {
			
			// If it's a CREATE TABLE or CREATE INDEX for whatsmeow, we need to enter a block to skip it all
			if strings.HasPrefix(trimmed, "CREATE TABLE") || strings.HasPrefix(trimmed, "CREATE INDEX") ||
			   strings.HasPrefix(trimmed, "ALTER TABLE") {
				inTypeBlock = true
				continue
			}
			
			// For simple DROP TABLE or other single line whatsmeow commands
			continue
		}
		if strings.HasPrefix(trimmed, "CREATE TYPE") || strings.HasPrefix(trimmed, "DO $$") || 
		   strings.HasPrefix(trimmed, "CREATE SEQUENCE") || strings.HasPrefix(trimmed, "CREATE FUNCTION") ||
		   strings.HasPrefix(trimmed, "CREATE OR REPLACE FUNCTION") {
			inTypeBlock = true
			if strings.Contains(trimmed, "$$") {
				inDoBlock = true
			}
			continue
		}
		
		// End of a block
		if inTypeBlock {
			// If it's a DO block or FUNCTION using $$, wait for $$;
			if inDoBlock {
				if strings.Contains(trimmed, "$$;") {
					inTypeBlock = false
					inDoBlock = false
				}
			} else {
				// Simple blocks end with ;
				if strings.HasSuffix(trimmed, ";") {
					inTypeBlock = false
				}
			}
			continue
		}
		
		// Remove CASCADE from DROP TABLE
		if strings.HasPrefix(trimmed, "DROP TABLE") && strings.Contains(trimmed, "CASCADE") {
			line = strings.ReplaceAll(line, "CASCADE", "")
		}

		// Handle ALTER TABLE specific Postgres syntax
		if strings.Contains(trimmed, "ALTER TABLE") || inAlterBlock {
			// SQLite doesn't support ALTER COLUMN, ADD CONSTRAINT, ADD PRIMARY KEY, or ADD FOREIGN KEY
			if strings.Contains(trimmed, "ALTER COLUMN") || strings.Contains(trimmed, "ADD CONSTRAINT") ||
			   strings.Contains(trimmed, "PRIMARY KEY") || strings.Contains(trimmed, "FOREIGN KEY") {
				continue 
			}
			
			// Remove IF NOT EXISTS/IF EXISTS for ALTER TABLE ONLY
			line = strings.ReplaceAll(line, "IF NOT EXISTS", "")
			line = strings.ReplaceAll(line, "IF EXISTS", "")
			line = strings.ReplaceAll(line, "public.", "") 

			// Detect start of multi-line ALTER
			if strings.Contains(trimmed, "ALTER TABLE") {
				fields := strings.Fields(trimmed)
				if len(fields) >= 3 {
					currentAlterTable = fields[2]
					// Remove quotes if present
					currentAlterTable = strings.Trim(currentAlterTable, "\"")
					inAlterBlock = true
				}
			}

			// If we are in an ALTER block and it's an ADD COLUMN
			if inAlterBlock && strings.Contains(trimmed, "ADD COLUMN") {
				// If the line already contains "ALTER TABLE", just fix the "public." and types
				if strings.Contains(trimmed, "ALTER TABLE") {
					line = strings.ReplaceAll(line, "public.", "")
					// Types will be replaced later by the general logic
				} else {
					// It's a multi-line ADD COLUMN, prepend ALTER TABLE
					p := strings.TrimSuffix(strings.TrimSpace(trimmed), ",")
					line = fmt.Sprintf("ALTER TABLE \"%s\" %s;", currentAlterTable, p)
				}
			}

			// End of ALTER block
			if strings.HasSuffix(trimmed, ";") {
				inAlterBlock = false
			}
		}

		// Remove Navicat/Postgres specific noise
		line = strings.ReplaceAll(line, "COLLATE \"pg_catalog\".\"default\"", "")
		// Remove any "pg_catalog"."..." pattern
		if strings.Contains(line, "\"pg_catalog\"") {
			start := strings.Index(line, "\"pg_catalog\"")
			after := line[start+len("\"pg_catalog\""):]
			if strings.HasPrefix(after, ".") {
				// Find the end of the next quoted part (e.g., ."text_ops")
				// after is like ."text_ops" ASC...
				secondPart := after[1:] // "text_ops" ASC...
				if strings.HasPrefix(secondPart, "\"") {
					endQuote := strings.Index(secondPart[1:], "\"")
					if endQuote != -1 {
						// Remove from start of "pg_catalog" to the end of second quote
						line = line[:start] + secondPart[endQuote+2:]
					} else {
						line = strings.ReplaceAll(line, "\"pg_catalog\"", "")
					}
				} else {
					line = strings.ReplaceAll(line, "\"pg_catalog\"", "")
				}
			} else {
				line = strings.ReplaceAll(line, "\"pg_catalog\"", "")
			}
		}
		line = strings.ReplaceAll(line, "ASC NULLS LAST", "")
		line = strings.ReplaceAll(line, "DESC NULLS FIRST", "")
		line = strings.ReplaceAll(line, "::regclass", "")
		line = strings.ReplaceAll(line, "::character varying", "")
		line = strings.ReplaceAll(line, "::text", "")
		line = strings.ReplaceAll(line, "::TEXT", "")
		line = strings.ReplaceAll(line, "::jsonb", "")
		line = strings.ReplaceAll(line, "::UUID", "")
		line = strings.ReplaceAll(line, "::uuid", "")
		
		// Remove USING btree and its parentheses content if needed, but simplest is to just remove the phrase
		if strings.Contains(line, "USING btree") {
			line = strings.ReplaceAll(line, "USING btree", "")
		}

		// Replace functions FIRST
		line = strings.ReplaceAll(line, "NOW()", "CURRENT_TIMESTAMP")
		line = strings.ReplaceAll(line, "now()", "CURRENT_TIMESTAMP")
		line = strings.ReplaceAll(line, "gen_random_uuid()", "(lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6))))")

		// Replace types with more specific patterns (using spaces or quotes)
		line = strings.ReplaceAll(line, "timestamptz(6)", "DATETIME")
		line = strings.ReplaceAll(line, "timestamptz", "DATETIME")
		line = strings.ReplaceAll(line, " uuid", " TEXT")
		line = strings.ReplaceAll(line, "(uuid", "(TEXT")
		line = strings.ReplaceAll(line, " bool", " BOOLEAN")
		line = strings.ReplaceAll(line, "(bool", "(BOOLEAN")
		line = strings.ReplaceAll(line, " bytea", " BLOB")
		line = strings.ReplaceAll(line, " jsonb", " TEXT")
		line = strings.ReplaceAll(line, " inet", " TEXT")
		// Detect ID columns with nextval and turn them into AUTOINCREMENT
		if strings.Contains(line, "\"id\"") && strings.Contains(line, "DEFAULT nextval") {
			if strings.Contains(line, "int4") {
				line = strings.ReplaceAll(line, "int4", "INTEGER PRIMARY KEY AUTOINCREMENT")
			} else if strings.Contains(line, "int8") {
				line = strings.ReplaceAll(line, "int8", "INTEGER PRIMARY KEY AUTOINCREMENT")
			} else if strings.Contains(line, "INTEGER") {
				line = strings.ReplaceAll(line, "INTEGER", "INTEGER PRIMARY KEY AUTOINCREMENT")
			} else if strings.Contains(line, "int") {
				line = strings.ReplaceAll(line, "int", "INTEGER PRIMARY KEY AUTOINCREMENT")
			}
			// Remove the DEFAULT nextval part
			start := strings.Index(line, "DEFAULT nextval")
			end := strings.Index(line[start:], ")")
			if end != -1 {
				line = line[:start] + line[start+end+1:]
			}
			// Remove NOT NULL if it's already there (Primary Key implies it)
			line = strings.ReplaceAll(line, "NOT NULL", "")
		} else {
			line = strings.ReplaceAll(line, " int4", " INTEGER")
			line = strings.ReplaceAll(line, " int8", " INTEGER")
			line = strings.ReplaceAll(line, " character varying", " varchar")
			
			// Remove DEFAULT nextval('...') for non-ID columns
			if strings.Contains(line, "DEFAULT nextval") {
				start := strings.Index(line, "DEFAULT nextval")
				end := strings.Index(line[start:], ")")
				if end != -1 {
					line = line[:start] + line[start+end+1:]
				}
			}
		}

		line = strings.ReplaceAll(line, "SERIAL PRIMARY KEY", "INTEGER PRIMARY KEY AUTOINCREMENT")

		filteredLines = append(filteredLines, line)
	}
	sql = strings.Join(filteredLines, "\n")

	// 3. Data type and syntax replacements (Case-insensitive)
	replacements := []string{
		"FALSE", "0",
		"false", "0",
		"TRUE", "1",
		"true", "1",
		"ILIKE", "LIKE",
		"ilike", "LIKE",
	}
	
	r := strings.NewReplacer(replacements...)
	translated := r.Replace(sql)
	
	return translated
}
