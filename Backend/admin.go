package main

import (
	"fmt"
	"net/http"
	"strings"

	db "blms/Database"
)

func nullableStringFromPtr(value *string) interface{} {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func adminPlayersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var player db.Player
	if err := decodeJSON(r, &player); err != nil {
		writeError(w, http.StatusBadRequest, "invalid player payload")
		return
	}

	if err := createPlayer(player); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create player")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "player created"})
}

func adminPlayerByIDHandler(w http.ResponseWriter, r *http.Request) {
	playerID, ok := parseIDFromPath(r.URL.Path, "/api/admin/players/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid player id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		player, err := db.GetPlayerByID(playerID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load player")
			return
		}
		if player == nil {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		writeJSON(w, http.StatusOK, player)
	case http.MethodPut:
		existing, err := db.GetPlayerByID(playerID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load player")
			return
		}
		if existing == nil {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		var player db.Player
		if err := decodeJSON(r, &player); err != nil {
			writeError(w, http.StatusBadRequest, "invalid player payload")
			return
		}
		player.PlayerID = playerID
		if err := updatePlayer(player); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update player")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "player updated"})
	case http.MethodDelete:
		if err := deletePlayer(playerID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete player")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "player deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func adminTeamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var club db.Club
	if err := decodeJSON(r, &club); err != nil {
		writeError(w, http.StatusBadRequest, "invalid team payload")
		return
	}

	if err := createClub(club); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create team")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "team created"})
}

func adminTeamByIDHandler(w http.ResponseWriter, r *http.Request) {
	clubID, ok := parseIDFromPath(r.URL.Path, "/api/admin/teams/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid team id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		club, err := db.GetClubByID(clubID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load team")
			return
		}
		if club == nil {
			writeError(w, http.StatusNotFound, "team not found")
			return
		}
		writeJSON(w, http.StatusOK, club)
	case http.MethodPut:
		existing, err := db.GetClubByID(clubID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load team")
			return
		}
		if existing == nil {
			writeError(w, http.StatusNotFound, "team not found")
			return
		}
		var club db.Club
		if err := decodeJSON(r, &club); err != nil {
			writeError(w, http.StatusBadRequest, "invalid team payload")
			return
		}
		club.ClubID = clubID
		if err := updateClub(club); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update team")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "team updated"})
	case http.MethodDelete:
		if err := deleteClub(clubID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete team")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "team deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func adminStadiumsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var stadium db.Stadium
	if err := decodeJSON(r, &stadium); err != nil {
		writeError(w, http.StatusBadRequest, "invalid stadium payload")
		return
	}

	if err := createStadium(stadium); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create stadium")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "stadium created"})
}

func adminStadiumByIDHandler(w http.ResponseWriter, r *http.Request) {
	stadiumID, ok := parseIDFromPath(r.URL.Path, "/api/admin/stadiums/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid stadium id")
		return
	}

	switch r.Method {
	case http.MethodPut:
		var stadium db.Stadium
		if err := decodeJSON(r, &stadium); err != nil {
			writeError(w, http.StatusBadRequest, "invalid stadium payload")
			return
		}
		stadium.StadiumID = stadiumID
		if err := updateStadium(stadium); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update stadium")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "stadium updated"})
	case http.MethodDelete:
		if err := deleteStadium(stadiumID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete stadium")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "stadium deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func adminManagersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var manager db.Manager
	if err := decodeJSON(r, &manager); err != nil {
		writeError(w, http.StatusBadRequest, "invalid manager payload")
		return
	}

	if err := createManager(manager); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create manager")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "manager created"})
}

func adminManagerByIDHandler(w http.ResponseWriter, r *http.Request) {
	managerID, ok := parseIDFromPath(r.URL.Path, "/api/admin/managers/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid manager id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		manager, err := db.GetManagerByID(managerID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load manager")
			return
		}
		if manager == nil {
			writeError(w, http.StatusNotFound, "manager not found")
			return
		}
		writeJSON(w, http.StatusOK, manager)
	case http.MethodPut:
		existing, err := db.GetManagerByID(managerID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load manager")
			return
		}
		if existing == nil {
			writeError(w, http.StatusNotFound, "manager not found")
			return
		}
		var manager db.Manager
		if err := decodeJSON(r, &manager); err != nil {
			writeError(w, http.StatusBadRequest, "invalid manager payload")
			return
		}
		manager.ManagerID = managerID
		if err := updateManager(manager); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update manager")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "manager updated"})
	case http.MethodDelete:
		if err := deleteManager(managerID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete manager")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "manager deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func adminStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var stat db.Stat
	if err := decodeJSON(r, &stat); err != nil {
		writeError(w, http.StatusBadRequest, "invalid stat payload")
		return
	}

	if err := createStat(stat); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create stat")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "stat created"})
}

func adminStatByIDHandler(w http.ResponseWriter, r *http.Request) {
	statID, ok := parseIDFromPath(r.URL.Path, "/api/admin/stats/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid stat id")
		return
	}

	switch r.Method {
	case http.MethodPut:
		var stat db.Stat
		if err := decodeJSON(r, &stat); err != nil {
			writeError(w, http.StatusBadRequest, "invalid stat payload")
			return
		}
		stat.StatID = statID
		if err := updateStat(stat); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update stat")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "stat updated"})
	case http.MethodDelete:
		if err := deleteStat(statID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete stat")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "stat deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func adminFixturesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var fixture db.Fixture
	if err := decodeJSON(r, &fixture); err != nil {
		writeError(w, http.StatusBadRequest, "invalid fixture payload")
		return
	}

	if err := createFixture(fixture); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create fixture")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "fixture created"})
}

func createStat(stat db.Stat) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("stat create: %w", err)
	}

	if stat.StatID > 0 {
		_, err = conn.Exec(`
			INSERT INTO team_statistics (
				team_stat_id, team_id, games_played, wins, losses, league_points,
				points_scored
			) VALUES (:1, :2, :3, :4, :5, :6, :7)`,
			stat.StatID, stat.TeamID, stat.GamesPlayed, stat.Wins, stat.Losses,
			stat.Points, stat.PointsScored,
		)
	} else {
		_, err = conn.Exec(`
			INSERT INTO team_statistics (
				team_id, games_played, wins, losses, league_points,
				points_scored
			) VALUES (:1, :2, :3, :4, :5, :6)`,
			stat.TeamID, stat.GamesPlayed, stat.Wins, stat.Losses,
			stat.Points, stat.PointsScored,
		)
	}
	if err != nil {
		return fmt.Errorf("insert stat: %w", err)
	}
	return nil
}

func updateStat(stat db.Stat) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("stat update: %w", err)
	}
	if stat.StatID <= 0 {
		return fmt.Errorf("stat update: invalid statId")
	}
	_, err = conn.Exec(`
		UPDATE team_statistics
		SET
			team_id = :1,
			games_played = :2,
			wins = :3,
			losses = :4,
			league_points = :5,
			points_scored = :6
		WHERE team_stat_id = :7`,
		stat.TeamID, stat.GamesPlayed, stat.Wins, stat.Losses,
		stat.Points, stat.PointsScored, stat.StatID,
	)
	if err != nil {
		return fmt.Errorf("update stat: %w", err)
	}
	return nil
}

