package database

// Query for update the username of the user
var queryUpdateGroupName = "UPDATE group_t SET GroupName = ? WHERE GroupId = ?;"

func (db *appdbimpl) SetGroupName(GroupId int, newName string) error {
	// Update the username of the user
	_, err := db.c.Exec(queryUpdateGroupName, newName, GroupId)
	if err != nil {
		return err
	}
	return nil
}
