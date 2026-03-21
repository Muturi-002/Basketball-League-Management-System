package database

// User maps the fields from the Users table (users.png).
type User struct {
    UserID       int64  `json:"userId"`
    FirstName    string `json:"firstName"`
    LastName     string `json:"lastName"`
    EmailAddress string `json:"emailAddress"`
    Password     string `json:"password"`
    FavTeamID    int64  `json:"favTeamId"`
}

// TODO: Implement data access operations for user authentication and retrieval.
func ListUsers() ([]User, error) {
    return nil, nil
}