func deleteStat(statID int64) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("stat delete: %w", err)
	}
	_, err = conn.Exec(`DELETE FROM team_statistics WHERE team_stat_id = :1`, statID)
	if err != nil {
		return fmt.Errorf("delete stat: %w", err)
	}
	return nil
}

func createFixture(fixture db.Fixture) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("fixture create: %w", err)
	}

	var homeScore interface{}
	var awayScore interface{}
	if fixture.HomeScore != nil {
		homeScore = *fixture.HomeScore
	}
	if fixture.AwayScore != nil {
		awayScore = *fixture.AwayScore
	}

	if fixture.FixtureID > 0 {
		_, err = conn.Exec(`
			INSERT INTO fixtures (
				fixture_id, fixture_date, fixture_time, home_team_id, home_team_logo,
				away_team_id, away_team_logo, fixture_location, fixture_status, home_score, away_score
			) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11)`,
			fixture.FixtureID, fixture.FixtureDate, fixture.FixtureTime, fixture.HomeTeamID, fixture.HomeTeamLogo,
			fixture.AwayTeamID, fixture.AwayTeamLogo, fixture.FixtureLocation, fixture.FixtureStatus, homeScore, awayScore,
		)
	} else {
		_, err = conn.Exec(`
			INSERT INTO fixtures (
				fixture_date, fixture_time, home_team_id, home_team_logo,
				away_team_id, away_team_logo, fixture_location, fixture_status, home_score, away_score
			) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)`,
			fixture.FixtureDate, fixture.FixtureTime, fixture.HomeTeamID, fixture.HomeTeamLogo,
			fixture.AwayTeamID, fixture.AwayTeamLogo, fixture.FixtureLocation, fixture.FixtureStatus, homeScore, awayScore,
		)
	}
	if err != nil {
		return fmt.Errorf("insert fixture: %w", err)
	}
	return nil
}

