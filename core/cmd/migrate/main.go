package main

import (
	"flag"
	"fmt"
	"log"

	"wacast/core/config"
	"wacast/core/database"
	"wacast/core/utils"
)

func main() {
	// Parse command line arguments
	statusCmd := flag.Bool("status", false, "Show migration status")
	helpCmd := flag.Bool("help", false, "Show help")
	flag.Parse()

	if *helpCmd {
		fmt.Println(`
Migration Status Tool

Usage:
  go run migrate.go -status         Show applied migrations
  go run migrate.go -help           Show this help

Description:
  This tool connects to the database and shows which migrations have been applied.
`)
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	err = utils.InitLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Connect to database
	db, err := database.InitDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("\n=== WACAST Migration Status ===\n")

	// Show status
	if *statusCmd || !*helpCmd {
		runner := database.NewMigrationRunner(db)

		// Initialize migration table
		err := runner.InitMigrationTable()
		if err != nil {
			log.Fatalf("Failed to initialize migration table: %v", err)
		}

		// Get applied migrations
		applied, err := runner.GetAppliedMigrations()
		if err != nil {
			log.Fatalf("Failed to get migrations: %v", err)
		}

		if len(applied) == 0 {
			fmt.Println("No migrations applied yet.")
		} else {
			fmt.Printf("Applied Migrations: %d\n\n", len(applied))
			for version, appliedAt := range applied {
				fmt.Printf("  [%s] Applied at: %s\n", version, appliedAt.Format("2006-01-02 15:04:05"))
			}
		}

		// Load pending migrations
		runner.LoadMigrationsFromDirectory("./migrations")
		fmt.Printf("\nAvailable Migrations: %d\n", len(runner.GetMigrations()))

		fmt.Println("\nStatus: OK ✓")
	}

	fmt.Println()
}
