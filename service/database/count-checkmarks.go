package database

// Query used to count member of a checkmarks
var queryCountCheckmarks = `SELECT COUNT(UserId) FROM checkmarks WHERE ConversationId = ? AND MessageId = ?`

// CountCheckmarks returns the number of member of a checkmarks
func (db *appdbimpl) CountCheckmarks(conversationID int, messageID int) (int, error) {
	var count int
	err := db.c.QueryRow(queryCountCheckmarks, conversationID, messageID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}