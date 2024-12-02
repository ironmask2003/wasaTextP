package database

import (
	"database/sql"
	"errors"
)

// Query for create a message in the message table
var queryAddMessage = "INSERT INTO message (MessageId, Message, SendUserId, ConversationId) VALUES (?, ?, ?, ?) "

// Query for get the last ID used in the conversation table
var queryLastMessageId = "SELECT MAX(MessageId) FROM message WHERE ConversationId = ?"

func GetLastElemMessage(db *appdbimpl, convId int) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryLastMessageId, convId)
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

func (db *appdbimpl) CreateMessage(m Message) (Message, error) {
	var msg Message

	// Set the user id, conversation id and the text
	msg.ConversationId = m.ConversationId
	msg.SenderUserId = m.SenderUserId
	msg.Text = m.Text

	// Getting the max id in the conversation table
	maxID, err := GetLastElemMessage(db, msg.ConversationId)
	if err != nil {
		return msg, err
	}

	// Setting the id of the new group
	msg.MessageId += maxID + 1

	// -- INSERT THE MESSAGE IN THE DATABSE -- //
	_, err = db.c.Exec(queryAddMessage, msg.MessageId, msg.Text, msg.SenderUserId, msg.ConversationId)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
