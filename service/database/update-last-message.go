package database

// Query used to update the last message in the conversation table
var queryUpdateLastMsg = `UPDATE conversation SET LastMessageId = ? WHERE ConversationId = ?`

func (db *appdbimpl) UpdateLastMessage(convId int, msgId int) error {
	// Update the last message in the conversation table
	_, err := db.c.Exec(queryUpdateLastMsg, msgId, convId)
	if err != nil {
		return err
	}
	return nil
}
