package database

// Query used to get all conversation with a group
var queryGetConversationsByUserGroupId = `SELECT ConversationId FROM conversation WHERE GroupId = ? AND UserId = ?`

func (db *appdbimpl) GetConversationsByUserGroup(groupId int, userId int) (int, error) {
	var convId int
	// Get all conversations of a group from the db
	err := db.c.QueryRow(queryGetConversationsByUserGroupId, groupId, userId).Scan(&convId)
	if err != nil {
		return 0, err
	}

	return convId, nil
}
