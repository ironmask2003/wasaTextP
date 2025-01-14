package database

// Query used to update comment
var queryUpdateComment = `UPDATE comment SET Comment = ? WHERE CommentId = ? AND MessageId = ? AND conversationId = ?`

func (db *appdbimpl) UpdateComment(comment string, commentId int, messageId int, conversationId int) error {
	_, err := db.c.Exec(queryUpdateComment, comment, commentId, messageId, conversationId)
	return err
}
