package database

// Query used to get a commentm in the comment table
var queryGetComment = `SELECT CommentId, Comment FROM comment WHERE CommentUserId = ? AND MessageId = ? AND ConversationId = ? AND UserId = ?`

func (db *appdbimpl) GetComment(commentUserId int, msgId int, convId int, userId int) (Comment, error) {
	var com Comment
	err := db.c.QueryRow(queryGetComment, commentUserId, msgId, convId, userId).Scan(&com.CommentId, &com.Comment)
	if err != nil {
		return com, err
	}
	return com, nil
}
