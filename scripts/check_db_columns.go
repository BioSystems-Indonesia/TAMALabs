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

	// Check work_orders table structure
	rows, err := db.Query("PRAGMA table_info(work_orders)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	log.Println("\nwork_orders table columns:")
	log.Println("----------------------------")

	foundCols := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name string
		var typ string
		var notnull int
		var dfltValue sql.NullString
		var pk int

		err = rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%2d | %-30s | %-15s | notnull: %d | pk: %d", cid, name, typ, notnull, pk)
		foundCols[name] = true
	}

	log.Println("\n----------------------------")
	log.Println("Checking for new columns:")
	newCols := []string{"visit_number", "specimen_collection_date", "result_release_date", "diagnosis"}
	for _, col := range newCols {
		if foundCols[col] {
			log.Printf("✅ %s - EXISTS", col)
		} else {
			log.Printf("❌ %s - MISSING", col)
		}
	}
}
