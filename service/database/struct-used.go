package database

// -- Structs used in the database for the user containning the userId and the username -- //
type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
}

// -- Structs for the group -- //
type Group struct {
	GroupId   int    `json:"groupId"`
	GroupName string `json:"groupName"`
}
