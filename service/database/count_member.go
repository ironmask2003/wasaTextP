package database

// Query used to count member of a conversation
var queryCountMember = `SELECT COUNT(UserId) FROM conversation_user WHERE ConversationId = ?`

// CountMember returns the number of member of a conversation
func (db *appdbimpl) CountMember(conversationID int) (int, error) {
	var count int
	err := db.c.QueryRow(queryCountMember, conversationID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}