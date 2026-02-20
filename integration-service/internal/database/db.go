// database/db.go
package database

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

var DB *gorm.DB
var dbFileName = getDBPath()

func Connect() {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	db.Exec("PRAGMA busy_timeout = 5000")
	db.Exec("PRAGMA journal_mode = WAL")

	DB = db
}

func getDBPath() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		return filepath.Join(localAppData, "TAMALabs", "database", "TAMALabs.db")
	}

	appData := os.Getenv("APPDATA")
	if appData != "" {
		return filepath.Join(appData, "TAMALabs", "database", "TAMALabs.db")
	}

	programData := os.Getenv("ProgramData")
	if programData != "" {
		slog.Warn("Using ProgramData for database - this may require admin privileges")
		return filepath.Join(programData, "TAMALabs", "database", "TAMALabs.db")
	}

	if runtime.GOOS == "windows" {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile != "" {
			return filepath.Join(userProfile, "AppData", "Local", "TAMALabs", "database", "TAMALabs.db")
		}
		slog.Warn("All environment variables empty, using hardcoded path")
		return `C:\Users\Public\TAMALabs\database\TAMALabs.db`
	}

	return filepath.Join(os.TempDir(), "TAMALabs", "database", "TAMALabs.db")
}
