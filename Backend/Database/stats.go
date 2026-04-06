package database

import (
	"database/sql"
	"fmt"
)

// Stat maps team-level league statistics.
type Stat struct {
	StatID        int64 `json:"statId"`
	TeamID        int64 `json:"teamId"`
	GamesPlayed   int   `json:"gamesPlayed"`
	Wins          int   `json:"wins"`
	Losses        int   `json:"losses"`
	Points        int   `json:"points"`
	PointsScored  int   `json:"pointsScored"`
	PointsAllowed int   `json:"pointsAllowed"`
}

func ListStats() ([]Stat, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("stats list: %w", err)
	}

	rows, err := conn.Query(`
		SELECT
			stat_id,
			team_id,
			games_played,
			wins,
			losses,
			points,
			points_scored,
			points_allowed
		FROM stats
		ORDER BY points DESC, wins DESC`)
	if err != nil {
		return nil, fmt.Errorf("query stats: %w", err)
	}
	defer rows.Close()

	stats := make([]Stat, 0)
	for rows.Next() {
		var stat Stat
		if err := rows.Scan(
			&stat.StatID,
			&stat.TeamID,
			&stat.GamesPlayed,
			&stat.Wins,
			&stat.Losses,
			&stat.Points,
			&stat.PointsScored,
			&stat.PointsAllowed,
		); err != nil {
			return nil, fmt.Errorf("scan stat row: %w", err)
		}
		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stats: %w", err)
	}

	return stats, nil
}

func GetStatByID(statID int64) (*Stat, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("stat get: %w", err)
	}

	var stat Stat
	err = conn.QueryRow(`
		SELECT
			stat_id,
			team_id,
			games_played,
			wins,
			losses,
			points,
			points_scored,
			points_allowed
		FROM stats
		WHERE stat_id = :1`, statID).Scan(
		&stat.StatID,
		&stat.TeamID,
		&stat.GamesPlayed,
		&stat.Wins,
		&stat.Losses,
		&stat.Points,
		&stat.PointsScored,
		&stat.PointsAllowed,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query stat by id: %w", err)
	}

	return &stat, nil
}
