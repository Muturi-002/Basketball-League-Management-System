package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

var db *sql.DB

func Connect(driverName, dsn string) error {
	driverName = strings.TrimSpace(driverName)
	dsn = strings.TrimSpace(dsn)

	if driverName == "" {
		return errors.New("driver name is required")
	}
	if dsn == "" {
		return errors.New("dsn is required")
	}

	if db != nil {
		return nil
	}

	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		return fmt.Errorf("open connection: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	db = conn
	return nil
}

func ConnectFromEnv() error {
	driver := strings.TrimSpace(os.Getenv("ORACLE_DRIVER"))
	if driver == "" {
		driver = "oracle"
	}

	dsn := strings.TrimSpace(os.Getenv("ORACLE_DSN"))
	if dsn == "" {
		dsn = strings.TrimSpace(os.Getenv("BLMS_ORACLE_DSN"))
	}
	if dsn == "" {
		return errors.New("missing ORACLE_DSN or BLMS_ORACLE_DSN")
	}

	return Connect(driver, dsn)
}

// EnsureConnected verifies that a database connection is active.
func EnsureConnected() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}

	if err := ConnectFromEnv(); err != nil {
		return nil, err
	}
	return db, nil
}

// DB exposes the shared *sql.DB instance.
func DB() *sql.DB {
	return db
}

// IsConnected reports whether a DB handle has been initialized.
func IsConnected() bool {
	return db != nil
}
