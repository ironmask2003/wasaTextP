package database

var queryInsertCheckmarks = `INSERT INTO checkmarks (ConversationId, MessageId, UserId) VALUES (?, ?, ?)`

// InsertCheckmarks inserts checkmarks into the database
func (db *appdbimpl) InsertCheckmarks(conversationID int, messageID int, userID int) error {
	_, err := db.c.Exec(queryInsertCheckmarks, conversationID, messageID, userID)
	if err != nil {
		return err
	}
	return nil
}