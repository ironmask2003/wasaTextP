package database

// Query to delete a comment in the comment table
var queryDeleteComment = `DELETE FROM comment WHERE CommentId = ? AND CommentUserId = ? AND MessageId = ? AND ConversationId = ? AND UserId = ?`

func (db *appdbimpl) DeleteComment(commentId int, commentUserId int, msgId int, convId int, userId int) error {
	// Delete the comment in the comment table
	_, err := db.c.Exec(queryDeleteComment, commentId, commentUserId, msgId, convId, userId)
	if err != nil {
		return err
	}
	return nil
}
