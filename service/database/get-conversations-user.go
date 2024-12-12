package database

import (
	"wasa.project/service/api/structs"
)

// Query used to get all conversations of a user
var queryGetConversations = `SELECT ConversationId FROM conversation_user WHERE UserId = ?`

func (db *appdbimpl) GetUserConversations(userId int) ([]structs.Conversation, error) {
	// Exec query
	rows, err := db.c.Query(queryGetConversations, userId)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	// All conversation
	var convs []structs.Conversation

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var conv structs.Conversation
		err = rows.Scan(&conv.ConversationId)
		if err != nil {
			return nil, err
		}
		conv, err = db.GetConversationById(conv.ConversationId)
		if err != nil {
			return nil, err
		}

		convs = append(convs, conv)
	}

	return convs, nil
}
