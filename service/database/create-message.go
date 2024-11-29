package database

import (
	"database/sql"
	"errors"
)

// Query used to insert a new message in the message table
var queryAddtMessage = "INSERT INTO message (MessageId, ConversationId, Text, SendTime, Status, SenderUserId) VALUES (?, ?, ?, ?, ?);"

// Query for take the max id in the message table
var queryMaxMessageID = "SELECT MAX(MessageId) FROM message WHERE ConversationId = ?" // Need to take the last message of that conversation

// Function used to take the last element of the message table
func GetLastElemMassage(convId int, db *appdbimpl) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryMaxMessageID, convId)
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

	// Check if the message is empty and set the text
	if m.Text == "" {
		msg.Text = "Photo without caption"
	} else {
		msg.Text = m.Text
	}

	// Set the status of the message
	msg.Status = m.Status
	msg.SenderUserId = m.SenderUserId
	msg.ConversationId = m.ConversationId

	// Set the id of the message
	maxId, err := GetLastElemMassage(m.ConversationId, db)
	if err != nil {
		return msg, err
	}
	msg.MessageId = maxId + 1

	// Inizialize transation with the db and check if the is open correctly
	tx, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return msg, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
		err = tx.Commit()
	}()

	// Insert in the conversation table the new conversation of the user
	_, err = tx.Exec(queryAddtMessage, msg.MessageId, msg.ConversationId, msg.Text, msg.SendTime, msg.Status, msg.SenderUserId)
	if err != nil {
		return msg, err
	}

	return msg, nil
}
