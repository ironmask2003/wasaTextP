package database

// Query used to update the last message from a conversation
var queryLastMessage = `UPDATE conversation SET LastMessageId = ? WHERE ConversationId = ?`

// Function used to update the last message of a conversation
func (db *appdbimpl) UpdateLastMessage(messageId int, conversationId int) error {
	_, err := db.c.Exec(queryLastMessage, messageId, conversationId)
	return err
}
