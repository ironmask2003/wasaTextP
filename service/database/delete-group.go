package database

import (
	"database/sql"
	"fmt"
	"os"
)

// Query used to delete a group from the group_t table
var queryDeleteGroup = `DELETE FROM group_t WHERE GroupId = ?`

func (db *appdbimpl) DeleteGroup(groupId int) error {
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

	// Delete the group in the group table
	_, err = tx.Exec(queryDeleteGroup, groupId)
	if err != nil {
		return err
	}

	// Remove all member from the group
	err = db.DeleteMember(groupId, tx)
	if err != nil {
		return err
	}

	// Delete the groups folder
	path := "./storage/groups/" + fmt.Sprint(groupId) + "/"
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil

}
