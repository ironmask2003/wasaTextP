package database

// Query used to get user of a conversation
var queryGetUsersConv = `SELECT UserId FROM conversation_user WHERE ConversationId = ? AND UserId != ?`

// Function
func (db *appdbimpl) GetUsersConv(convId int, userId int) (User, error) {
	var user User
	err := db.c.QueryRow(queryGetUsersConv, convId, userId).Scan(&user.UserId)
	return user, err
}
