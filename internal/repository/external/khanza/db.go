package khanza

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	_ "github.com/go-sql-driver/mysql"
)

// DB represents a MySQL database connection wrapper
type DB struct {
	conn   *sql.DB
	config *config.Schema
}

// NewBridgeDB creates a new MySQL database connection using configuration
func NewBridgeDB(cfg *config.Schema) (*DB, error) {
	// Build MySQL connection string
	dsn := cfg.KhanzaBridgeDatabaseDSN

	// Open database connection
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("Connection to khanza database success")

	return &DB{
		conn:   conn,
		config: cfg,
	}, nil
}

// NewMainDB creates a new MySQL database connection using configuration
func NewMainDB(cfg *config.Schema) (*DB, error) {
	// Build MySQL connection string
	dsn := cfg.KhanzaMainDatabaseDSN

	// Open database connection
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	slog.Info("Connection to khanza database success")

	return &DB{
		conn:   conn,
		config: cfg,
	}, nil
}

// GetConnection returns the underlying sql.DB connection
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Ping tests the database connection
func (db *DB) Ping() error {
	return db.conn.Ping()
}
