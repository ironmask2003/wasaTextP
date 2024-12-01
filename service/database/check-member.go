package database

// Query used to check if a user is member of a group
var queryMemberGroup = `SELECT UserId, GroupId FROM user_group WHERE UserId = ? AND GroupId = ?`

// Function used to check id a user is member of a group
func (db *appdbimpl) CheckMember(userId int, groupId int) (bool, error) {
	var user int
	var group int
	err := db.c.QueryRow(queryMemberGroup, userId, groupId).Scan(&user, &group)
	if err != nil {
		return false, err
	}
	return true, nil
}
