package database

// Query used to check if a user is member of a group
var queryMemberGroup = `SELECT UserId, GroupId FROM user_group WHERE UserId = ? AND GroupId = ?`

// Function used to check id a user is member of a group
func (db *appdbimpl) CheckMember(userId int, groupId int) (bool, error) {
	_, err := db.c.Exec(queryMemberGroup, userId, groupId)
	if err != nil {
		return false, err
	}
	return true, nil
}
