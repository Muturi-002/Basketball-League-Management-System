package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Player maps every column from the players table.
// PlayerPhoto is an OCI Object Storage URL stored as VARCHAR(200) — not binary.
type Player struct {
	PlayerID            int64     `json:"playerId"`
	FirstName           string    `json:"firstName"`
	MiddleName          *string   `json:"middleName,omitempty"`
	LastName            string    `json:"lastName"`
	PlayerPhoto         *string   `json:"playerPhoto,omitempty"`
	DateOfBirth         time.Time `json:"dateOfBirth"`
	Age                 int       `json:"age"`
	TeamID              int64     `json:"teamId"`
	PlayerHistory       *string   `json:"playerHistory,omitempty"`
	JerseyNumber        int       `json:"jerseyNumber"`
	PositionCode        string    `json:"positionCode"`
	CurrentSeasonStats  *string   `json:"currentSeasonStats,omitempty"`
	CareerStats         *string   `json:"careerStats,omitempty"`
}

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
			player_history,
			jersey_number,
			position_code,
			current_season_stats,
			career_stats
		FROM players
		ORDER BY player_id`)
	if err != nil {
		return nil, fmt.Errorf("query players: %w", err)
	}
	defer rows.Close()

	players := make([]Player, 0)
	for rows.Next() {
		var player Player
		var (
			middleName         sql.NullString
			playerPhoto        sql.NullString
			playerHistory      sql.NullString
			currentSeasonStats sql.NullString
			careerStats        sql.NullString
		)

		if err := rows.Scan(
			&player.PlayerID,
			&player.FirstName,
			&middleName,
			&player.LastName,
			&playerPhoto,
			&player.DateOfBirth,
			&player.Age,
			&player.TeamID,
			&playerHistory,
			&player.JerseyNumber,
			&player.PositionCode,
			&currentSeasonStats,
			&careerStats,
		); err != nil {
			return nil, fmt.Errorf("scan player row: %w", err)
		}

		player.MiddleName = nullStringPtr(middleName)
		player.PlayerPhoto = nullStringPtr(playerPhoto)
		player.PlayerHistory = nullStringPtr(playerHistory)
		player.CurrentSeasonStats = nullStringPtr(currentSeasonStats)
		player.CareerStats = nullStringPtr(careerStats)

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
	var (
		middleName         sql.NullString
		playerPhoto        sql.NullString
		playerHistory      sql.NullString
		currentSeasonStats sql.NullString
		careerStats        sql.NullString
	)

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
			player_history,
			jersey_number,
			position_code,
			current_season_stats,
			career_stats
		FROM players
		WHERE player_id = :1`, playerID).Scan(
		&player.PlayerID,
		&player.FirstName,
		&middleName,
		&player.LastName,
		&playerPhoto,
		&player.DateOfBirth,
		&player.Age,
		&player.TeamID,
		&playerHistory,
		&player.JerseyNumber,
		&player.PositionCode,
		&currentSeasonStats,
		&careerStats,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query player by id: %w", err)
	}

	player.MiddleName = nullStringPtr(middleName)
	player.PlayerPhoto = nullStringPtr(playerPhoto)
	player.PlayerHistory = nullStringPtr(playerHistory)
	player.CurrentSeasonStats = nullStringPtr(currentSeasonStats)
	player.CareerStats = nullStringPtr(careerStats)

	return &player, nil
}

// nullStringPtr converts a sql.NullString to a *string.
func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
