package database

import (
	"database/sql"
	"errors"
)

// Query for create a conversation in the conversation table
var queryAddConversationGroup = "INSERT INTO conversation (ConversationId, UserId, GroupId) VALUES (?, ?, ?)"

var queryAddConversationUser = "INSERT INTO conversation (ConversationId, UserId, SenderUserId) VALUES (?, ?, ?)"

// Query for get the last ID used in the conversation table
var queryLastConversationId = "SELECT MAX(ConversationId) FROM conversation WHERE UserId = ?"

func GetLastElemConversation(db *appdbimpl, userId int) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryLastConversationId, userId)
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

func (db *appdbimpl) CreateConversation(c Conversation) (Conversation, error) {
	var conv Conversation

	// Set the user id
	conv.UserId = c.UserId

	// Getting the max id in the conversation table
	maxID, err := GetLastElemConversation(db, conv.UserId)
	if err != nil {
		return conv, err
	}

	// Setting the id of the new group
	conv.ConversationId += maxID + 1

	// Check if is a conversation with a group or a user
	if c.GroupId != 0 {
		conv.GroupId = c.GroupId
		users, err := db.GetMembers(conv.GroupId)
		if err != nil {
			return conv, err
		}
		for i := 0; i < len(users); i++ {
			if users[i].UserId != conv.UserId {
				max_user, err := GetLastElemConversation(db, users[i].UserId)
				if err != nil {
					return conv, err
				}
				_, err = db.c.Exec(queryAddConversationGroup, max_user+1, users[i].UserId, conv.GroupId)
				if err != nil {
					return conv, err
				}
			}
		}
		// -- INSERT THE CONVERSATION IN THE DATABSE -- //
		_, err = db.c.Exec(queryAddConversationGroup, conv.ConversationId, conv.UserId, conv.GroupId)
		if err != nil {
			return conv, err
		}
	} else {
		conv.SenderUserId = c.SenderUserId
		// -- INSERT THE CONVERSATION IN THE DATABSE -- //
		_, err = db.c.Exec(queryAddConversationUser, conv.ConversationId, conv.UserId, conv.SenderUserId)
		if err != nil {
			return conv, err
		}

	}

	return conv, nil
}
