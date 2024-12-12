package database

// Query used to update comment
var queryUpdateComment = `UPDATE comment SET Coment = ? WHERE CommentId = ? AND MessageId = ? AND conversationId = ?`

func (db *appdbimpl) UpdateComment(commentId int, messageId int, conversationId int) error {
	_, err := db.c.Exec(queryUpdateComment, commentId, messageId, conversationId)
	return err
}
