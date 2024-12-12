package api

import (
	"regexp"
	"wasa.project/service/api/imageFunctions"
	"wasa.project/service/database"
)

type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

// Check if the userName respect the regex
func (u *User) IsValid() bool {
	validUser := regexp.MustCompile(`^[a-z][a-z0-9]{2,13}$`)
	return validUser.MatchString(u.Username)
}

// Function used to insert a user in the database

// This function convert the User struct used in the api package in the User struct used in the database package
func (u *User) ConvertUserForDB() database.User {
	return database.User{
		UserId:   u.UserId,
		Username: u.Username,
	}
}

// This function convert the User struct used in the database package in the User struct used in the api package
func (u *User) ConvertUserFromDB(user database.User) error {
	u.UserId = user.UserId
	u.Username = user.Username
	profilePhoto, err := imageFunctions.ImageToBase64(imageFunctions.SetDefaultPhoto(user.UserId))
	if err != nil {
		return err
	}
	u.Photo = profilePhoto
	return nil
}
