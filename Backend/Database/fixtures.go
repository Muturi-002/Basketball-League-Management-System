package database

import "time"

// Fixture maps the fields from the Fixtures table (fixtures.png).
type Fixture struct {
    FixtureID       int64     `json:"fixtureId"`
    FixtureDate     time.Time `json:"fixtureDate"`
    FixtureTime     time.Time `json:"fixtureTime"`
    HomeTeamID      int64     `json:"homeTeamId"`
    HomeTeamLogo    []byte    `json:"homeTeamLogo"`
    AwayTeamID      int64     `json:"awayTeamId"`
    AwayTeamLogo    []byte    `json:"awayTeamLogo"`
    FixtureLocation string    `json:"fixtureLocation"`
}

// TODO: Implement data access operations for fixtures.
func ListFixtures() ([]Fixture, error) {
    return nil, nil
}
