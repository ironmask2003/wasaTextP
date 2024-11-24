package database

import (
	"fmt"
	"io"
	"os"
	"wasa.project/service/api/functions"
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

	// Set default photo profile
	source, err := os.Open("./storage/default_propic.jpg") // Open the img file
	if err != nil {
		return user, err
	}
	defer source.Close()

	destination, err := os.Create(functions.SetDefaultPhoto(user.UserId)) // Create the path where the photo will be saved
	if err != nil {
		return user, err
	}
	defer destination.Close() // Close the user folder

	_, err = io.Copy(destination, source) // Copy the photo in the user folder
	if err != nil {
		return user, err
	}

	// Insert user in the database
	_, err = db.c.Exec(queryAddUser, user.UserId, user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}
