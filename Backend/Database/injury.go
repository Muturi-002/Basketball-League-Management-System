package database

import "time"

// Injury maps the fields from the Injuries table (injury.png).
type Injury struct {
    InjuryID               int64  `json:"injuryId"`
    InjuryName             string `json:"injuryName"`
    InjuryType             string `json:"injuryType"`
    InjuryDescription      string `json:"injuryDescription"`
    InjuryResources        string `json:"injuryResources"`
    InjuryVideoDescription string `json:"injuryVideoDescription"`
}

// InjuredPlayer maps the fields from the Injured Players table (injured-players.png).
type InjuredPlayer struct {
    InjuredPlayerID      int64     `json:"injuredPlayerId"`
    PlayerID             int64     `json:"playerId"`
    InjuryID             int64     `json:"injuryId"`
    InjuryName           string    `json:"injuryName"`
    ExpectedTimeOfReturn time.Time `json:"expectedTimeOfReturn"`
}

// TODO: Implement data access operations for injuries and injured players.
func ListInjuries() ([]Injury, error) {
    return nil, nil
}

func ListInjuredPlayers() ([]InjuredPlayer, error) {
    return nil, nil
}
