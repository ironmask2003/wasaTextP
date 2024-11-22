package database

import (
	"fmt"
	"os"
)

// Query for add a new user in the user table
var queryAddUser = "INSERT INTO user (UserId, Username) VALUES (?, ?);"

// Query for take the max id in the user table
var queryMaxUserID = "SELECT MAX(UserId) FROM user"

func (db *appdbimpl) CreateUser(u User) (User, error) {
	var user User
	user.Username = u.Username

	// Getting the max id in the user table
	maxID, err := db.GetLastElem(queryMaxUserID)
	if err != nil {
		return user, err
	}

	// Setting the id of the new user
	user.UserId += maxID + 1

	// Craetion of the user folder
	path := "./storage/" + fmt.Sprint(user.UserId) + "/conversations"
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return user, err
	}

	// Insert user in the database
	_, err = db.c.Exec(queryAddUser, user.UserId, user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}
