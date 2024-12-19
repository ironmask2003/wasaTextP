package database

<<<<<<< HEAD
import (
	"database/sql"
)

=======
>>>>>>> 779b51a (Modified table and function)
// Query used to update the last message from a conversation
var queryLastMessage = `UPDATE conversation SET LastMessageId = ? WHERE ConversationId = ?`

// Function used to update the last message of a conversation
func (db *appdbimpl) UpdateLastMessage(messageId int, conversationId int) error {
<<<<<<< HEAD
	value := sql.NullInt64{Int64: int64(messageId), Valid: !(messageId == 0)}
	_, err := db.c.Exec(queryLastMessage, value, conversationId)
=======
	_, err := db.c.Exec(queryLastMessage, messageId, conversationId)
>>>>>>> 779b51a (Modified table and function)
	return err
}
