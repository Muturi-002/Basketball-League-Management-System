package database

import "fmt"

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
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("users list: %w", err)
	}

	rows, err := conn.Query(`
        SELECT
            user_id,
            first_name,
            last_name,
            email_address,
            password,
            fav_team_id
        FROM users
        ORDER BY user_id`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.EmailAddress,
			&user.Password,
			&user.FavTeamID,
		); err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}
