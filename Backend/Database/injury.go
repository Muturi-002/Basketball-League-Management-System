package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Injury maps every column from the injuries table.
type Injury struct {
	InjuryID               int64   `json:"injuryId"`
	InjuryName             string  `json:"injuryName"`
	InjuryType             string  `json:"injuryType"`
	InjuryDescription      *string `json:"injuryDescription,omitempty"`
	InjuryResources        *string `json:"injuryResources,omitempty"`
	InjuryVideoDescription *string `json:"injuryVideoDescription,omitempty"`
}

// InjuredPlayer maps every column from the injured_players table.
// ExpectedTimeOfReturn is nullable (DATE, may be NULL for open-ended injuries).
type InjuredPlayer struct {
	InjuredPlayerID      int64      `json:"injuredPlayerId"`
	PlayerID             int64      `json:"playerId"`
	InjuryID             int64      `json:"injuryId"`
	InjuryName           string     `json:"injuryName"`
	ExpectedTimeOfReturn *time.Time `json:"expectedTimeOfReturn,omitempty"`
	InjuryStartDate      time.Time  `json:"injuryStartDate"`
	InjuryStatus         string     `json:"injuryStatus"`
}

func ListInjuries() ([]Injury, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("injuries list: %w", err)
	}

	rows, err := conn.Query(`
		SELECT
			injury_id,
			injury_name,
			injury_type,
			injury_description,
			injury_resources,
			injury_video_description
		FROM injuries
		ORDER BY injury_id`)
	if err != nil {
		return nil, fmt.Errorf("query injuries: %w", err)
	}
	defer rows.Close()

	injuries := make([]Injury, 0)
	for rows.Next() {
		var injury Injury
		var (
			description      sql.NullString
			resources        sql.NullString
			videoDescription sql.NullString
		)

		if err := rows.Scan(
			&injury.InjuryID,
			&injury.InjuryName,
			&injury.InjuryType,
			&description,
			&resources,
			&videoDescription,
		); err != nil {
			return nil, fmt.Errorf("scan injury row: %w", err)
		}

		injury.InjuryDescription = nullStringPtr(description)
		injury.InjuryResources = nullStringPtr(resources)
		injury.InjuryVideoDescription = nullStringPtr(videoDescription)

		injuries = append(injuries, injury)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate injuries: %w", err)
	}

	return injuries, nil
}

func GetInjuryByID(injuryID int64) (*Injury, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("injury get: %w", err)
	}

	var injury Injury
	var (
		description      sql.NullString
		resources        sql.NullString
		videoDescription sql.NullString
	)

	err = conn.QueryRow(`
		SELECT
			injury_id,
			injury_name,
			injury_type,
			injury_description,
			injury_resources,
			injury_video_description
		FROM injuries
		WHERE injury_id = :1`, injuryID).Scan(
		&injury.InjuryID,
		&injury.InjuryName,
		&injury.InjuryType,
		&description,
		&resources,
		&videoDescription,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query injury by id: %w", err)
	}

	injury.InjuryDescription = nullStringPtr(description)
	injury.InjuryResources = nullStringPtr(resources)
	injury.InjuryVideoDescription = nullStringPtr(videoDescription)

	return &injury, nil
}

func ListInjuredPlayers() ([]InjuredPlayer, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("injured players list: %w", err)
	}

	rows, err := conn.Query(`
		SELECT
			injured_player_id,
			player_id,
			injury_id,
			injury_name,
			expected_time_of_return,
			injury_start_date,
			injury_status
		FROM injured_players
		ORDER BY injury_start_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("query injured players: %w", err)
	}
	defer rows.Close()

	injured := make([]InjuredPlayer, 0)
	for rows.Next() {
		var ip InjuredPlayer
		var expectedReturn sql.NullTime

		if err := rows.Scan(
			&ip.InjuredPlayerID,
			&ip.PlayerID,
			&ip.InjuryID,
			&ip.InjuryName,
			&expectedReturn,
			&ip.InjuryStartDate,
			&ip.InjuryStatus,
		); err != nil {
			return nil, fmt.Errorf("scan injured player row: %w", err)
		}

		if expectedReturn.Valid {
			ip.ExpectedTimeOfReturn = &expectedReturn.Time
		}

		injured = append(injured, ip)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate injured players: %w", err)
	}

	return injured, nil
}

// ListCurrentInjuredPlayers returns only players with injury_status = 'CURRENT'.
func ListCurrentInjuredPlayers() ([]InjuredPlayer, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("current injured players list: %w", err)
	}

	rows, err := conn.Query(`
		SELECT
			injured_player_id,
			player_id,
			injury_id,
			injury_name,
			expected_time_of_return,
			injury_start_date,
			injury_status
		FROM injured_players
		WHERE injury_status = 'CURRENT'
		ORDER BY injury_start_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("query current injured players: %w", err)
	}
	defer rows.Close()

	injured := make([]InjuredPlayer, 0)
	for rows.Next() {
		var ip InjuredPlayer
		var expectedReturn sql.NullTime

		if err := rows.Scan(
			&ip.InjuredPlayerID,
			&ip.PlayerID,
			&ip.InjuryID,
			&ip.InjuryName,
			&expectedReturn,
			&ip.InjuryStartDate,
			&ip.InjuryStatus,
		); err != nil {
			return nil, fmt.Errorf("scan injured player row: %w", err)
		}

		if expectedReturn.Valid {
			ip.ExpectedTimeOfReturn = &expectedReturn.Time
		}

		injured = append(injured, ip)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate current injured players: %w", err)
	}

	return injured, nil
}

func GetInjuredPlayerByID(injuredPlayerID int64) (*InjuredPlayer, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("injured player get: %w", err)
	}

	var ip InjuredPlayer
	var expectedReturn sql.NullTime

	err = conn.QueryRow(`
		SELECT
			injured_player_id,
			player_id,
			injury_id,
			injury_name,
			expected_time_of_return,
			injury_start_date,
			injury_status
		FROM injured_players
		WHERE injured_player_id = :1`, injuredPlayerID).Scan(
		&ip.InjuredPlayerID,
		&ip.PlayerID,
		&ip.InjuryID,
		&ip.InjuryName,
		&expectedReturn,
		&ip.InjuryStartDate,
		&ip.InjuryStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query injured player by id: %w", err)
	}

	if expectedReturn.Valid {
		ip.ExpectedTimeOfReturn = &expectedReturn.Time
	}

	return &ip, nil
}
