package database

import (
	"database/sql"
)

// Query used to delete a conversation
var queryDeleteConv = `DELETE FROM conversation WHERE ConversationId = ?`

// Query used to delete a conversation
func (db *appdbimpl) DeleteConv(convId int) error {
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
	_, err = tx.Exec(queryDeleteConv, convId)
	if err != nil {
		return err
	}

	return nil

}
