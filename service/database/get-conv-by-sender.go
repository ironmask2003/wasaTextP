package database

// Query used to get all conversation with a group
var queryGetConversationsBySender = `SELECT ConversationId FROM conversation WHERE SenderUserId = ? AND UserId = ?`

func (db *appdbimpl) GetConversationsBySender(senderId int, userId int) (int, error) {
	var convId int
	// Get all conversations of a group from the db
	err := db.c.QueryRow(queryGetConversationsBySender, senderId, userId).Scan(&convId)
	if err != nil {
		return 0, err
	}
	return convId, nil
}
