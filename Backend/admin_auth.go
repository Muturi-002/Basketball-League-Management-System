package main

import (
	db "blms/Database"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const adminSessionCookieName = "blms_admin_session"
const adminSessionTTL = 8 * time.Hour

type adminSessionStore struct {
	mu     sync.Mutex
	tokens map[string]time.Time
}

var adminSessions = adminSessionStore{tokens: make(map[string]time.Time)}

type adminCredential struct {
	Username     string
	PasswordHash string
}

func adminLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if r.URL.Path != "/admin-login.html" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "../Frontend/admin-login.html")
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var payload struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		CreatePassword  bool   `json:"createPassword,omitempty"`
		ConfirmPassword string `json:"confirmPassword,omitempty"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid login payload")
		return
	}

	if err := db.ValidateUsername(payload.Username, 3, 32); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.ValidatePassword(payload.Password, 8, 72); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	authStatus := "admin logged in"

	if payload.CreatePassword {
		if payload.Password != payload.ConfirmPassword {
			writeError(w, http.StatusBadRequest, "passwords do not match")
			return
		}

		created, err := upsertAdminPassword(payload.Username, payload.Password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save admin account")
			return
		}

		authStatus = "admin password updated"
		if created {
			authStatus = "admin account created"
		}
	}

	if !payload.CreatePassword {
		ok, err := authenticateAdmin(payload.Username, payload.Password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to authenticate admin")
			return
		}
		if !ok {
			writeError(w, http.StatusUnauthorized, "invalid admin credentials")
			return
		}
	}

	token, err := generateAdminToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	expiresAt := time.Now().Add(adminSessionTTL)
	storeAdminToken(token, expiresAt)

	http.SetCookie(w, &http.Cookie{
		Name:     adminSessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{"status": authStatus})
}

func adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token := adminTokenFromRequest(r)
	if token != "" {
		deleteAdminToken(token)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     adminSessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{"status": "admin logged out"})
}

func requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAdminAuthenticated(r) {
			writeError(w, http.StatusUnauthorized, "admin login required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isAdminAuthenticated(r *http.Request) bool {
	token := adminTokenFromRequest(r)
	if token == "" {
		return false
	}

	adminSessions.mu.Lock()
	defer adminSessions.mu.Unlock()

	expiresAt, ok := adminSessions.tokens[token]
	if !ok {
		return false
	}
	if time.Now().After(expiresAt) {
		delete(adminSessions.tokens, token)
		return false
	}
	return true
}

func adminTokenFromRequest(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}

	cookie, err := r.Cookie(adminSessionCookieName)
	if err == nil {
		return strings.TrimSpace(cookie.Value)
	}

	return ""
}

func generateAdminToken() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func storeAdminToken(token string, expiresAt time.Time) {
	adminSessions.mu.Lock()
	defer adminSessions.mu.Unlock()
	adminSessions.tokens[token] = expiresAt
}

func deleteAdminToken(token string) {
	adminSessions.mu.Lock()
	defer adminSessions.mu.Unlock()
	delete(adminSessions.tokens, token)
}

func getAdminCredential(username string) (*adminCredential, error) {
	conn, err := db.EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("admin credential get: %w", err)
	}

	if err := db.ValidateUsername(username, 3, 80); err != nil {
		return nil, err
	}

	var credential adminCredential
	err = conn.QueryRow(`
		SELECT username, password_hash
		FROM admin_users
		WHERE LOWER(username) = LOWER(:1)
	`, username).Scan(&credential.Username, &credential.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query admin credential: %w", err)
	}

	return &credential, nil
}

func authenticateAdmin(username, password string) (bool, error) {
	if err := db.ValidatePassword(password, 8, 72); err != nil {
		return false, err
	}

	credential, err := getAdminCredential(username)
	if err != nil {
		return false, fmt.Errorf("admin auth: %w", err)
	}
	if credential == nil {
		return false, nil
	}
	if !adminPasswordHashMatches(credential.PasswordHash, password) {
		return false, nil
	}

	return true, nil
}

func upsertAdminPassword(username, password string) (bool, error) {
	conn, err := db.EnsureConnected()
	if err != nil {
		return false, fmt.Errorf("admin password save: %w", err)
	}
	username = strings.TrimSpace(username)

	if err := db.ValidateUsername(username, 3, 80); err != nil {
		return false, err
	}
	if err := db.ValidatePassword(password, 8, 72); err != nil {
		return false, err
	}

	hashedPassword := adminHashPassword(password)
	tx, err := conn.Begin()
	if err != nil {
		return false, fmt.Errorf("begin admin password transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var existingCount int64
	err = tx.QueryRow(`
		SELECT COUNT(*)
		FROM admin_users
		WHERE LOWER(username) = LOWER(:1)
	`, username).Scan(&existingCount)
	if err != nil {
		return false, fmt.Errorf("check admin user: %w", err)
	}

	if existingCount > 0 {
		if _, err := tx.Exec(`
			UPDATE admin_users
			SET password_hash = :1
			WHERE LOWER(username) = LOWER(:2)
		`, hashedPassword, username); err != nil {
			return false, fmt.Errorf("update admin password: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("commit admin password update: %w", err)
		}
		return false, nil
	}

	if _, err := tx.Exec(`
		INSERT INTO admin_users (username, password_hash)
		VALUES (:1, :2)
	`, username, hashedPassword); err != nil {
		return false, fmt.Errorf("insert admin password: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit admin password insert: %w", err)
	}

	return true, nil
}

func adminPasswordHashMatches(storedHash, password string) bool {
	return strings.EqualFold(strings.TrimSpace(storedHash), adminHashPassword(password))
}

func adminHashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}
