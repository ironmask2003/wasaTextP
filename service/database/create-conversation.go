package database

import (
	"database/sql"
	"errors"
)

// Query used ty insert a new conversation in the conversation table
var queryAddConversation = `INSERT INTO conversation (ConversationId, SenderUserId, GroupId, UserId, LastMessageId) VALUES (?, ?, ?, ?, ?);`

// Query for take the max id in the conversation
var queryMaxConversationID = `SELECT MAX(ConversationId) FROM conversation WHERE userId = ?;`

// Funciont used to take max ID in the conversation table
func GetLastElemConversation(userId int, db *appdbimpl) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryMaxConversationID, userId)
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

func (db *appdbimpl) CreateConversation(c Conversation, m Message) (Conversation, error) {
	var conv Conversation
	conv.UserId = c.UserId // Set the user who want to create a converastion

	// Check if the conversation is with a user or a group
	if c.GroupId != 0 {
		// Set the sender group id
		conv.GroupId = c.GroupId
		// Se the sende user id at 0
		conv.SenderUserId = 0
	} else {
		// Set the sender user id
		conv.SenderUserId = c.SenderUserId
		// Set the group id at 0
		conv.GroupId = 0
		// Set the message id in the conversation
		conv.LastMessageId = m.MessageId
	}

	// Getting the max id in the conversation table
	maxID, err := GetLastElemConversation(conv.UserId, db)
	if err != nil {
		return conv, err
	}

	// Setting the id of the new conversation
	conv.ConversationId += maxID + 1

	// Inizialize transation with the db and check if the is open correctly
	tx, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return conv, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
		err = tx.Commit()
	}()

	// Insert in the conversation table the new conversation of the user
	_, err = db.c.Exec(queryAddConversation, conv.ConversationId, conv.SenderUserId, conv.GroupId, conv.UserId, conv.LastMessageId)
	if err != nil {
		return conv, err
	}

	return conv, nil
}
