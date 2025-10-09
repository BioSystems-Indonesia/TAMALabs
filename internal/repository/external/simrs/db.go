package simrs

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

// DB is a database connection wrapper
type DB struct {
	conn *sql.DB
}

// NewDB creates a new DB instance
func NewDB(dsn string) (*DB, error) {
	conn, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// GetConnection returns the database connection
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Ping tests the database connection
func (db *DB) Ping() error {
	return db.conn.Ping()
}
