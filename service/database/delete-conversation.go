package database

var queryDeleteConversation = `DELETE FROM conversation WHERE ConversationId = ?;`

func (db *appdbimpl) DeleteConversation(convId int) error {
	// Delete the conversation
	_, err := db.c.Exec(queryDeleteConversation, convId)
	if err != nil {
		return err
	}
	return nil
}
