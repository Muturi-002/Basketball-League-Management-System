package database

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"net/url"

	_ "github.com/sijms/go-ora/v2"
)

var (
	db      *sql.DB
	dbOnce  sync.Once
	dbErr   error
	envOnce sync.Once
)

// Connect opens and pings a database connection using the given driver and DSN.
func Connect(driverName, connString string) error {
	driverName = strings.TrimSpace(driverName)
	connString = strings.TrimSpace(connString)

	if driverName == "" {
		return errors.New("driver name is required")
	}
	if connString == "" {
		return errors.New("connection string is required")
	}

	dbOnce.Do(func() {
		conn, err := sql.Open(driverName, connString)
		if err != nil {
			dbErr = fmt.Errorf("open connection: %w", err)
			return
		}
		if err := conn.Ping(); err != nil {
			dbErr = fmt.Errorf("ping database: %w", err)
			return
		}
		db = conn
	})

	return dbErr
}

// ConnectFromEnv loads the .env file (if present) and opens the database
func ConnectFromEnv() error {
	envOnce.Do(func() {
		_ = loadEnvFile(findEnvFile())
	})

	driver := strings.TrimSpace(os.Getenv("ORACLE_DRIVER"))
	if driver == "" {
		// go-ora registers the driver as "oracle"
		driver = "oracle"
	}

	connString, err := buildConnString()
	if err != nil {
		return err
	}

	return Connect(driver, connString)
}

// EnsureConnected returns the active *sql.DB
func EnsureConnected() (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	if err := ConnectFromEnv(); err != nil {
		return nil, err
	}
	return db, nil
}

// DB returns the shared *sql.DB instance (may be nil before Connect is called).
func DB() *sql.DB { return db }

// IsConnected reports whether a DB handle has been initialised.
func IsConnected() bool { return db != nil }

// findEnvFile walks up from the current working directory looking for a .env file.
func findEnvFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return ".env"
	}
	for {
		candidate := filepath.Join(dir, ".env")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ".env"
}

// loadEnvFile reads key=value pairs from path into the process environment
func loadEnvFile(path string) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}
		if len(value) >= 2 {
			q := value[0]
			if (q == '"' || q == '\'') && value[len(value)-1] == q {
				value = value[1 : len(value)-1]
			}
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// buildConnString assembles a go-ora JDBC-like URL from environment variables. It prefers building an "oracle://user:pass@..." URL; if the ORACLE_CONNECTION_STRING contains a full connect descriptor the descriptor is URL-escaped and passed as the `connectString` query parameter which go-ora understands.
func buildConnString() (string, error) {
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	connStr := strings.TrimSpace(os.Getenv("ORACLE_CONNECTION_STRING"))

	if user == "" || password == "" || connStr == "" {
		return "", errors.New("missing required env vars: DB_USER, DB_PASSWORD, ORACLE_CONNECTION_STRING")
	}

	configDir := strings.TrimSpace(os.Getenv("ORACLE_CONFIG_DIR"))

	// Extracting host, port and service_name from a full descriptor like (description=...(address=(protocol=tcps)(port=1522)(host=...))(connect_data=(service_name=...)))
	extract := func(key string) string {
		lower := strings.ToLower(connStr)
		k := strings.ToLower(key) + "="
		i := strings.Index(lower, k)
		if i < 0 {
			return ""
		}
		i += len(k)
		j := i
		for j < len(lower) {
			c := lower[j]
			if c == ')' || c == '(' || c == ' ' || c == '\n' || c == '\r' || c == '\t' {
				break
			}
			j++
		}
		return connStr[i:j]
	}

	host := extract("host")
	port := extract("port")
	svc := extract("service_name")

	if host != "" && port != "" && svc != "" {
		// Build simple URL form with go-ora supported options
		jdbc := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", url.PathEscape(user), url.PathEscape(password), host, port, svc)
		opts := url.Values{}
		if configDir != "" {
			opts.Set("WALLET", configDir)
			opts.Set("AUTH TYPE", "TCPS")
			opts.Set("SSL", "TRUE")
			opts.Set("SSL VERIFY", "TRUE")
		}
		if encoded := opts.Encode(); encoded != "" {
			jdbc += "?" + encoded
		}
		return jdbc, nil
	}

	// Fallback: pass the full descriptor as connectString query parameter
	escaped := url.QueryEscape(connStr)
	jdbc := fmt.Sprintf("oracle://%s:%s@?connectString=%s", url.PathEscape(user), url.PathEscape(password), escaped)
	if configDir != "" {
		jdbc += fmt.Sprintf("&WALLET=%s&AUTH+TYPE=TCPS&SSL=TRUE&SSL+VERIFY=TRUE", url.QueryEscape(configDir))
	}
	return jdbc, nil
}
