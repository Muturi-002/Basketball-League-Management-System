package database

import (
	"database/sql"
	"fmt"
)

// Manager maps every column from the team_managers table.
// ManagerPhoto is an OCI Object Storage URL stored as VARCHAR(200) — not binary.
type Manager struct {
	ManagerID      int64   `json:"managerId"`
	FirstName      string  `json:"firstName"`
	MiddleName     *string `json:"middleName,omitempty"`
	LastName       string  `json:"lastName"`
	ManagerPhoto   *string `json:"managerPhoto,omitempty"`
	ManagerType    string  `json:"managerType"`
	TeamID         int64   `json:"teamId"`
	ManagerHistory *string `json:"managerHistory,omitempty"`
}

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
		var (
			middleName     sql.NullString
			managerPhoto   sql.NullString
			managerHistory sql.NullString
		)

		if err := rows.Scan(
			&manager.ManagerID,
			&manager.FirstName,
			&middleName,
			&manager.LastName,
			&managerPhoto,
			&manager.ManagerType,
			&manager.TeamID,
			&managerHistory,
		); err != nil {
			return nil, fmt.Errorf("scan manager row: %w", err)
		}

		manager.MiddleName = nullStringPtr(middleName)
		manager.ManagerPhoto = nullStringPtr(managerPhoto)
		manager.ManagerHistory = nullStringPtr(managerHistory)

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
	var (
		middleName     sql.NullString
		managerPhoto   sql.NullString
		managerHistory sql.NullString
	)

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
		&managerPhoto,
		&manager.ManagerType,
		&manager.TeamID,
		&managerHistory,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query manager by id: %w", err)
	}

	manager.MiddleName = nullStringPtr(middleName)
	manager.ManagerPhoto = nullStringPtr(managerPhoto)
	manager.ManagerHistory = nullStringPtr(managerHistory)

	return &manager, nil
}
