package database

import (
	"errors"
)

// Query used to check if a user is member of a group
var queryMemberGroup = `SELECT UserId FROM user_group WHERE UserId = ? AND GroupId = ?`

// Function used to check id a user is member of a group
func (db *appdbimpl) CheckMember(userId int, groupId int) (bool, error) {
	var id int
	err := db.c.QueryRow(queryMemberGroup, userId, groupId).Scan(&id)
	if err != nil {
		return false, errors.New("User is not a member of the group")
	}
	return id != 0, nil
}
