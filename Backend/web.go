package main

import (
    "encoding/json"
    "log"
    "net/http"

    db "blms/Database"
)

// TODO: wire in database.Connect once DSN and driver are finalized.

func main() {
    mux := http.NewServeMux()

    // Serve static frontend files.
    fileServer := http.FileServer(http.Dir("../Frontend"))
    mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

    // Map root to Frontend/home.html
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "../Frontend/home.html")
    })

    // JSON API endpoints that expose data from the database package.
    mux.HandleFunc("/api/players", playersHandler)
    mux.HandleFunc("/api/teams", teamsHandler)
    mux.HandleFunc("/api/managers", managersHandler)
    mux.HandleFunc("/api/injuries", injuriesHandler)
    mux.HandleFunc("/api/stats", statsHandler)
    mux.HandleFunc("/api/fixtures", fixturesHandler)
    mux.HandleFunc("/api/auth/users", usersHandler)

    addr := ":8080"
    log.Printf("Starting server on %s...", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatal(err)
    }
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(value)
}

func playersHandler(w http.ResponseWriter, r *http.Request) {
    players, err := db.ListPlayers()
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load players"})
        return
    }
    if players == nil {
        players = []db.Player{}
    }
    writeJSON(w, http.StatusOK, players)
}

func teamsHandler(w http.ResponseWriter, r *http.Request) {
    clubs, err := db.ListClubs()
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load clubs"})
        return
    }
    if clubs == nil {
        clubs = []db.Club{}
    }
    writeJSON(w, http.StatusOK, clubs)
}

func managersHandler(w http.ResponseWriter, r *http.Request) {
    managers, err := db.ListManagers()
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load managers"})
        return
    }
    if managers == nil {
        managers = []db.Manager{}
    }
    writeJSON(w, http.StatusOK, managers)
}

func injuriesHandler(w http.ResponseWriter, r *http.Request) {
    injuries, err := db.ListInjuries()
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load injuries"})
        return
    }
    if injuries == nil {
        injuries = []db.Injury{}
    }
    writeJSON(w, http.StatusOK, injuries)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
    stats, err := db.ListStats()
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load stats"})
        return
    }
    if stats == nil {
        stats = []db.Stat{}
    }
    writeJSON(w, http.StatusOK, stats)
}

func fixturesHandler(w http.ResponseWriter, r *http.Request) {
    fixtures, err := db.ListFixtures()
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load fixtures"})
        return
    }
    if fixtures == nil {
        fixtures = []db.Fixture{}
    }
    writeJSON(w, http.StatusOK, fixtures)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    users, err := db.ListUsers()
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load users"})
        return
    }
    if users == nil {
        users = []db.User{}
    }
    writeJSON(w, http.StatusOK, users)
}