func adminFixtureByIDHandler(w http.ResponseWriter, r *http.Request) {
	fixtureID, ok := parseIDFromPath(r.URL.Path, "/api/admin/fixtures/")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid fixture id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		fixture, err := db.GetFixtureByID(fixtureID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load fixture")
			return
		}
		if fixture == nil {
			writeError(w, http.StatusNotFound, "fixture not found")
			return
		}
		writeJSON(w, http.StatusOK, fixture)
	case http.MethodPut:
		existing, err := db.GetFixtureByID(fixtureID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load fixture")
			return
		}
		if existing == nil {
			writeError(w, http.StatusNotFound, "fixture not found")
			return
		}
		var fixture db.Fixture
		if err := decodeJSON(r, &fixture); err != nil {
			writeError(w, http.StatusBadRequest, "invalid fixture payload")
			return
		}
		fixture.FixtureID = fixtureID
		if fixture.HomeTeamLogo == "" {
			fixture.HomeTeamLogo = existing.HomeTeamLogo
		}
		if fixture.AwayTeamLogo == "" {
			fixture.AwayTeamLogo = existing.AwayTeamLogo
		}
		if fixture.FixtureLocation == "" {
			fixture.FixtureLocation = existing.FixtureLocation
		}
		if err := updateFixture(fixture); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update fixture")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "fixture updated"})
	case http.MethodDelete:
		if err := deleteFixture(fixtureID); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to delete fixture")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "fixture deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func createPlayer(player db.Player) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("player create: %w", err)
	}

	middleName := nullableStringFromPtr(player.MiddleName)
	if player.PlayerID > 0 {
		_, err = conn.Exec(`
            INSERT INTO players (
                player_id,
                first_name,
                middle_name,
                last_name,
                player_photo,
                date_of_birth,
                age,
                team_id,
                player_history
            ) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)`,
			player.PlayerID,
			player.FirstName,
			middleName,
			player.LastName,
			player.PlayerPhoto,
			player.DateOfBirth,
			player.Age,
			player.TeamID,
			player.PlayerHistory,
		)
	} else {
		_, err = conn.Exec(`
            INSERT INTO players (
                first_name,
                middle_name,
                last_name,
                player_photo,
                date_of_birth,
                age,
                team_id,
                player_history
            ) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`,
			player.FirstName,
			middleName,
			player.LastName,
			player.PlayerPhoto,
			player.DateOfBirth,
			player.Age,
			player.TeamID,
			player.PlayerHistory,
		)
	}
	if err != nil {
		return fmt.Errorf("insert player: %w", err)
	}
	return nil
}
func updatePlayer(player db.Player) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("player update: %w", err)
	}
	if player.PlayerID <= 0 {
		return fmt.Errorf("player update: invalid playerId")
	}

	middleName := nullableStringFromPtr(player.MiddleName)
	_, err = conn.Exec(`
        UPDATE players
        SET
            first_name = :1,
            middle_name = :2,
            last_name = :3,
            player_photo = :4,
            date_of_birth = :5,
            age = :6,
            team_id = :7,
            player_history = :8
        WHERE player_id = :9`,
		player.FirstName,
		middleName,
		player.LastName,
		player.PlayerPhoto,
		player.DateOfBirth,
		player.Age,
		player.TeamID,
		player.PlayerHistory,
		player.PlayerID,
	)
	if err != nil {
		return fmt.Errorf("update player: %w", err)
	}
	return nil
}
func deletePlayer(playerID int64) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("player delete: %w", err)
	}
	_, err = conn.Exec(`DELETE FROM players WHERE player_id = :1`, playerID)
	if err != nil {
		return fmt.Errorf("delete player: %w", err)
	}
	return nil
}

