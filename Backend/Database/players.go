package database

import (
	"database/sql"
	"fmt"
	"time"
)

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
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("players list: %w", err)
	}

	rows, err := conn.Query(`
        SELECT
            player_id,
            first_name,
            middle_name,
            last_name,
            player_photo,
            date_of_birth,
            age,
            team_id,
            player_history
        FROM players
        ORDER BY player_id`)
	if err != nil {
		return nil, fmt.Errorf("query players: %w", err)
	}
	defer rows.Close()

	players := make([]Player, 0)
	for rows.Next() {
		var player Player
		var middleName sql.NullString

		if err := rows.Scan(
			&player.PlayerID,
			&player.FirstName,
			&middleName,
			&player.LastName,
			&player.PlayerPhoto,
			&player.DateOfBirth,
			&player.Age,
			&player.TeamID,
			&player.PlayerHistory,
		); err != nil {
			return nil, fmt.Errorf("scan player row: %w", err)
		}

		player.MiddleName = nullStringPtr(middleName)
		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate players: %w", err)
	}

	return players, nil
}

func GetPlayerByID(playerID int64) (*Player, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("player get: %w", err)
	}

	var player Player
	var middleName sql.NullString

	err = conn.QueryRow(`
        SELECT
            player_id,
            first_name,
            middle_name,
            last_name,
            player_photo,
            date_of_birth,
            age,
            team_id,
            player_history
        FROM players
        WHERE player_id = :1`, playerID).Scan(
		&player.PlayerID,
		&player.FirstName,
		&middleName,
		&player.LastName,
		&player.PlayerPhoto,
		&player.DateOfBirth,
		&player.Age,
		&player.TeamID,
		&player.PlayerHistory,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query player by id: %w", err)
	}

	player.MiddleName = nullStringPtr(middleName)
	return &player, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
