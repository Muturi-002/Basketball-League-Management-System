package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type User struct {
	UserID       string `json:"-"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
	FavTeamID    string `json:"favTeamId,omitempty"`
}

type AuthUser struct {
	UserID       string `json:"-"`
	Username     string `json:"username"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	EmailAddress string `json:"emailAddress"`
	FavTeamID    string `json:"favTeamId,omitempty"`
}

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
		var favTeamID sql.NullInt64
		if err := rows.Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.EmailAddress,
			&user.Password,
			&favTeamID,
		); err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}
		user.FavTeamID = nullInt64ToString(favTeamID)
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func GetUserByID(userID string) (*User, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("user get: %w", err)
	}
	userID = normalizeUserID(userID)

	var user User
	var favTeamID sql.NullInt64
	err = conn.QueryRow(`
		SELECT
			user_id,
			first_name,
			last_name,
			email_address,
			password,
			fav_team_id
		FROM users
		WHERE UPPER(user_id) = UPPER(:1)`, userID).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.EmailAddress,
		&user.Password,
		&favTeamID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	user.FavTeamID = nullInt64ToString(favTeamID)

	return &user, nil
}

func AuthenticateUser(identifier, password string) (*AuthUser, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("user auth: %w", err)
	}

	if err := ValidateLoginIdentifier(identifier); err != nil {
		return nil, err
	}
	if err := ValidatePassword(password, 8, 72); err != nil {
		return nil, err
	}

	var user User
	var favTeamID sql.NullInt64
	err = conn.QueryRow(`
		SELECT
			user_id,
			first_name,
			last_name,
			email_address,
			password,
			fav_team_id
		FROM users
		WHERE UPPER(user_id) = UPPER(:1)
		   OR LOWER(email_address) = LOWER(:2)`, identifier, identifier).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.EmailAddress,
		&user.Password,
		&favTeamID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user for auth: %w", err)
	}

	if !passwordHashMatches(user.Password, password) {
		return nil, nil
	}

	return &AuthUser{
		UserID:       user.UserID,
		Username:     user.UserID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		EmailAddress: user.EmailAddress,
		FavTeamID:    nullInt64ToString(favTeamID),
	}, nil
}

func CreateUser(user User) (*AuthUser, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}

	userID := normalizeUserID(user.UserID)
	firstName := strings.TrimSpace(user.FirstName)
	lastName := strings.TrimSpace(user.LastName)
	emailAddress := strings.ToLower(strings.TrimSpace(user.EmailAddress))
	password := user.Password
	favTeamID := strings.TrimSpace(user.FavTeamID)

	if err := ValidateUsername(userID, 3, 8); err != nil {
		return nil, err
	}
	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("first and last name are required")
	}
	if err := ValidateEmailAddress(emailAddress); err != nil {
		return nil, err
	}
	if err := ValidatePassword(password, 8, 72); err != nil {
		return nil, err
	}

	var duplicateCount int64
	err = conn.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE UPPER(user_id) = UPPER(:1)
		   OR LOWER(email_address) = LOWER(:2)`, userID, emailAddress).Scan(&duplicateCount)
	if err != nil {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if duplicateCount > 0 {
		return nil, fmt.Errorf("user already exists")
	}

	hashedPassword := hashPassword(password)
	tx, err := conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin user create transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var favTeamValue any
	if favTeamID != "" {
		teamID, parseErr := strconv.ParseInt(favTeamID, 10, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("favorite team id must be numeric")
		}
		favTeamValue = teamID
	}

	if _, err := tx.Exec(`
		INSERT INTO users (
			user_id,
			first_name,
			last_name,
			email_address,
			password,
			fav_team_id
		) VALUES (:1, :2, :3, :4, :5, :6)`,
		userID,
		firstName,
		lastName,
		emailAddress,
		hashedPassword,
		favTeamValue,
	); err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit user create: %w", err)
	}

	authUser := &AuthUser{
		UserID:       userID,
		Username:     userID,
		FirstName:    firstName,
		LastName:     lastName,
		EmailAddress: emailAddress,
		FavTeamID:    favTeamID,
	}

	go func() {
		if err := sendAccountCreatedEmail(authUser.EmailAddress, authUser.FirstName, authUser.Username); err != nil {
			fmt.Printf("account email failed for %s: %v\n", authUser.EmailAddress, err)
		}
	}()

	return authUser, nil
}

func normalizeUserID(userID string) string {
	return strings.ToUpper(strings.TrimSpace(userID))
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func passwordHashMatches(storedHash, password string) bool {
	return strings.EqualFold(strings.TrimSpace(storedHash), hashPassword(password))
}

func nullInt64ToString(value sql.NullInt64) string {
	if !value.Valid {
		return ""
	}
	return strconv.FormatInt(value.Int64, 10)
}
