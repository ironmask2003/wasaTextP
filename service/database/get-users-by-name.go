package database

var query_GETUSERS = `SELECT UserId, Username FROM user WHERE Username = ?;`

func (db *appdbimpl) SearchUsers(userID int, search string) (User, error) {
	var user User
	err := db.c.QueryRow(query_GETUSERS, search).Scan(&user.UserId, &user.Username)
	return user, err
}
