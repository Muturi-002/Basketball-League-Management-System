package database

import "time"

// Player maps the fields from the Players table (players.png).
type Player struct {
    PlayerID      int64     `json:"playerId"`
    FirstName     string    `json:"firstName"`
    MiddleName    *string   `json:"middleName,omitempty"`
    LastName      string    `json:"lastName"`
    PlayerPhoto   []byte    `json:"playerPhoto"` // binary image data (JSON-encoded as base64)
    DateOfBirth   time.Time `json:"dateOfBirth"`
    Age           int       `json:"age"`
    TeamID        int64     `json:"teamId"`
    PlayerHistory string    `json:"playerHistory"`
}

// TODO: Implement database-backed operations (queries, inserts, updates)
// using the shared connection from DB().
func ListPlayers() ([]Player, error) {
    // Placeholder implementation for Step 1 draft.
    return nil, nil
}
