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
	FixtureStatus   string    `json:"fixtureStatus"`
	HomeScore       *int64    `json:"homeScore,omitempty"`
	AwayScore       *int64    `json:"awayScore,omitempty"`
}

func ListFixtures() ([]Fixture, error) {
	return queryList("fixtures list", "query fixtures", "iterate fixtures", `
		SELECT
			fixture_id,
			fixture_date,
			fixture_time,
			home_team_id,
			home_team_logo,
			away_team_id,
			away_team_logo,
			fixture_location,
			fixture_status,
			home_score,
			away_score
		FROM fixtures
		ORDER BY fixture_date, fixture_time`, func(rows *sql.Rows) (Fixture, error) {
		var fixture Fixture
		var (
			fixtureStatus sql.NullString
			homeScore     sql.NullInt64
			awayScore     sql.NullInt64
		)
		if err := rows.Scan(
			&fixture.FixtureID,
			&fixture.FixtureDate,
			&fixture.FixtureTime,
			&fixture.HomeTeamID,
			&fixture.HomeTeamLogo,
			&fixture.AwayTeamID,
			&fixture.AwayTeamLogo,
			&fixture.FixtureLocation,
			&fixtureStatus,
			&homeScore,
			&awayScore,
		); err != nil {
			return Fixture{}, fmt.Errorf("scan fixture row: %w", err)
		}
		fixture.FixtureStatus = fixtureStatus.String
		fixture.HomeScore = nullInt64Ptr(homeScore)
		fixture.AwayScore = nullInt64Ptr(awayScore)
		return fixture, nil
	})
}

func GetFixtureByID(fixtureID int64) (*Fixture, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("fixture get: %w", err)
	}

	var fixture Fixture
	var (
		fixtureStatus sql.NullString
		homeScore     sql.NullInt64
		awayScore     sql.NullInt64
	)
	err = conn.QueryRow(`
		SELECT
			fixture_id,
			fixture_date,
			fixture_time,
			home_team_id,
			home_team_logo,
			away_team_id,
			away_team_logo,
			fixture_location,
			fixture_status,
			home_score,
			away_score
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
		&fixtureStatus,
		&homeScore,
		&awayScore,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query fixture by id: %w", err)
	}

	fixture.FixtureStatus = fixtureStatus.String
	fixture.HomeScore = nullInt64Ptr(homeScore)
	fixture.AwayScore = nullInt64Ptr(awayScore)

	return &fixture, nil
}
