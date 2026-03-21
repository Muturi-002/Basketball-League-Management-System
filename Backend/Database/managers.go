package database

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
	return nil, nil
}
