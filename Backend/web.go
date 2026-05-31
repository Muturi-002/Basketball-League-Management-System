package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	db "blms/Database"
)

func main() {
	if err := db.ConnectFromEnv(); err != nil {
		log.Printf("Database connection not ready at startup: %v", err)
	} else {
		log.Println("Database connection established.")
	}

	publicMux := http.NewServeMux()
	adminMux := http.NewServeMux()

	// Public read APIs.
	publicMux.HandleFunc("/api/health", healthHandler)
	publicMux.HandleFunc("/api/players", playersHandler)
	publicMux.HandleFunc("/api/players/", playerByIDHandler)
	publicMux.HandleFunc("/api/teams", teamsHandler)
	publicMux.HandleFunc("/api/teams/", teamByIDHandler)
	publicMux.HandleFunc("/api/stadiums", stadiumsHandler)
	publicMux.HandleFunc("/api/stadiums/", stadiumByIDHandler)
	publicMux.HandleFunc("/api/managers", managersHandler)
	publicMux.HandleFunc("/api/managers/", managerByIDHandler)
	publicMux.HandleFunc("/api/injuries", injuriesHandler)
	publicMux.HandleFunc("/api/injured-players", injuredPlayersHandler)
	publicMux.HandleFunc("/api/stats", statsHandler)
	publicMux.HandleFunc("/api/stats/", statByIDHandler)
	publicMux.HandleFunc("/api/fixtures", fixturesHandler)
	publicMux.HandleFunc("/api/fixtures/", fixtureByIDHandler)
	publicMux.HandleFunc("/api/auth/login", loginHandler)
	publicMux.HandleFunc("/api/auth/register", registerHandler)
	publicMux.HandleFunc("/api/auth/users", usersHandler)
	publicMux.HandleFunc("/api/auth/users/", userByIDHandler)

	// Admin login/logout and protected APIs.
	adminMux.HandleFunc("/api/admin/login", adminLoginHandler)
	adminMux.HandleFunc("/api/admin/logout", adminLogoutHandler)
	adminMux.Handle("/api/admin/players", requireAdmin(http.HandlerFunc(adminPlayersHandler)))
	adminMux.Handle("/api/admin/players/", requireAdmin(http.HandlerFunc(adminPlayerByIDHandler)))
	adminMux.Handle("/api/admin/teams", requireAdmin(http.HandlerFunc(adminTeamsHandler)))
	adminMux.Handle("/api/admin/teams/", requireAdmin(http.HandlerFunc(adminTeamByIDHandler)))
	adminMux.Handle("/api/admin/stadiums", requireAdmin(http.HandlerFunc(adminStadiumsHandler)))
	adminMux.Handle("/api/admin/stadiums/", requireAdmin(http.HandlerFunc(adminStadiumByIDHandler)))
	adminMux.Handle("/api/admin/managers", requireAdmin(http.HandlerFunc(adminManagersHandler)))
	adminMux.Handle("/api/admin/managers/", requireAdmin(http.HandlerFunc(adminManagerByIDHandler)))
	adminMux.Handle("/api/admin/stats", requireAdmin(http.HandlerFunc(adminStatsHandler)))
	adminMux.Handle("/api/admin/stats/", requireAdmin(http.HandlerFunc(adminStatByIDHandler)))
	adminMux.Handle("/api/admin/fixtures", requireAdmin(http.HandlerFunc(adminFixturesHandler)))
	adminMux.Handle("/api/admin/fixtures/", requireAdmin(http.HandlerFunc(adminFixtureByIDHandler)))
	adminMux.HandleFunc("/admin-login.html", adminLoginPageHandler)
	adminMux.HandleFunc("/", adminPageHandler)
	adminMux.HandleFunc("/styles.css", frontendHandler)
	adminMux.HandleFunc("/auth.html", frontendHandler)

	// Frontend pages and assets.
	publicMux.HandleFunc("/", frontendHandler)

	publicAddr := ":4000"
	adminAddr := ":4900"

	go func() {
		log.Printf("Starting admin server on %s...", adminAddr)
		if err := http.ListenAndServe(adminAddr, withCORS(adminMux)); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Starting public server on %s...", publicAddr)
	if err := http.ListenAndServe(publicAddr, withCORS(publicMux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseIDFromPath(path, prefix string) (int64, bool) {
	if !strings.HasPrefix(path, prefix) {
		return 0, false
	}
	rawID := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if rawID == "" || strings.Contains(rawID, "/") {
		return 0, false
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func decodeJSON(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dest)
}

func writeList[T any](w http.ResponseWriter, r *http.Request, load func() ([]T, error), errorMessage string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	items, err := load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, errorMessage)
		return
	}
	if items == nil {
		items = []T{}
	}
	writeJSON(w, http.StatusOK, items)
}

func writeByID[T any](w http.ResponseWriter, r *http.Request, prefix, invalidIDMessage, loadErrorMessage, notFoundMessage string, load func(int64) (*T, error)) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, ok := parseIDFromPath(r.URL.Path, prefix)
	if !ok {
		writeError(w, http.StatusBadRequest, invalidIDMessage)
		return
	}

	item, err := load(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, loadErrorMessage)
		return
	}
	if item == nil {
		writeError(w, http.StatusNotFound, notFoundMessage)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/internal/") {
		http.NotFound(w, r)
		return
	}

	frontendFiles := map[string]bool{
		"home.html":        true,
		"teams.html":       true,
		"players.html":     true,
		"injury.html":      true,
		"stats.html":       true,
		"fixtures.html":    true,
		"auth.html":        true,
		"admin-login.html": true,
		"manager.html":     true,
		"styles.css":       true,
	}

	page := "home.html"
	if r.URL.Path != "/" {
		candidate := strings.TrimPrefix(r.URL.Path, "/")
		if !frontendFiles[candidate] {
			http.NotFound(w, r)
			return
		}
		page = candidate
	}

	http.ServeFile(w, r, "../Frontend/"+page)
}

func adminPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/admin-login.html", http.StatusFound)
		return
	}
	if r.URL.Path != "/admin.html" {
		http.NotFound(w, r)
		return
	}
	if !isAdminAuthenticated(r) {
		http.Redirect(w, r, "/admin-login.html", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "../Frontend/admin.html")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	_, err := db.EnsureConnected()
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status": "disconnected",
			"error":  err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "connected"})
}

func playersHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListPlayers, "failed to load players")
}

func playerByIDHandler(w http.ResponseWriter, r *http.Request) {
	writeByID(w, r, "/api/players/", "invalid player id", "failed to load player", "player not found", db.GetPlayerByID)
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListClubs, "failed to load clubs")
}

func teamByIDHandler(w http.ResponseWriter, r *http.Request) {
	writeByID(w, r, "/api/teams/", "invalid team id", "failed to load team", "team not found", db.GetClubByID)
}

func stadiumsHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListStadiums, "failed to load stadiums")
}

func stadiumByIDHandler(w http.ResponseWriter, r *http.Request) {
	writeByID(w, r, "/api/stadiums/", "invalid stadium id", "failed to load stadium", "stadium not found", db.GetStadiumByID)
}

func managersHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListManagers, "failed to load managers")
}

func managerByIDHandler(w http.ResponseWriter, r *http.Request) {
	writeByID(w, r, "/api/managers/", "invalid manager id", "failed to load manager", "manager not found", db.GetManagerByID)
}

func injuriesHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListInjuries, "failed to load injuries")
}

func injuredPlayersHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListInjuredPlayers, "failed to load injured players")
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListStats, "failed to load stats")
}

func statByIDHandler(w http.ResponseWriter, r *http.Request) {
	writeByID(w, r, "/api/stats/", "invalid stat id", "failed to load stat", "stat not found", db.GetStatByID)
}

func fixturesHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListFixtures, "failed to load fixtures")
}

func fixtureByIDHandler(w http.ResponseWriter, r *http.Request) {
	writeByID(w, r, "/api/fixtures/", "invalid fixture id", "failed to load fixture", "fixture not found", db.GetFixtureByID)
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type registerRequest struct {
	Username     string `json:"username"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
	FavTeamID    string `json:"favTeamId"`
}

func authMemberPayload(user *db.AuthUser) map[string]string {
	return map[string]string{
		"username":     user.UserID,
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"emailAddress": user.EmailAddress,
		"favTeamId":    user.FavTeamID,
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid login payload")
		return
	}
	if err := db.ValidateLoginIdentifier(req.Identifier); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.ValidatePassword(req.Password, 8, 72); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := db.AuthenticateUser(req.Identifier, req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to authenticate member")
		return
	}
	if user == nil {
		writeError(w, http.StatusUnauthorized, "invalid member credentials")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"member":  authMemberPayload(user),
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req registerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid registration payload")
		return
	}
	if err := db.ValidateUsername(req.Username, 3, 8); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.ValidateEmailAddress(req.EmailAddress); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := db.ValidatePassword(req.Password, 8, 72); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := db.CreateUser(db.User{
		UserID:       req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		EmailAddress: req.EmailAddress,
		Password:     req.Password,
		FavTeamID:    req.FavTeamID,
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "already exists"):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Member account created",
		"member":  authMemberPayload(user),
	})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	writeList(w, r, db.ListUsers, "failed to load users")
}

func userByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/auth/users/"), "/")
	if userID == "" || strings.Contains(userID, "/") {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := db.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load user")
		return
	}
	if user == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, user)
}
