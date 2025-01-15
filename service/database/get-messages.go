package database

import (
	"wasa.project/service/api/structs"
)

// Query used to get the messages of a conversation from the database
var queryGetMessages = `SELECT MessageId, Message, Status, SenderUserId, SendTime, COALESCE(Photo, "") FROM message WHERE ConversationId = ?`

func (db *appdbimpl) GetMessages(convId int) ([]structs.Message, error) {
	// Create a new slice of messages
	var messages []structs.Message

	// Execute the query to get the messages of a conversation
	rows, err := db.c.Query(queryGetMessages, convId)
	if err != nil {
		return nil, err
	}
	defer func() { rows.Close() }()

	// Iterate
	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		// Create a new message
		var message structs.Message
		// Scan the values of the message
		err := rows.Scan(&message.MessageId, &message.Text, &message.Status, &message.SenderUserId, &message.SendTime, &message.Photo)
		if err != nil {
			return nil, err
		}
		// Append the message to the slice
		messages = append(messages, message)
	}
	return messages, nil
}
