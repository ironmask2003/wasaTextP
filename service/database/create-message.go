package database

// Query used to insert a new message in the message table
var queryAddtMessage = "INSERT INTO message (MessageId, Message, SendTime, Status, SenderUserId) VALUES (?, ?, ?, ?, ?);"

// Query for take the max id in the message table
var queryMaxMessageID = "SELECT MAX(UserId) FROM message"

func (db *appdbimpl) CreateMessage(message Message) (Message, error) {
	var msg Message
	msg.Text = message.Text

	// Getting the max id in the message table
	maxID, err := db.GetLastElem(queryMaxMessageID)
	if err != nil {
		return msg, err
	}

	// Setting the id of the new user
	msg.MessageId += maxID + 1

	return msg, nil
}
