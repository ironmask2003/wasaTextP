package database

import (
	"database/sql"
	"errors"
)

// Query used to insert comment in a table
var queryAddComment = "INSERT INTO comment (CommentId, CommentUserId, Comment, MessageId, UserId, ConversationId) VALUES (?, ?, ?, ?, ?, ?)"

// Query used to get the last id in the comment table
var queryGetLastComment = "SELECT MAX(CommentId) FROM comment WHERE MessageId = ? AND UserId = ? AND ConversationId = ?"

// Function used to get last element in the comment table
func GetLastElemComment(db *appdbimpl, userId int, conversationId int, messageId int) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryGetLastComment, messageId, userId, conversationId)
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

func (db *appdbimpl) CreateComment(c Comment) (Comment, error) {
	var com Comment
	com.ConversationId = c.ConversationId
	com.CommentUserId = c.CommentUserId
	com.UserId = c.UserId
	com.MessageId = c.MessageId
	com.Comment = c.Comment

	// Get the las id in the comment table
	maxId, err := GetLastElemComment(db, com.UserId, com.ConversationId, com.MessageId)
	if err != nil {
		return com, err
	}

	com.CommentId = maxId + 1

	// Insert the comment in the table
	_, err = db.c.Exec(queryAddComment, com.CommentId, com.CommentUserId, com.Comment, com.MessageId, com.UserId, com.ConversationId)
	if err != nil {
		return com, err
	}

	return com, nil
}
