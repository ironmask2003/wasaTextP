package database

// Query used to update status of message
var queryUpdateMessage = `UPDATE message SET Status = ? WHERE MessageId = ? AND ConversationId = ? AND Status = "Sended"`

func (db *appdbimpl) UpdateStatusMessage(msgId int, convId int) error {
	_, err := db.c.Exec(queryUpdateMessage, "Readed", msgId, convId)
	if err != nil {
		return err
	}
	return nil
}
