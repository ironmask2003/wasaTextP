package database

import (
	"wasa.project/service/api/structs"
)

// Query used to take the conversation from the db with the id
<<<<<<< HEAD
var queryGetConv = `SELECT ConversationId, IsGroup, COALESCE(LastMessageId, 0) AS LastMessageId FROM conversation WHERE ConversationId = ?`
=======
var queryGetConv = `SELECT ConversationId, IsGroup FROM conversation WHERE ConversationId = ?`
>>>>>>> 779b51a (Modified table and function)

func (db *appdbimpl) GetConversationById(convId int) (structs.Conversation, error) {
	var conv structs.Conversation
	// Exec the query
	err := db.c.QueryRow(queryGetConv, convId).Scan(&conv.ConversationId, &conv.IsGroup, &conv.LastMessageId)
	if err != nil {
		return structs.Conversation{}, err
	}
	return conv, nil
}
