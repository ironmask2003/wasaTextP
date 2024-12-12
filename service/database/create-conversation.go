package database

import (
	"database/sql"
	"errors"
	"wasa.project/service/api/structs"
)

// Query used to create a conversation
var queryAddConversation = `INSERT INTO conversation (ConversationId, IsGroup) VALUES (?, ?);`

// Query used to take the max id from the conversation table
var queryGetMaxConvId = `SELECT MAX(ConversationId) FROM conversation;`

// Function used to get the last id in the conversation table
func GetMaxConversationId(db *appdbimpl) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryGetMaxConvId)
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

func (db *appdbimpl) CreateConversation(c structs.Conversation) (structs.Conversation, error) {
	// New conversation
	var newConv structs.Conversation
	// Set the value of the new conversation
	newConv.IsGroup = c.IsGroup
	newConv.LastMessageId = c.LastMessageId
	// Get the id of the new conversation
	maxID, err := GetMaxConversationId(db)
	if err != nil {
		return structs.Conversation{}, err
	}
	// Set the id of the new conversation
	newConv.ConversationId = maxID + 1
	// Execute the query to create the conversation
	_, err = db.c.Exec(queryAddConversation, newConv.ConversationId, newConv.IsGroup)
	if err != nil {
		return structs.Conversation{}, err
	}
	// Return the new conversation
	return structs.Conversation{
		ConversationId: newConv.ConversationId,
		IsGroup:        newConv.IsGroup,
		LastMessageId:  newConv.LastMessageId,
	}, nil
}
