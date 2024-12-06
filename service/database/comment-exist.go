// Control if a comment exist
package database

// Query check if a comment exist
var queryCommentExist = "SELECT CommentId FROM comment WHERE CommentUserId = ? AND MessageId = ? AND UserId = ? AND ConversationId = ?"

func (db *appdbimpl) ExistComment(messageId int, userId int, convId int, cUserId int) (bool, error) {
	var _commentId int
	err := db.c.QueryRow(queryCommentExist, cUserId, messageId, userId, convId).Scan(&_commentId)
	if err != nil {
		return false, err
	}
	return true, nil
}
