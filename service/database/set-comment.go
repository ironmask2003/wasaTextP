package database

// Query to update a comment in the comment table
var querySetComment = `UPDATE comment SET Comment = ? WHERE (CommentId, CommentUserId, MessageId, ConversationId, UserId) = (?, ?, ?, ?, ?)`

func (db *appdbimpl) SetComment(commentId int, commentUserId int, msgId int, convId int, userId int, newComment string) error {
	// Upadte the comment in the comment table
	_, err := db.c.Exec(querySetComment, newComment, commentId, commentUserId, msgId, convId, userId)
	if err != nil {
		return err
	}
	return nil
}
