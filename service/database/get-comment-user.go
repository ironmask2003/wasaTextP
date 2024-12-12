package database

import (
	"wasa.project/service/api/structs"
)

// Query used to take the comment of a user
var queryGetCommentByUser = `SELECT CommentId, Comment FROM comment WHERE CommentUserId = ? AND MessageId = ? AND ConversationId = ?`

func (db *appdbimpl) GetCommentByUser(userId int, messageId int, convId int) (structs.Comment, error) {
	var comment structs.Comment
	err := db.c.QueryRow(queryGetCommentByUser, userId, messageId, convId).Scan(&comment.CommentId, &comment.Comment)
	return comment, err
}
