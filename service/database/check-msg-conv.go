package database

// Query used to check if a message is part of a conversation
var queryCheckMessage = `SELECT MessageId FROM message WHERE ConversationId = ? AND SendUserId = ? AND MessageId = ?`

// Function used to check if a message is part of a conversation
func (db *appdbimpl) CheckMessageConv(msgId int, convId int, userId int) (bool, error) {
	_, err := db.c.Exec(queryCheckMessage, convId, userId, msgId)
	if err != nil {
		return false, err
	}
	return true, nil
}
