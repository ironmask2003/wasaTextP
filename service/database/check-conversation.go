package database

import (
	"database/sql"
	"errors"
)

// Query to search a conversation in the database with a specified user
var queryFindConversationWithUser = `SELECT ConversationId FROM conversation WHERE SendUserId = ? AND UserId = ?;`

func (db *appdbimpl) CheckIfExistConversationWithUser(userId int, senderUsedId int) (bool, error) {
	var existsConversation string
	err := db.c.QueryRow(queryFindConversationWithUser, senderUsedId, userId).Scan(&existsConversation)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	// If exist the function return true (existName != "")
	return existsConversation != "", err
}

// Query to search a conversation in the database with a specified group
var queryFindConversationWithGroup = `SELECT ConversationId FROM conversation WHERE GroupId = ? AND UserId = ?;`

func (db *appdbimpl) CheckIfExistConversationWithGroup(userId int, groupId int) (bool, error) {
	var existsConversation string
	err := db.c.QueryRow(queryFindConversationWithGroup, groupId, userId).Scan(&existsConversation)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	// If exist the function return true (existName != "")
	return existsConversation != "", err
}
