package database

import (
	"wasa.project/service/api/structs"
)

// Query used to take the conversation from the db with the id
var queryGetConv = `SELECT ConversationId, COALESCE(GroupId, 0), COALESCE(LastMessageId, 0) AS LastMessageId FROM conversation WHERE ConversationId = ?`

func (db *appdbimpl) GetConversationById(convId int) (structs.Conversation, error) {
	var conv structs.Conversation
	// Exec the query
	err := db.c.QueryRow(queryGetConv, convId).Scan(&conv.ConversationId, &conv.GroupId, &conv.LastMessageId)
	if err != nil {
		return structs.Conversation{}, err
	}
	return conv, nil
}
