package database

var queryCommentExistId = `SELECT CommentId FROM comment WHERE CommentId = ? AND UserId = ? AND MessageId = ? AND ConversationId = ? AND CommentUserId = ?;`

func (db *appdbimpl) ExistCommentWithId(commentId int, messageId int, userId int, convId int, cUserId int) (bool, error) {
	var _commentId int
	err := db.c.QueryRow(queryCommentExistId, commentId, userId, messageId, convId, cUserId).Scan(&_commentId)
	if err != nil {
		return false, err
	}
	return true, nil

}
