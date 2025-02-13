package	database

// Update status of the message
var queryUpdateStatus = `UPDATE message SET Status = ? WHERE ConversationId = ? AND MessageId = ?`

// UpdateStatus updates the status of the message
func (db *appdbimpl) UpdateStatus(conversationID int, messageID int, status string) error {
	_, err := db.c.Exec(queryUpdateStatus, status, conversationID, messageID)
	if err != nil {
		return err
	}
	return nil
}