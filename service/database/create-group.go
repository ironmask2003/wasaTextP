package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"wasa.project/service/api/imageFunctions"
)

// Query for create a group in the group table
var queryAddGroup = "INSERT INTO group_t (GroupId, GroupName) VALUES (?, ?)"

// Query for get the last ID used in the group table
var queryLastGroupID = "SELECT MAX(GroupId) FROM group_t"

func GetLastElemGroup(db *appdbimpl) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(queryLastGroupID)
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

func (db *appdbimpl) CreateGroup(group Group, userId int, convId int) (Group, error) {

	// -- SET THE NEW GROUP STRUCT -- //
	var g Group
	g.GroupName = group.GroupName

	// Getting the max id in the gorup table
	maxID, err := GetLastElemGroup(db)
	if err != nil {
		return g, err
	}

	// Setting the id of the new group
	g.GroupId += maxID + 1

	// -- SET THE PHOTO OF THE GROUP AND THE FOLDER -- //
	// Craetion of the group folder
	path := "./storage/groups/" + fmt.Sprint(g.GroupId)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return g, err
	}

	// Set default photo of the group
	source, err := os.Open("./storage/default_profile_photo.jpg") // Open the img file
	if err != nil {
		return g, err
	}
	defer source.Close()

	destination, err := os.Create(imageFunctions.SetDefaultPhotoGroup(g.GroupId)) // Create the path where the photo will be saved
	if err != nil {
		return g, err
	}
	defer destination.Close() // Close the user folder

	_, err = io.Copy(destination, source) // Copy the photo in the group folder
	if err != nil {
		return g, err
	}

	// -- INSERT THE GROUP IN THE DATABSE -- //
	_, err = db.c.Exec(queryAddGroup, g.GroupId, g.GroupName)
	if err != nil {
		return g, err
	}

	err = db.AddUserGroup(userId, g.GroupId)
	if err != nil {
		return g, err
	}

	err = db.AddUserConv(convId, userId)
	if err != nil {
		return g, err
	}

	return g, nil
}
