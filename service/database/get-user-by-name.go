package database

// Query for search the user by username
var queryFindUserByUsername = `SELECT UserId, Username FROM user WHERE Username = ?;`

func (db *appdbimpl) GetUserByName(username string) (User, error) {
	var user User
	err := db.c.QueryRow(queryFindUserByUsername, username).Scan(&user.UserId, &user.Username)
	return user, err
}
