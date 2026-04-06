package database

import (
	"database/sql"
	"fmt"
)

type Club struct {
	ClubID              int64  `json:"clubId"`
	TeamName            string `json:"teamName"`
	TeamLogo            []byte `json:"teamLogo"` // logo image data
	ClubManagerID       int64  `json:"clubManagerId"`
	ClubManagerPhoto    []byte `json:"clubManagerPhoto"` // manager photo
	AssClubManager      string `json:"assClubManager"`
	AssClubManagerPhoto []byte `json:"assClubManagerPhoto"` // assistant manager photo
	ClubLocationID      int64  `json:"clubLocationId"`
	ClubHistory         string `json:"clubHistory"`
}

// Stadium maps the fields from the Stadium/Location table (location.png).
type Stadium struct {
	StadiumID       int64  `json:"stadiumId"`
	StadiumName     string `json:"stadiumName"`
	StadiumLocation string `json:"stadiumLocation"`
	AssociatedClub  string `json:"associatedClub"`
}

// TODO: Implement data access for both clubs and stadiums.
func ListClubs() ([]Club, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("clubs list: %w", err)
	}

	rows, err := conn.Query(`
        SELECT
            club_id,
            team_name,
            team_logo,
            club_manager_id,
            club_manager_photo,
            ass_club_manager,
            ass_club_manager_photo,
            club_location_id,
            club_history
        FROM clubs
        ORDER BY club_id`)
	if err != nil {
		return nil, fmt.Errorf("query clubs: %w", err)
	}
	defer rows.Close()

	clubs := make([]Club, 0)
	for rows.Next() {
		var club Club
		if err := rows.Scan(
			&club.ClubID,
			&club.TeamName,
			&club.TeamLogo,
			&club.ClubManagerID,
			&club.ClubManagerPhoto,
			&club.AssClubManager,
			&club.AssClubManagerPhoto,
			&club.ClubLocationID,
			&club.ClubHistory,
		); err != nil {
			return nil, fmt.Errorf("scan club row: %w", err)
		}
		clubs = append(clubs, club)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate clubs: %w", err)
	}

	return clubs, nil
}

func ListStadiums() ([]Stadium, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("stadiums list: %w", err)
	}

	rows, err := conn.Query(`
        SELECT
            stadium_id,
            stadium_name,
            stadium_location,
            associated_club
        FROM stadiums
        ORDER BY stadium_id`)
	if err != nil {
		return nil, fmt.Errorf("query stadiums: %w", err)
	}
	defer rows.Close()

	stadiums := make([]Stadium, 0)
	for rows.Next() {
		var stadium Stadium
		if err := rows.Scan(
			&stadium.StadiumID,
			&stadium.StadiumName,
			&stadium.StadiumLocation,
			&stadium.AssociatedClub,
		); err != nil {
			return nil, fmt.Errorf("scan stadium row: %w", err)
		}
		stadiums = append(stadiums, stadium)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stadiums: %w", err)
	}

	return stadiums, nil
}

func GetClubByID(clubID int64) (*Club, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("club get: %w", err)
	}

	var club Club
	err = conn.QueryRow(`
        SELECT
            club_id,
            team_name,
            team_logo,
            club_manager_id,
            club_manager_photo,
            ass_club_manager,
            ass_club_manager_photo,
            club_location_id,
            club_history
        FROM clubs
        WHERE club_id = :1`, clubID).Scan(
		&club.ClubID,
		&club.TeamName,
		&club.TeamLogo,
		&club.ClubManagerID,
		&club.ClubManagerPhoto,
		&club.AssClubManager,
		&club.AssClubManagerPhoto,
		&club.ClubLocationID,
		&club.ClubHistory,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query club by id: %w", err)
	}

	return &club, nil
}

func GetStadiumByID(stadiumID int64) (*Stadium, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("stadium get: %w", err)
	}

	var stadium Stadium
	err = conn.QueryRow(`
        SELECT
            stadium_id,
            stadium_name,
            stadium_location,
            associated_club
        FROM stadiums
        WHERE stadium_id = :1`, stadiumID).Scan(
		&stadium.StadiumID,
		&stadium.StadiumName,
		&stadium.StadiumLocation,
		&stadium.AssociatedClub,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query stadium by id: %w", err)
	}

	return &stadium, nil
}
