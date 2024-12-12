package database

import (
	"database/sql"
	"errors"
	"wasa.project/service/api/structs"
)

// Query used to insert new comment in the db
var queryAddComment = `INSERT INTO comment (CommentId, CommentUserId, Comment, MessageId, ConversationId) VALUES (?, ?, ?, ?, ?)`

// Query used to take last id in the conversation table
var queryGetMaxComment = `SELECT MAX(CommentId) FROM comment WHERE MessageId = ? AND ConversationId = ?`

// Function used to take the max id in the comment table
func GetMaxCommentId(db *appdbimpl, messageId int, convId int) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryGetMaxComment, messageId, convId)
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

func (db *appdbimpl) CreateComment(c structs.Comment) (structs.Comment, error) {
	// New comment
	var newComment structs.Comment
	// Set values of new comment
	newComment.CommentUserId = c.CommentUserId
	newComment.Comment = c.Comment
	newComment.MessageId = c.MessageId
	newComment.ConversationId = c.ConversationId

	// Get the id of the comment
	maxId, err := GetMaxCommentId(db, newComment.MessageId, newComment.ConversationId)
	if err != nil {
		return structs.Comment{}, err
	}

	// Set the id
	newComment.CommentId = maxId + 1

	// Exec query
	_, err = db.c.Exec(queryAddComment, newComment.MessageId, newComment.ConversationId, newComment.Comment, newComment.MessageId, newComment.ConversationId)
	if err != nil {
		return structs.Comment{}, err
	}

	return structs.Comment{
		CommentId:      newComment.CommentId,
		CommentUserId:  newComment.CommentUserId,
		Comment:        newComment.Comment,
		MessageId:      newComment.MessageId,
		ConversationId: newComment.ConversationId,
	}, nil
}