func createClub(club db.Club) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("club create: %w", err)
	}

	if club.ClubID > 0 {
		_, err = conn.Exec(`
			INSERT INTO clubs (
				club_id,
				team_name,
				team_logo,
				club_manager_id,
				club_manager_photo,
				club_location_id,
				club_history
			) VALUES (:1, :2, :3, :4, :5, :6, :7)`,
			club.ClubID,
			club.TeamName,
			club.TeamLogo,
			club.ClubManagerID,
			club.ClubManagerPhoto,
			club.ClubLocationID,
			club.ClubHistory,
		)
	} else {
		_, err = conn.Exec(`
			INSERT INTO clubs (
				team_name, team_logo, club_manager_id, club_manager_photo,
				club_location_id, club_history
			) VALUES (:1, :2, :3, :4, :5, :6)`,
			club.TeamName, club.TeamLogo, club.ClubManagerID,
			club.ClubManagerPhoto, club.ClubLocationID, club.ClubHistory,
		)
	}
	if err != nil {
		return fmt.Errorf("insert club: %w", err)
	}
	return nil
}

func updateClub(club db.Club) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("club update: %w", err)
	}
	if club.ClubID <= 0 {
		return fmt.Errorf("club update: invalid clubId")
	}

	_, err = conn.Exec(`
        UPDATE clubs
		SET
			team_name = :1,
			team_logo = :2,
			club_manager_id = :3,
			club_manager_photo = :4,
			club_location_id = :5,
			club_history = :6
		WHERE club_id = :7`,
		club.TeamName, club.TeamLogo, club.ClubManagerID, club.ClubManagerPhoto,
		club.ClubLocationID, club.ClubHistory, club.ClubID,
	)
	if err != nil {
		return fmt.Errorf("update club: %w", err)
	}
	return nil
}

func deleteClub(clubID int64) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("club delete: %w", err)
	}
	_, err = conn.Exec(`DELETE FROM clubs WHERE club_id = :1`, clubID)
	if err != nil {
		return fmt.Errorf("delete club: %w", err)
	}
	return nil
}

func createStadium(stadium db.Stadium) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("stadium create: %w", err)
	}

	if stadium.StadiumID > 0 {
		_, err = conn.Exec(`
            INSERT INTO stadiums (
                stadium_id,
                stadium_name,
                stadium_location,
                associated_club
            ) VALUES (:1, :2, :3, :4)`,
			stadium.StadiumID,
			stadium.StadiumName,
			stadium.StadiumLocation,
			stadium.AssociatedClub,
		)
	} else {
		_, err = conn.Exec(`
            INSERT INTO stadiums (
                stadium_name,
                stadium_location,
                associated_club
            ) VALUES (:1, :2, :3)`,
			stadium.StadiumName,
			stadium.StadiumLocation,
			stadium.AssociatedClub,
		)
	}
	if err != nil {
		return fmt.Errorf("insert stadium: %w", err)
	}
	return nil
}

func updateStadium(stadium db.Stadium) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("stadium update: %w", err)
	}
	if stadium.StadiumID <= 0 {
		return fmt.Errorf("stadium update: invalid stadiumId")
	}

	_, err = conn.Exec(`
        UPDATE stadiums
        SET
            stadium_name = :1,
            stadium_location = :2,
            associated_club = :3
        WHERE stadium_id = :4`,
		stadium.StadiumName,
		stadium.StadiumLocation,
		stadium.AssociatedClub,
		stadium.StadiumID,
	)
	if err != nil {
		return fmt.Errorf("update stadium: %w", err)
	}
	return nil
}

