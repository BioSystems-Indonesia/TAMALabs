//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

func main() {
	// Get database path
	var dbFileName string
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = os.Getenv("APPDATA")
	}
	if localAppData == "" && runtime.GOOS == "windows" {
		localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
	}

	if localAppData != "" {
		dbFileName = filepath.Join(localAppData, "TAMALabs", "database", "TAMALabs.db")
	} else {
		dbFileName = filepath.Join("data", "TAMALabs.db")
	}

	log.Printf("Using database: %s", dbFileName)

	// Open database
	db, err := sql.Open("sqlite", dbFileName+"?_parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check current schema_migrations state
	var version int
	var dirty bool
	err = db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty)
	if err != nil {
		log.Printf("Error reading schema_migrations: %v", err)
		log.Fatal(err)
	}

	log.Printf("Current version: %d, Dirty: %v", version, dirty)

	if dirty {
		log.Println("Database is in dirty state, cleaning up...")

		// First, try to rollback the failed migration by dropping the new columns if they exist
		log.Println("Attempting to drop any partially created columns...")
		_, err = db.Exec(`
			ALTER TABLE work_orders DROP COLUMN IF EXISTS visit_number;
			ALTER TABLE work_orders DROP COLUMN IF EXISTS specimen_collection_date;
			ALTER TABLE work_orders DROP COLUMN IF EXISTS result_release_date;
			ALTER TABLE work_orders DROP COLUMN IF EXISTS diagnosis;
		`)
		// SQLite doesn't support DROP COLUMN in older versions, ignore error
		if err != nil {
			log.Printf("Note: Could not drop columns (may not exist or SQLite version issue): %v", err)
		}

		// Set dirty flag to false and keep at previous version
		prevVersion := version - 1
		if prevVersion < 0 {
			prevVersion = 0
		}

		log.Printf("Setting version to %d and marking as clean...", prevVersion)
		_, err = db.Exec("UPDATE schema_migrations SET version = ?, dirty = 0", prevVersion)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("✅ Successfully cleaned dirty state!")
		log.Println("Database is now at version:", prevVersion)
		log.Println("You can now restart the application to apply migrations.")
	} else {
		log.Println("Database is not dirty. No action needed.")
	}
}
