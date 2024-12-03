package database

// Query used to check if a message is part of a conversation
var queryCheckMessage = `SELECT MessageId FROM message WHERE ConversationId = ? AND SendUserId = ? AND MessageId = ?`

func (db *appdbimpl) CheckMessageConv(msgId int, convId int, userId int) (bool, error) {
	var msg Message
	err := db.c.QueryRow(queryCheckMessage, convId, userId, msgId).Scan(&msg.MessageId)
	if err != nil {
		return false, err
	}
	return true, nil
}
