package database


type Club struct {
    ClubID              int64  `json:"clubId"`
    TeamName            string `json:"teamName"`
    TeamLogo            []byte `json:"teamLogo"`           // logo image data
    ClubManagerID       int64  `json:"clubManagerId"`
    ClubManagerPhoto    []byte `json:"clubManagerPhoto"`   // manager photo
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
    return nil, nil
}

func ListStadiums() ([]Stadium, error) {
    return nil, nil
}
