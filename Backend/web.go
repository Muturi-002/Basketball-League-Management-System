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

	mux := http.NewServeMux()

	// Public read APIs.
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/players", playersHandler)
	mux.HandleFunc("/api/players/", playerByIDHandler)
	mux.HandleFunc("/api/teams", teamsHandler)
	mux.HandleFunc("/api/teams/", teamByIDHandler)
	mux.HandleFunc("/api/stadiums", stadiumsHandler)
	mux.HandleFunc("/api/stadiums/", stadiumByIDHandler)
	mux.HandleFunc("/api/managers", managersHandler)
	mux.HandleFunc("/api/managers/", managerByIDHandler)
	mux.HandleFunc("/api/injuries", injuriesHandler)
	mux.HandleFunc("/api/injured-players", injuredPlayersHandler)
	mux.HandleFunc("/api/stats", statsHandler)
	mux.HandleFunc("/api/stats/", statByIDHandler)
	mux.HandleFunc("/api/fixtures", fixturesHandler)
	mux.HandleFunc("/api/fixtures/", fixtureByIDHandler)
	mux.HandleFunc("/api/auth/users", usersHandler)

	// Admin write APIs - redirect to admin service on port 4900
	mux.HandleFunc("/api/admin/players", adminRedirectHandler)
	mux.HandleFunc("/api/admin/players/", adminRedirectHandler)
	mux.HandleFunc("/api/admin/teams", adminRedirectHandler)
	mux.HandleFunc("/api/admin/teams/", adminRedirectHandler)
	mux.HandleFunc("/api/admin/stadiums", adminRedirectHandler)
	mux.HandleFunc("/api/admin/stadiums/", adminRedirectHandler)
	mux.HandleFunc("/api/admin/managers", adminRedirectHandler)
	mux.HandleFunc("/api/admin/managers/", adminRedirectHandler)
	mux.HandleFunc("/api/admin/stats", adminRedirectHandler)
	mux.HandleFunc("/api/admin/stats/", adminRedirectHandler)
	mux.HandleFunc("/api/admin/fixtures", adminRedirectHandler)
	mux.HandleFunc("/api/admin/fixtures/", adminRedirectHandler)

	// Frontend pages and assets.
	mux.HandleFunc("/internal/admin.html", adminPageHandler)
	mux.HandleFunc("/", frontendHandler)

	addr := ":4000"
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
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

func frontendHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/internal/") {
		http.NotFound(w, r)
		return
	}

	frontendFiles := map[string]bool{
		"home.html":     true,
		"teams.html":    true,
		"players.html":  true,
		"injury.html":   true,
		"stats.html":    true,
		"fixtures.html": true,
		"auth.html":     true,
		"manager.html":  true,
		"styles.css":    true,
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
	http.ServeFile(w, r, "../Frontend/admin.html")
}

func adminRedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Construct the redirect URL to the admin service on port 4900
	redirectURL := "http://localhost:4900" + r.URL.Path
	if r.URL.RawQuery != "" {
		redirectURL += "?" + r.URL.RawQuery
	}

	// Preserve the request method for the redirect (307 Temporary Redirect)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
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
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	players, err := db.ListPlayers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load players")
		return
	}
	if players == nil {
		players = []db.Player{}
	}
	writeJSON(w, http.StatusOK, players)
}

func playerByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	playerID, ok := parseIDFromPath(r.URL.Path, "/api/players/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid player id")
		return
	}

	player, err := db.GetPlayerByID(playerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load player")
		return
	}
	if player == nil {
		writeError(w, http.StatusNotFound, "player not found")
		return
	}
	writeJSON(w, http.StatusOK, player)
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	clubs, err := db.ListClubs()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load clubs")
		return
	}
	if clubs == nil {
		clubs = []db.Club{}
	}
	writeJSON(w, http.StatusOK, clubs)
}

func teamByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	clubID, ok := parseIDFromPath(r.URL.Path, "/api/teams/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid team id")
		return
	}

	club, err := db.GetClubByID(clubID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load team")
		return
	}
	if club == nil {
		writeError(w, http.StatusNotFound, "team not found")
		return
	}
	writeJSON(w, http.StatusOK, club)
}

func stadiumsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stadiums, err := db.ListStadiums()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load stadiums")
		return
	}
	if stadiums == nil {
		stadiums = []db.Stadium{}
	}
	writeJSON(w, http.StatusOK, stadiums)
}

func stadiumByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stadiumID, ok := parseIDFromPath(r.URL.Path, "/api/stadiums/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid stadium id")
		return
	}

	stadium, err := db.GetStadiumByID(stadiumID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load stadium")
		return
	}
	if stadium == nil {
		writeError(w, http.StatusNotFound, "stadium not found")
		return
	}
	writeJSON(w, http.StatusOK, stadium)
}

func managersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	managers, err := db.ListManagers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load managers")
		return
	}
	if managers == nil {
		managers = []db.Manager{}
	}
	writeJSON(w, http.StatusOK, managers)
}

func managerByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	managerID, ok := parseIDFromPath(r.URL.Path, "/api/managers/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid manager id")
		return
	}

	manager, err := db.GetManagerByID(managerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load manager")
		return
	}
	if manager == nil {
		writeError(w, http.StatusNotFound, "manager not found")
		return
	}
	writeJSON(w, http.StatusOK, manager)
}

func injuriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	injuries, err := db.ListInjuries()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load injuries")
		return
	}
	if injuries == nil {
		injuries = []db.Injury{}
	}
	writeJSON(w, http.StatusOK, injuries)
}

func injuredPlayersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	injuredPlayers, err := db.ListInjuredPlayers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load injured players")
		return
	}
	if injuredPlayers == nil {
		injuredPlayers = []db.InjuredPlayer{}
	}
	writeJSON(w, http.StatusOK, injuredPlayers)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stats, err := db.ListStats()
	if err != nil {
		log.Printf("stats handler: failed to load stats: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to load stats")
		return
	}
	if stats == nil {
		stats = []db.Stat{}
	}
	writeJSON(w, http.StatusOK, stats)
}

func statByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	statID, ok := parseIDFromPath(r.URL.Path, "/api/stats/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid stat id")
		return
	}

	stat, err := db.GetStatByID(statID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load stat")
		return
	}
	if stat == nil {
		writeError(w, http.StatusNotFound, "stat not found")
		return
	}
	writeJSON(w, http.StatusOK, stat)
}

func fixturesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	fixtures, err := db.ListFixtures()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load fixtures")
		return
	}
	if fixtures == nil {
		fixtures = []db.Fixture{}
	}
	writeJSON(w, http.StatusOK, fixtures)
}

func fixtureByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	fixtureID, ok := parseIDFromPath(r.URL.Path, "/api/fixtures/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid fixture id")
		return
	}

	fixture, err := db.GetFixtureByID(fixtureID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load fixture")
		return
	}
	if fixture == nil {
		writeError(w, http.StatusNotFound, "fixture not found")
		return
	}
	writeJSON(w, http.StatusOK, fixture)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	users, err := db.ListUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load users")
		return
	}
	if users == nil {
		users = []db.User{}
	}
	writeJSON(w, http.StatusOK, users)
}