func deleteStadium(stadiumID int64) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("stadium delete: %w", err)
	}
	_, err = conn.Exec(`DELETE FROM stadiums WHERE stadium_id = :1`, stadiumID)
	if err != nil {
		return fmt.Errorf("delete stadium: %w", err)
	}
	return nil
}

func createManager(manager db.Manager) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("manager create: %w", err)
	}
	middleName := nullableStringFromPtr(manager.MiddleName)
	if manager.ManagerID > 0 {
		_, err = conn.Exec(`
			INSERT INTO team_managers (
				manager_id, first_name, middle_name, last_name,
				manager_photo, manager_type, team_id, manager_history
			) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`,
			manager.ManagerID, manager.FirstName, middleName, manager.LastName,
			manager.ManagerPhoto, manager.ManagerType, manager.TeamID,
			manager.ManagerHistory,
		)
	} else {
		_, err = conn.Exec(`
			INSERT INTO team_managers (
				first_name, middle_name, last_name, manager_photo, manager_type,
				team_id, manager_history
			) VALUES (:1, :2, :3, :4, :5, :6, :7)`,
			manager.FirstName, middleName, manager.LastName,
			manager.ManagerPhoto, manager.ManagerType, manager.TeamID,
			manager.ManagerHistory,
		)
	}
	if err != nil {
		return fmt.Errorf("insert manager: %w", err)
	}
	return nil
}

func updateManager(manager db.Manager) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("manager update: %w", err)
	}
	if manager.ManagerID <= 0 {
		return fmt.Errorf("manager update: invalid managerId")
	}

	middleName := nullableStringFromPtr(manager.MiddleName)
	_, err = conn.Exec(`
		UPDATE team_managers
		SET
			first_name = :1,
			middle_name = :2,
			last_name = :3,
			manager_photo = :4,
			manager_type = :5,
			team_id = :6,
			manager_history = :7
		WHERE manager_id = :8`,
		manager.FirstName, middleName, manager.LastName, manager.ManagerPhoto,
		manager.ManagerType, manager.TeamID, manager.ManagerHistory,
		manager.ManagerID,
	)
	if err != nil {
		return fmt.Errorf("update manager: %w", err)
	}
	return nil
}

func deleteManager(managerID int64) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("manager delete: %w", err)
	}
	_, err = conn.Exec(`DELETE FROM team_managers WHERE manager_id = :1`, managerID)
	if err != nil {
		return fmt.Errorf("delete manager: %w", err)
	}
	return nil
}

func updateFixture(fixture db.Fixture) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("fixture update: %w", err)
	}
	if fixture.FixtureID <= 0 {
		return fmt.Errorf("fixture update: invalid fixtureId")
	}

	_, err = conn.Exec(`
        UPDATE fixtures
        SET
            fixture_date = :1,
            fixture_time = :2,
            home_team_id = :3,
            home_team_logo = :4,
            away_team_id = :5,
            away_team_logo = :6,
            fixture_location = :7
        WHERE fixture_id = :8`,
		fixture.FixtureDate, fixture.FixtureTime, fixture.HomeTeamID, fixture.HomeTeamLogo,
		fixture.AwayTeamID, fixture.AwayTeamLogo, fixture.FixtureLocation,
		fixture.FixtureID,
	)
	if err != nil {
		return fmt.Errorf("update fixture: %w", err)
	}
	return nil
}

func deleteFixture(fixtureID int64) error {
	conn, err := db.EnsureConnected()
	if err != nil {
		return fmt.Errorf("fixture delete: %w", err)
	}
	_, err = conn.Exec(`DELETE FROM fixtures WHERE fixture_id = :1`, fixtureID)
	if err != nil {
		return fmt.Errorf("delete fixture: %w", err)
	}
	return nil
}
