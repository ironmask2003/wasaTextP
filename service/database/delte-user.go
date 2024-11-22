package database

import (
	"database/sql"
	"fmt"
	"os"
)

// Query for delete a user in the user table
var queryDeleteUser = "DELETE FROM user WHERE UserId = ?;"

func (db *appdbimpl) DeleteUser(UserId int) error {

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

	// Delete the user in the user table
	_, err = tx.Exec(queryDeleteUser, UserId)
	if err != nil {
		return err
	}

	// Delete the user folder
	path := "./storage/" + fmt.Sprint(UserId) + "/"
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}
