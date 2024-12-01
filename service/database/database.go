/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// -- USER OPERATION -- //

	// Creation of new user in the user table
	CreateUser(u User) (User, error)

	// Change the username of the user
	SetMyUsername(UserId int, newUsername string) error

	// Delete a user in the user table
	DeleteUser(UserId int) error

	// Check if the username is alredy used
	CheckIfExist(username string) (bool, error)

	// Get User information from the db with the username
	GetUserByName(username string) (User, error)

	// Get User information from the db with the id
	GetUserById(userId int) (User, error)

	// Delete a member from a group
	LeaveGroup(UserId int, GroupId int) error

	// Get all members of a group
	GetMembers(groupId int) ([]User, error)

	// Set new group name
	SetGroupName(GroupId int, newName string) error

	// -- GROUP OPERATION -- //

	// Add a user to a group
	AddUserGroup(userId int, groupId int) error

	// Create Group
	CreateGroup(group Group, userId int) (Group, error)

	// Get Groiup information from the db with the id
	GetGroupById(groupId int) (Group, error)

	// Check if a user is member of a group
	CheckMember(userId int, groupId int) (bool, error)

	// Delete group
	DeleteGroup(groupId int) error

	// Delete all the user from the user_group table there are member of the group
	DeleteMember(groupId int, tx *sql.Tx) error

	Ping() error
}

type appdbimpl struct {
	c   *sql.DB
	ctx context.Context
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {

	// Check if the database is nil (required)
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	/// Check if the database is empty
	var tableSQL uint8
	err := db.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type='table'").Scan(&tableSQL)
	if err != nil {
		return nil, fmt.Errorf("error checking if database is empty: %w", err)
	}

	// Check of the number of table is corret (there are 5 tables)
	// if the number of table is not 5, we creating missing tables
	if tableSQL != 5 {

		// Craetion of the user tabel
		_, err = db.Exec(userTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user: %w", err)
		}

		// Creation of the message table
		_, err = db.Exec(messageTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure message: %w", err)
		}

		// Creation of the group table
		_, err = db.Exec(groupTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure group: %w", err)
		}

		// Creation of the user_group table
		_, err = db.Exec(userGroupTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user and group: %w", err)
		}

		// Creation of the conversation table
		_, err = db.Exec(conversationTableSQL)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure conversation: %w", err)
		}
	}

	return &appdbimpl{
		c:   db,
		ctx: context.Background(),
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
