package database

// Stat is a placeholder structure for league statistics, to be refined
// once the concrete stats schema is available.
type Stat struct {
	ID   int64  `json:"id"`
	Note string `json:"note"`
}

// TODO: Replace this placeholder with real statistics fields and queries.
func ListStats() ([]Stat, error) {
	return nil, nil
}
