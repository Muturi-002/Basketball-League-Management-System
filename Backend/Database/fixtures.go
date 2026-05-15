package database

import (
	"database/sql"
	"fmt"
	"time"
)

// Fixture maps the fields from the Fixtures table (fixtures.png).
type Fixture struct {
	FixtureID       int64     `json:"fixtureId"`
	FixtureDate     time.Time `json:"fixtureDate"`
	FixtureTime     time.Time `json:"fixtureTime"`
	HomeTeamID      int64     `json:"homeTeamId"`
	HomeTeamLogo    string    `json:"homeTeamLogo"`
	AwayTeamID      int64     `json:"awayTeamId"`
	AwayTeamLogo    string    `json:"awayTeamLogo"`
	FixtureLocation string    `json:"fixtureLocation"`
}

// TODO: Implement data access operations for fixtures.
func ListFixtures() ([]Fixture, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("fixtures list: %w", err)
	}

	rows, err := conn.Query(`
        SELECT
            fixture_id,
            fixture_date,
            fixture_time,
            home_team_id,
            home_team_logo,
            away_team_id,
            away_team_logo,
            fixture_location
        FROM fixtures
        ORDER BY fixture_date, fixture_time`)
	if err != nil {
		return nil, fmt.Errorf("query fixtures: %w", err)
	}
	defer rows.Close()

	fixtures := make([]Fixture, 0)
	for rows.Next() {
		var fixture Fixture
		if err := rows.Scan(
			&fixture.FixtureID,
			&fixture.FixtureDate,
			&fixture.FixtureTime,
			&fixture.HomeTeamID,
			&fixture.HomeTeamLogo,
			&fixture.AwayTeamID,
			&fixture.AwayTeamLogo,
			&fixture.FixtureLocation,
		); err != nil {
			return nil, fmt.Errorf("scan fixture row: %w", err)
		}
		fixtures = append(fixtures, fixture)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate fixtures: %w", err)
	}

	return fixtures, nil
}

func GetFixtureByID(fixtureID int64) (*Fixture, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("fixture get: %w", err)
	}

	var fixture Fixture
	err = conn.QueryRow(`
        SELECT
            fixture_id,
            fixture_date,
            fixture_time,
            home_team_id,
            home_team_logo,
            away_team_id,
            away_team_logo,
            fixture_location
        FROM fixtures
        WHERE fixture_id = :1`, fixtureID).Scan(
		&fixture.FixtureID,
		&fixture.FixtureDate,
		&fixture.FixtureTime,
		&fixture.HomeTeamID,
		&fixture.HomeTeamLogo,
		&fixture.AwayTeamID,
		&fixture.AwayTeamLogo,
		&fixture.FixtureLocation,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query fixture by id: %w", err)
	}

	return &fixture, nil
}
