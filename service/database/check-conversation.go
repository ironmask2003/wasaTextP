package database

import (
	"database/sql"
	"errors"
)

// Query to search a conversation in the database
var queryFindConversation = `SELECT ConversationId FROM conversation WHERE ConversationId = ? AND UserId = ?;`

func (db *appdbimpl) CheckIfExistConversation(userId int, convId int) (bool, error) {
	var existsConversation string
	err := db.c.QueryRow(queryFindConversation, convId, userId).Scan(&existsConversation)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	// If exist the function return true (existName != "")
	return existsConversation != "", err
}
