package database

import (
	"database/sql"
	"fmt"
)

// Manager maps the fields from the Team Managers table (managers.png).
type Manager struct {
	ManagerID      int64   `json:"managerId"`
	FirstName      string  `json:"firstName"`
	MiddleName     *string `json:"middleName,omitempty"`
	LastName       string  `json:"lastName"`
	ManagerPhoto   []byte  `json:"managerPhoto"` // binary photo data
	ManagerType    string  `json:"managerType"`
	TeamID         int64   `json:"teamId"`
	ManagerHistory string  `json:"managerHistory"`
}

// TODO: Implement data access operations for managers.
func ListManagers() ([]Manager, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("managers list: %w", err)
	}

	rows, err := conn.Query(`
		SELECT
			manager_id,
			first_name,
			middle_name,
			last_name,
			manager_photo,
			manager_type,
			team_id,
			manager_history
		FROM team_managers
		ORDER BY manager_id`)
	if err != nil {
		return nil, fmt.Errorf("query managers: %w", err)
	}
	defer rows.Close()

	managers := make([]Manager, 0)
	for rows.Next() {
		var manager Manager
		var middleName sql.NullString

		if err := rows.Scan(
			&manager.ManagerID,
			&manager.FirstName,
			&middleName,
			&manager.LastName,
			&manager.ManagerPhoto,
			&manager.ManagerType,
			&manager.TeamID,
			&manager.ManagerHistory,
		); err != nil {
			return nil, fmt.Errorf("scan manager row: %w", err)
		}

		manager.MiddleName = nullStringPtr(middleName)
		managers = append(managers, manager)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate managers: %w", err)
	}

	return managers, nil
}

func GetManagerByID(managerID int64) (*Manager, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("manager get: %w", err)
	}

	var manager Manager
	var middleName sql.NullString

	err = conn.QueryRow(`
		SELECT
			manager_id,
			first_name,
			middle_name,
			last_name,
			manager_photo,
			manager_type,
			team_id,
			manager_history
		FROM team_managers
		WHERE manager_id = :1`, managerID).Scan(
		&manager.ManagerID,
		&manager.FirstName,
		&middleName,
		&manager.LastName,
		&manager.ManagerPhoto,
		&manager.ManagerType,
		&manager.TeamID,
		&manager.ManagerHistory,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query manager by id: %w", err)
	}

	manager.MiddleName = nullStringPtr(middleName)
	return &manager, nil
}
