package database

// Query used to take groups where the user taked in input is part
var queryGetGroups = `SELECT GroupId FROM user_group WHERE UserId = ?`

func (db *appdbimpl) GetGroups(userId int) ([]Group, error) {
	// List of the groups
	var groups []Group

	return groups, nil
}
