package database

// Query to take the message from a conversation
var queryGetMessage = `SELECT Message FROM message WHERE ConversationId = ? AND MessageId = ? AND SendUserId = ?`

func (db *appdbimpl) GetMessage(userId int, convId int, messageId int) (Message, error) {
	var msg Message
	// Get all conversations of a group from the db
	err := db.c.QueryRow(queryGetMessage, convId, messageId, userId).Scan(&msg.Text)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
