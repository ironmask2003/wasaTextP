package database

// Query for update the group id
var queryUpdateGroupId = `UPDATE conversation SET GroupId = ? WHERE ConversationId = ?`

// Function used to update the group id of a conversation
func (db *appdbimpl) UpdateGroupId(groupId int, conversationId int) error {
	_, err := db.c.Exec(queryUpdateGroupId, groupId, conversationId)
	return err
}
