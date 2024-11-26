package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"wasa.project/service/api/imageFunctions"
)

// Query for add a new user in the user table
var queryAddUser = "INSERT INTO user (UserId, Username) VALUES (?, ?);"

// Query for take the max id in the user table
var queryMaxUserID = "SELECT MAX(UserId) FROM user"

// Function used to take the last element of a user
func GetLastElem(db *appdbimpl) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryMaxUserID)
	if err != nil {
		return 0, err
	}

	var maxID int
	for row.Next() {
		if row.Err() != nil {
			return 0, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}

	return maxID, nil
}

func (db *appdbimpl) CreateUser(u User) (User, error) {
	var user User
	user.Username = u.Username

	// Getting the max id in the user table
	maxID, err := GetLastElem(db)
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
	source, err := os.Open("./storage/default_profile_photo.jpg") // Open the img file
	if err != nil {
		return user, err
	}
	defer source.Close()

	destination, err := os.Create(imageFunctions.SetDefaultPhoto(user.UserId)) // Create the path where the photo will be saved
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
