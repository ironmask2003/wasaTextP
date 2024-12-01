package database

// Query used to add the user and the group in the user_group table
var queryAddLink = "INSERT INTO user_group (GroupId, UserId) VALUES (?, ?)"

func (db *appdbimpl) AddUserGroup(UserId int, GroupId int) error {
	_, err := db.c.Exec(queryAddLink, GroupId, UserId)
	if err != nil {
		return err
	}
	return nil
}
