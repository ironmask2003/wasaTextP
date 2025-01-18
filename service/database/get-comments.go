package database

import (
	"wasa.project/service/api/structs"
)

// Query used to take the comment of message
var queryGetMsgComments = `SELECT CommentId, Comment, CommentUserId FROM comment WHERE MessageId = ? AND ConversationId = ?`

// function to get all comments of a message
func (db *appdbimpl) GetMsgComments(msgId int, convId int) ([]structs.RspComment, error) {
	// Exec query
	rows, err := db.c.Query(queryGetMsgComments, msgId, convId)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	// All conversation
	var comments []structs.RspComment

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var cmt structs.RspComment
		err = rows.Scan(&cmt.CommentId, &cmt.Comment, &cmt.CommentUserId)
		if err != nil {
			return nil, err
		}

		var usr User
		usr, err = db.GetUserById(cmt.CommentUserId)
		if err != nil {
			return nil, err
		}

		cmt.CommentUsername = usr.Username

		comments = append(comments, cmt)
	}

	return comments, nil

}
