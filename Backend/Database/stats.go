package database

import (
	"database/sql"
	"fmt"
)

// Stat maps team-level league statistics.
type Stat struct {
	StatID       int64 `json:"statId"`
	TeamID       int64 `json:"teamId"`
	GamesPlayed  int   `json:"gamesPlayed"`
	Wins         int   `json:"wins"`
	Losses       int   `json:"losses"`
	Points       int   `json:"points"`
	PointsScored int   `json:"pointsScored"`
}

func ListStats() ([]Stat, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("stats list: %w", err)
	}

	rows, err := conn.Query(`
		SELECT
			team_stat_id,
			team_id,
			games_played,
			wins,
			losses,
			league_points,
			points_scored
		FROM team_statistics
		ORDER BY league_points DESC, wins DESC`)
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
			team_stat_id,
			team_id,
			games_played,
			wins,
			losses,
			league_points,
			points_scored
		FROM team_statistics
		WHERE team_stat_id = :1`, statID).Scan(
		&stat.StatID,
		&stat.TeamID,
		&stat.GamesPlayed,
		&stat.Wins,
		&stat.Losses,
		&stat.Points,
		&stat.PointsScored,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query stat by id: %w", err)
	}

	return &stat, nil
}
