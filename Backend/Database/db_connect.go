package database

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"

	_ "github.com/godror/godror"
)

var db *sql.DB
var envOnce sync.Once

func Connect(driverName, conn_string string) error {
	driverName = strings.TrimSpace(driverName)
	conn_string = strings.TrimSpace(conn_string)

	if driverName == "" {
		return errors.New("driver name is required")
	}
	if conn_string == "" {
		return errors.New("conn_string is required")
	}

	if db != nil {
		return nil
	}

	conn, err := sql.Open(driverName, conn_string)
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
	envOnce.Do(func() {
		_ = loadEnvFile(resolveEnvPath())
	})

	driver := strings.TrimSpace(os.Getenv("ORACLE_DRIVER"))
	if driver == "" {
		driver = "godror"
	}

	conn_string, err := buildStringFromEnv()
	if err != nil {
		return err
	}

	return Connect(driver, conn_string)
}

func resolveEnvPath() string {
	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		for _, candidate := range []string{
			filepath.Join(filepath.Dir(sourceFile), "..", "..", ".env"),
			filepath.Join(filepath.Dir(sourceFile), "..", ".env"),
			".env",
		} {
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return candidate
			}
		}
	}

	return ".env"
}

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
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

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
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func buildStringFromEnv() (string, error) {

	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	connectionString := strings.TrimSpace(os.Getenv("ORACLE_CONNECTION_STRING"))
	configDir := strings.TrimSpace(os.Getenv("ORACLE_CONFIG_DIR"))
	if configDir == "" {
		configDir = strings.TrimSpace(os.Getenv("TNS_PATH"))
	}
	libDir := strings.TrimSpace(os.Getenv("ORACLE_LIB_DIR"))
	if libDir == "" {
		libDir = strings.TrimSpace(os.Getenv("ORACLE_CLIENT_LIB_DIR"))
	}

	if user == "" || password == "" || connectionString == "" {
		return "", errors.New("missing DB_USER/DB_PASSWORD/ORACLE_CONNECTION_STRING. Check your environment variables")
	}

	if configDir != "" {
		_ = os.Setenv("TNS_ADMIN", configDir)
	}
	if libDir != "" {
		current := strings.TrimSpace(os.Getenv("LD_LIBRARY_PATH"))
		if current == "" {
			_ = os.Setenv("LD_LIBRARY_PATH", libDir)
		} else if !strings.Contains(current, libDir) {
			_ = os.Setenv("LD_LIBRARY_PATH", libDir+":"+current)
		}
	}

	conn_string := fmt.Sprintf(`user="%s" password="%s" connectString="%s"`, user, password, connectionString)
	if configDir != "" {
		conn_string += fmt.Sprintf(` configDir="%s"`, configDir)
	}
	if libDir != "" {
		conn_string += fmt.Sprintf(` libDir="%s"`, libDir)
	}

	return conn_string, nil
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

func init() {
	// If we've already re-exec'd with BLMS_LAUNCHED=1, skip.
	if os.Getenv("BLMS_LAUNCHED") == "1" {
		return
	}

	// Prefer ORACLE_LIB_DIR; otherwise fall back to ../oracle-client when running from Backend.
	desired := strings.TrimSpace(os.Getenv("ORACLE_LIB_DIR"))
	if desired == "" {
		desired = filepath.Join("..", "oracle-client")
	}

	ld := strings.TrimSpace(os.Getenv("LD_LIBRARY_PATH"))
	if ld == "" || !strings.Contains(ld, desired) {
		var newLD string
		if ld == "" {
			newLD = desired
		} else {
			newLD = desired + ":" + ld
		}

		env := os.Environ()
		// Replace or append LD_LIBRARY_PATH in env list
		replaced := false
		for i, e := range env {
			if strings.HasPrefix(e, "LD_LIBRARY_PATH=") {
				env[i] = "LD_LIBRARY_PATH=" + newLD
				replaced = true
				break
			}
		}
		if !replaced {
			env = append(env, "LD_LIBRARY_PATH="+newLD)
		}
		env = append(env, "BLMS_LAUNCHED=1")

		exe, err := os.Executable()
		if err != nil {
			log.Printf("failed to determine executable for re-exec: %v", err)
			return
		}

		// Re-exec the current process so the dynamic loader sees LD_LIBRARY_PATH.
		if err := syscall.Exec(exe, os.Args, env); err != nil {
			log.Printf("re-exec failed: %v", err)
		}
	}
}