package database

import (
	"database/sql"
	"fmt"
	"log"
)

type Club struct {
	ClubID           int64   `json:"clubId"`
	TeamName         string  `json:"teamName"`
	TeamLogo         *string `json:"teamLogo,omitempty"`
	ClubManagerID    *int64  `json:"clubManagerId,omitempty"`
	ClubManagerPhoto *string `json:"clubManagerPhoto,omitempty"`
	ClubLocationID   *int64  `json:"clubLocationId,omitempty"`
	ClubHistory      *string `json:"clubHistory,omitempty"`
}

type Stadium struct {
	StadiumID       int64  `json:"stadiumId"`
	StadiumName     string `json:"stadiumName"`
	StadiumLocation string `json:"stadiumLocation"`
	AssociatedClub  string `json:"associatedClub"`
}

func ListClubs() ([]Club, error) {
	return queryList("clubs list", "query clubs", "iterate clubs", `
		SELECT
			club_id,
			team_name,
			team_logo,
			club_manager_id,
			club_manager_photo,
			club_location_id,
			club_history
		FROM clubs
		ORDER BY club_id`, func(rows *sql.Rows) (Club, error) {
		var club Club
		var (
			teamLogo         sql.NullString
			clubManagerID    sql.NullInt64
			clubManagerPhoto sql.NullString
			clubLocationID   sql.NullInt64
			clubHistory      sql.NullString
		)

		if err := rows.Scan(
			&club.ClubID,
			&club.TeamName,
			&teamLogo,
			&clubManagerID,
			&clubManagerPhoto,
			&clubLocationID,
			&clubHistory,
		); err != nil {
			log.Printf("scan club row error: %v", err)
			return Club{}, fmt.Errorf("scan club row: %w", err)
		}

		club.TeamLogo = nullStringPtr(teamLogo)
		club.ClubManagerID = nullInt64Ptr(clubManagerID)
		club.ClubManagerPhoto = nullStringPtr(clubManagerPhoto)
		club.ClubLocationID = nullInt64Ptr(clubLocationID)
		club.ClubHistory = nullStringPtr(clubHistory)

		return club, nil
	})
}

func GetClubByID(clubID int64) (*Club, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("club get: %w", err)
	}

	var club Club
	var (
		teamLogo         sql.NullString
		clubManagerID    sql.NullInt64
		clubManagerPhoto sql.NullString
		clubLocationID   sql.NullInt64
		clubHistory      sql.NullString
	)

	err = conn.QueryRow(`
		SELECT
			club_id,
			team_name,
			team_logo,
			club_manager_id,
			club_manager_photo,
			club_location_id,
			club_history
		FROM clubs
		WHERE club_id = :1`, clubID).Scan(
		&club.ClubID,
		&club.TeamName,
		&teamLogo,
		&clubManagerID,
		&clubManagerPhoto,
		&clubLocationID,
		&clubHistory,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("query club by id error: %v", err)
		return nil, fmt.Errorf("query club by id: %w", err)
	}

	club.TeamLogo = nullStringPtr(teamLogo)
	club.ClubManagerID = nullInt64Ptr(clubManagerID)
	club.ClubManagerPhoto = nullStringPtr(clubManagerPhoto)

	club.ClubLocationID = nullInt64Ptr(clubLocationID)
	club.ClubHistory = nullStringPtr(clubHistory)

	return &club, nil
}

func ListStadiums() ([]Stadium, error) {
	return queryList("stadiums list", "query stadiums", "iterate stadiums", `
		SELECT
			stadium_id,
			stadium_name,
			stadium_location,
			associated_club
		FROM stadiums
		ORDER BY stadium_id`, func(rows *sql.Rows) (Stadium, error) {
		var stadium Stadium
		if err := rows.Scan(
			&stadium.StadiumID,
			&stadium.StadiumName,
			&stadium.StadiumLocation,
			&stadium.AssociatedClub,
		); err != nil {
			return Stadium{}, fmt.Errorf("scan stadium row: %w", err)
		}
		return stadium, nil
	})
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
