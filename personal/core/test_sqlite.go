//go:build ignore
package main

import (
	"fmt"
	"wacast/core/config"
	"wacast/core/database"
)

func main() {
	fmt.Println("=== Testing SQLite Connection ===")
	
	cfg := &config.DatabaseConfig{
		Host: "", // Empty host triggers SQLite in this codebase
	}

	db, err := database.InitDatabase(cfg)
	if err != nil {
		fmt.Printf("ERROR: Connection failed: %v\n", err)
		return
	}
	defer db.Close()

	fmt.Println("SUCCESS: Connected to SQLite database.")
	
	// Test basic query
	var version string
	err = db.GetConnection().QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		fmt.Printf("ERROR: Query failed: %v\n", err)
		return
	}
	fmt.Printf("SQLite Engine Version: %s\n", version)
	fmt.Println("================================")
}
