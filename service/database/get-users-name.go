package database

import (
	_ "github.com/mattn/go-sqlite3"
)

var queryGetUsers = `SELECT UserId, Username FROM user WHERE Username regexp ?`

func (db *appdbimpl) SearchUsers(search string) ([]User, error) {
	var users []User

	rows, err := db.c.Query(queryGetUsers, "^"+search)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var u User
		if err := rows.Scan(&u.UserId, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	defer func() { err = rows.Close() }()

	return users, err
}
