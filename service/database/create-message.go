package database

// Query used to insert a new message in the message table
var queryAddtMessage = "INSERT INTO message (MessageId, ConversationId, Message, SendTime, Status, SenderUserId) VALUES (?, ?, ?, ?, ?);"

// Query for take the max id in the message table
var queryMaxMessageID = "SELECT MAX(MessageId) FROM message WHERE ConversationId = ?;"

func (db *appdbimpl) CreateMessage(message Message) (Message, error) {
	var msg Message
	msg.Text = message.Text

	return msg, nil
}
