package database

import (
	"database/sql"
	"fmt"
)

// NOTE: This file intentionally does not import or assume any specific
// third-party Oracle driver. The caller is expected to register a driver
// elsewhere in the application using only the tooling you choose.

var db *sql.DB

// Connect initializes a shared database handle for the application.
// The driverName argument should match the name of a registered database
// driver (for example, one you configure manually when wiring Oracle ATP).
func Connect(driverName, dsn string) error {
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

// DB exposes the shared *sql.DB instance.
func DB() *sql.DB {
    return db
}
