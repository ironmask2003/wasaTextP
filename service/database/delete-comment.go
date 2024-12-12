package database

import (
	"database/sql"
)

// Query used to delete comment
var queryDeleteComment = `DELETE FROM comment WHERE CommentId = ? AND MessageId = ? AND ConversationId = ?`

func (db *appdbimpl) DeleteComment(commentId int, messageId int, convId int) error {
	tx, err := db.c.BeginTx(db.ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		}
		err = tx.Commit()
	}()

	// Exec query
	_, err = tx.Exec(queryDeleteComment, commentId, messageId, convId)
	if err != nil {
		return err
	}

	return nil

}
