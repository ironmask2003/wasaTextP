package database

import (
	"database/sql"
	"errors"
	"wasa.project/service/api/structs"
)

// Query used to add the message in the database
var queryAddMessage = `INSERT INTO message (MessageId, Message, Status, ConversationId, SenderUserId, Photo) VALUES (?, ?, ?, ?, ?, ?)`

// Query used to get the last id in the message table
var queryGetLastIdMessage = `SELECT MAX(MessageId) FROM message WHERE ConversationId = ?`

// Function used to get the last id in the conversation table
func (db *appdbimpl) GetMaxMessageId(convId int) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryGetLastIdMessage, convId)
	if err != nil {
		return 0, err
	}

	var maxID int
	for row.Next() {
		if row.Err() != nil {
			return 0, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}

	return maxID, nil
}

func (db *appdbimpl) CreateMessage(msg structs.Message) (structs.Message, error) {
	// New message
	var newMsg structs.Message

	// Set the value of the new message
	newMsg.Text = msg.Text
	newMsg.ConversationId = msg.ConversationId
	newMsg.Status = msg.Status
	newMsg.SenderUserId = msg.SenderUserId
	newMsg.Photo = msg.Photo

	// Get the last id
	maxId, err := db.GetMaxMessageId(newMsg.ConversationId)
	if err != nil {
		return structs.Message{}, err
	}

	// Set the id of the new message
	newMsg.MessageId = maxId + 1

	// Execute the query to create the conversation
	_, err = db.c.Exec(queryAddMessage, newMsg.MessageId, newMsg.Text, newMsg.Status, newMsg.ConversationId, newMsg.SenderUserId, newMsg.Photo)
	if err != nil {
		return structs.Message{}, err
	}

	// Return the new conversation
	return structs.Message{
		MessageId:      newMsg.MessageId,
		Text:           newMsg.Text,
		SendTime:       newMsg.SendTime,
		Status:         newMsg.Status,
		SenderUserId:   newMsg.SenderUserId,
		ConversationId: newMsg.ConversationId,
		Photo:          newMsg.Photo,
	}, nil
}
