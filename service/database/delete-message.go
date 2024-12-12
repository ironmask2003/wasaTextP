package database

import (
	"database/sql"
)

// Query used to remeove a message from the db
var queryDeleteMessage = `DELETE FROM message WHERE MessageId = ? AND ConversationId = ?`

func (db *appdbimpl) DeleteMessage(messageId int, convId int) error {
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
	_, err = tx.Exec(queryDeleteMessage, messageId, convId)
	if err != nil {
		return err
	}

	return nil
}
