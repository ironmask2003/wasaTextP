package database

// Query used for search a conversation of a group
var queryGetConvGroup = `SELECT ConversationId FROM conversation WHERE GroupId = ?`

// Function used to get a conversation of a group
func (db *appdbimpl) GetConvGroup(groupId int) (int, error) {
	var convId int
	err := db.c.QueryRow(queryGetConvGroup, groupId).Scan(&convId)
	if err != nil {
		return 0, err
	}
	return convId, nil
}
