package database

import (
	"database/sql"
	"errors"
)

// Query for search user with a specified username
var queryFindUsername = `SELECT username FROM user WHERE username = ?`

func (db *appdbimpl) CheckIfExist(username string) (bool, error) {
	var existsName string
	err := db.c.QueryRow(queryFindUsername, username).Scan(&existsName)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	// If exist the function return true (existName != "")
	return existsName != "", err
}
