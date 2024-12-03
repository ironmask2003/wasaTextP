package database

// Query to take the message from a conversation
var queryDeleteMessage = `DELETE FROM message WHERE ConversationId = ? AND MessageId = ? AND SendUserId = ?`

func (db *appdbimpl) DeleteMessage(userId int, convId int, messageId int) error {
	// Delete the message
	_, err := db.c.Exec(queryDeleteMessage, convId, messageId, userId)
	if err != nil {
		return err
	}
	return nil
}
