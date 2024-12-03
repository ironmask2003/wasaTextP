package database

// Query used to update the last message in the conversation table
var queryUpdateLastMsg = `UPDATE conversation SET LastMessageId = ? WHERE ConversationId = ? AND UserId = ?`

func (db *appdbimpl) UpdateLastMessage(convId int, userId int, msgId int) error {
	// Update the last message in the conversation table
	_, err := db.c.Exec(queryUpdateLastMsg, msgId, convId, userId)
	if err != nil {
		return err
	}
	return nil
}
