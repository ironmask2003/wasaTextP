package database

var queryFindUserById = `SELECT UserId, Username FROM user WHERE UserId = ?;`

func (db *appdbimpl) GetUserById(userId int) (User, error) {
	var user User
	err := db.c.QueryRow(queryFindUserById, userId).Scan(&user.UserId, &user.Username)
	return user, err
}
