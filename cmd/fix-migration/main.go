package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open database connection
	db, err := sql.Open("sqlite3", "../../tmp/biosystem-lims.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Create migrate driver
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal("Failed to create migrate driver:", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	// Check current version and dirty state
	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal("Failed to get migration version:", err)
	}

	fmt.Printf("Current migration version: %d, dirty: %t\n", version, dirty)

	if dirty {
		fmt.Println("Database is in dirty state. Forcing version...")

		// Force to the previous version
		prevVersion := version - 1
		err = m.Force(int(prevVersion))
		if err != nil {
			log.Fatal("Failed to force migration version:", err)
		}

		fmt.Printf("Forced migration to version: %d\n", prevVersion)

		// Now try to migrate up again
		fmt.Println("Attempting to migrate up...")
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to migrate up:", err)
		}

		fmt.Println("Migration completed successfully!")
	} else {
		fmt.Println("Database is not dirty. No action needed.")
	}
}
