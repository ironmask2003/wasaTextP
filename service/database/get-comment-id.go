package database

var queryCommentExistId = `SELECT CommentId FROM comment WHERE CommentId = ? AND UserId = ? AND MessageId = ? AND ConversationId = ? AND CommentUserId = ?;`

func (db *appdbimpl) ExistCommentWithId(commentId int, messageId int, userId int, convId int, cUserId int) (bool, error) {
	_, err := db.c.Exec(queryCommentExistId, commentId, userId, messageId, convId, cUserId)
	if err != nil {
		return false, err
	}
	return true, nil

}
