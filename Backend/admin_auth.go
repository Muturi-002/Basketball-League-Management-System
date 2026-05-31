package main

import (
    "crypto/rand"
    "encoding/base64"
    db "blms/Database"
    "net/http"
    "os"
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

func adminCredentials() (string, string) {
    username := strings.TrimSpace(os.Getenv("BLMS_ADMIN_USERNAME"))
    password := strings.TrimSpace(os.Getenv("BLMS_ADMIN_PASSWORD"))
    if username == "" {
        username = "admin"
    }
    if password == "" {
        password = "admin123"
    }
    return username, password
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
        Username string `json:"username"`
        Password string `json:"password"`
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

    // Authenticate against admin_users DB table first
    ok, err := db.AuthenticateAdmin(payload.Username, payload.Password)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to authenticate admin")
        return
    }
    if !ok {
        // Fallback: allow env-based credentials for quick local overrides
        expectedUsername, expectedPassword := adminCredentials()
        if payload.Username != expectedUsername || payload.Password != expectedPassword {
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

    writeJSON(w, http.StatusOK, map[string]string{"status": "admin logged in"})
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
