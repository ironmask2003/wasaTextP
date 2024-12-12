package database

import (
	"wasa.project/service/api/structs"
)

// Query used to find a message by its id in the database
var queryFindMessageById = `SELECT MessageId, Message, SenderUserId, SendTime FROM message WHERE MessageId = ? AND ConversationId = ?`

func (db *appdbimpl) GetMessageById(messageId int, convId int) (structs.Message, error) {
	var message structs.Message
	err := db.c.QueryRow(queryFindMessageById, messageId, convId).Scan(&message.MessageId, &message.Text, &message.SenderUserId, &message.SendTime)
	return message, err
}
